package bundling

import (
	"context"
	"santapan/domain"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// PostgresRepositoryQueries defines the methods for querying the category repository.
type PostgresRepositoryQueries interface {
	GetByID(ctx context.Context, id int64) (domain.Bundling, error)
	FetchBundlingMenuByBundlingID(ctx context.Context, bundlingID int64) ([]domain.BundlingMenu, error)
	Fetch(ctx context.Context, cursor string, num int64) (res []domain.Bundling, nextCursor string, err error)
}

type PostgresRepositoryQueriesMenu interface {
	GetByID(ctx context.Context, id int64) (domain.Menu, error)
}

// PostgresRepositoryCommand defines the methods for executing commands on the category repository.
type PostgresRepositoryCommand interface {
}

//go:generate mockery --name CategoryRepository
type Service struct {
	postgresRepoQuery     PostgresRepositoryQueries
	postgresRepoCommand   PostgresRepositoryCommand
	postgresRepoQueryMenu PostgresRepositoryQueriesMenu
}

// NewService will create a new category service object.
func NewService(pq PostgresRepositoryQueries, pc PostgresRepositoryCommand, pm PostgresRepositoryQueriesMenu) *Service {
	return &Service{
		postgresRepoQuery:     pq,
		postgresRepoCommand:   pc,
		postgresRepoQueryMenu: pm,
	}
}

// Fill Bundling And Menu Data
func (c *Service) fillBundlingMenu(ctx context.Context, bundlings []domain.BundlingMenu) ([]domain.BundlingMenu, error) {
	g, ctx := errgroup.WithContext(ctx)
	mapMenu := make(map[int64]domain.Menu)
	mapBundling := make(map[int64]domain.Bundling)

	// Map all menu
	for _, bundling := range bundlings {
		mapMenu[bundling.Menu.ID] = domain.Menu{}
	}

	// Map all bundling
	for _, bundling := range bundlings {
		mapBundling[bundling.Bundling.ID] = domain.Bundling{}
	}

	// Fetch all menu
	chanMenu := make(chan domain.Menu)
	for menuID := range mapMenu {
		menuID := menuID
		g.Go(func() error {
			menu, err := c.postgresRepoQueryMenu.GetByID(ctx, menuID)
			if err != nil {
				return err
			}
			chanMenu <- menu
			return nil
		})
	}

	// Fetch all bundling
	chanBundling := make(chan domain.Bundling)
	for bundlingID := range mapBundling {
		bundlingID := bundlingID
		g.Go(func() error {
			bundling, err := c.postgresRepoQuery.GetByID(ctx, bundlingID)
			if err != nil {
				return err
			}
			chanBundling <- bundling
			return nil
		})
	}

	// Fill bundling and menu
	for i := 0; i < len(mapMenu)+len(mapBundling); i++ {
		select {
		case menu := <-chanMenu:
			mapMenu[menu.ID] = menu
		case bundling := <-chanBundling:
			mapBundling[bundling.ID] = bundling
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// Fill bundling and menu to bundlingMenu
	for i := range bundlings {
		bundlings[i].Menu = mapMenu[bundlings[i].Menu.ID]
		bundlings[i].Bundling = mapBundling[bundlings[i].Bundling.ID]
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return bundlings, nil
}

// Fetch retrieves all categories.
func (c *Service) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Bundling, nextCursor string, err error) {
	res, nextCursor, err = c.postgresRepoQuery.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}
	return
}

// FetchBundlingMenuByBundlingID retrieves all bundling menu by bundling ID.
func (c *Service) FetchBundlingMenuByBundlingID(ctx context.Context, bundlingID int64) ([]domain.BundlingMenu, error) {
	res, err := c.postgresRepoQuery.FetchBundlingMenuByBundlingID(ctx, bundlingID)
	logrus.Print(err)
	if err != nil {
		return nil, err
	}

	res, err = c.fillBundlingMenu(ctx, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetByID retrieves a category by its ID.
func (c *Service) GetByID(ctx context.Context, id int64) (domain.Bundling, error) {
	banner, err := c.postgresRepoQuery.GetByID(ctx, id)
	if err != nil {
		return domain.Bundling{}, err
	}
	return banner, nil
}

-- Drop `bundling_menu` first since it references `bundling` and `menu`
DROP TABLE IF EXISTS bundling_menu CASCADE;

-- Drop `bundling` table
DROP TABLE IF EXISTS bundling CASCADE;

-- Drop `menu` table
DROP TABLE IF EXISTS menu CASCADE;

-- Drop ENUM type for `bundling_type` if it exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bundling_type_enum') THEN
        DROP TYPE bundling_type_enum;
    END IF;
END;
$$;

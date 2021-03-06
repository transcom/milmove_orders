-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

-- https://stackoverflow.com/questions/8092086/create-postgresql-role-user-if-it-doesnt-exist
DO
$do$
BEGIN
  IF NOT EXISTS (
    SELECT -- SELECT list can stay empty for this
    FROM   pg_catalog.pg_roles
    WHERE  rolname = 'ecs_user') THEN

    -- New local user with password
    -- In Production this would be
    -- CREATE USER ecs_user WITH LOGIN;
    CREATE USER ecs_user WITH PASSWORD 'mysecretpassword';

    -- rds_iam is an empty role in development but not in production
    GRANT rds_iam TO ecs_user;

    -- Local user should have same privs as primary user
    GRANT master TO ecs_user;

    -- Modify existing tables, sequences, and functions
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ecs_user;
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ecs_user;
    GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO ecs_user;

    -- Modify future tables, sequences, and functions
    ALTER DEFAULT PRIVILEGES FOR ROLE master GRANT ALL PRIVILEGES ON TABLES TO ecs_user;
    ALTER DEFAULT PRIVILEGES FOR ROLE master GRANT ALL PRIVILEGES ON SEQUENCES TO ecs_user;
    ALTER DEFAULT PRIVILEGES FOR ROLE master GRANT ALL PRIVILEGES ON FUNCTIONS TO ecs_user;

  END IF;
END
$do$;

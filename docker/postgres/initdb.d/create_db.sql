CREATE
OR REPLACE FUNCTION create_crypto_db() RETURNS void AS
$$
BEGIN
    CREATE
EXTENSION dblink;
   IF
NOT EXISTS (
      SELECT
      FROM   pg_catalog.pg_roles
      WHERE  rolname = 'crypto') THEN

CREATE ROLE crypto WITH SUPERUSER LOGIN PASSWORD 'crypto';
END IF;
   IF
EXISTS (SELECT FROM pg_database WHERE datname = 'crypto') THEN
      RAISE NOTICE 'Database already exists';
ELSE
      PERFORM dblink_exec('dbname=' || current_database()
                        , 'CREATE DATABASE crypto');
END IF;
   GRANT ALL PRIVILEGES ON DATABASE
crypto TO crypto;
END
$$
LANGUAGE plpgsql VOLATILE;

SELECT create_crypto_db();


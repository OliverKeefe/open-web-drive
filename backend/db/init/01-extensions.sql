-- TODO: create a custom dockerfile image and install pg_uuidv7 to handle
-- uuidv7 generation at the database level instead of the application layer.
-- CREATE EXTENSION IF NOT EXISTS "pg_uuidv7";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
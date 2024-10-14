CREATE TABLE IF NOT EXISTS "session"."schema_migrations" (
    version bigint PRIMARY KEY,
    dirty boolean NOT NULL
);

ALTER TABLE IF EXISTS "session"."schema_migrations" OWNER to postgres;
GRANT ALL ON TABLE "session"."schema_migrations" TO postgres;
GRANT ALL ON TABLE "session"."schema_migrations" TO session_manager;
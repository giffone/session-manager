CREATE TABLE IF NOT EXISTS "session"."errors" (
	id serial NOT NULL,
	unique_key varchar(255) NOT NULL UNIQUE,
	error_body varchar(255) NOT NULL,
	request_data jsonb NULL,
	created timestamptz NOT NULL DEFAULT LOCALTIMESTAMP,
	CONSTRAINT errors_pkey PRIMARY KEY (id)
);

-- Permissions

ALTER TABLE "session"."errors" OWNER TO postgres;
GRANT ALL ON TABLE "session"."errors" TO postgres;
GRANT ALL ON TABLE "session"."errors" TO session_manager;
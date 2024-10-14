CREATE TABLE IF NOT EXISTS "session"."on_platform" (
	session_id uuid NOT NULL,
	session_type varchar(20) NULL,
	login varchar(50) NOT NULL,
	start_date_time timestamptz NOT NULL DEFAULT LOCALTIMESTAMP,
	end_date_time timestamptz NOT NULL DEFAULT LOCALTIMESTAMP,
	created timestamptz NOT NULL DEFAULT LOCALTIMESTAMP,
	CONSTRAINT unique_session_id_type UNIQUE (session_id, session_type),
	CONSTRAINT on_platform_login_fkey FOREIGN KEY (login) REFERENCES "env_tracker"."users"(login),
	CONSTRAINT on_platform_session_id_fkey FOREIGN KEY (session_id) REFERENCES "session"."on_campus"(id)
);

-- Permissions

ALTER TABLE "session"."on_platform" OWNER TO postgres;
GRANT ALL ON TABLE "session"."on_platform" TO postgres;
GRANT ALL ON TABLE "session"."on_platform" TO session_manager;
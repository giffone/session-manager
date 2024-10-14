CREATE TABLE IF NOT EXISTS "session"."on_campus" (
	id uuid NOT NULL,
	comp_name varchar(30) NOT NULL,
	ip_addr varchar(20) NULL,
	login varchar(50) NOT NULL,
	next_ping_sec int4 NOT NULL,
	start_date_time timestamptz NOT NULL DEFAULT LOCALTIMESTAMP,
	end_date_time timestamptz NOT NULL DEFAULT LOCALTIMESTAMP,
	created timestamptz NOT NULL DEFAULT LOCALTIMESTAMP,
	CONSTRAINT on_campus_pkey PRIMARY KEY (id),
	CONSTRAINT on_campus_comp_name_fkey FOREIGN KEY (comp_name) REFERENCES "env_tracker"."computers"(comp_name),
	CONSTRAINT on_campus_login_fkey FOREIGN KEY (login) REFERENCES "env_tracker"."users"(login)
);

-- Permissions

ALTER TABLE "session"."on_campus" OWNER TO postgres;
GRANT ALL ON TABLE "session"."on_campus" TO postgres;
GRANT ALL ON TABLE "session"."on_campus" TO session_manager;
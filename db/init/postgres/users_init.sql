CREATE DATABASE pinterest;
CREATE USER postgres;
ALTER ROLE postgres SUPERUSER PASSWORD "postgres";
DROP TABLE IF EXISTS public.users;

CREATE TABLE public.users (
	user_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	email varchar(255) NOT NULL UNIQUE,
	nickname varchar(255) NOT NULL UNIQUE,
	"password" varchar(255) NOT NULL,
	register_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP
);

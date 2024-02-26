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

DROP TABLE IF EXISTS public.pins;
CREATE TABLE public.pins (
	pin_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	author_id bigint NOT NULL,
	created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	caption TEXT,
	click_url TEXT,
	content_url TEXT,
	CONSTRAINT fk_author FOREIGN KEY(author_id) REFERENCES public.users(user_id)
);

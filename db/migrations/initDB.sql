DROP TABLE IF EXISTS public.image;
CREATE TABLE public.image (
	image_id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	"name" TEXT NOT NULL UNIQUE,
	created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS public.user;
CREATE TABLE public.user (
	user_id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	nickname TEXT NOT NULL UNIQUE CHECK(length(nickname) BETWEEN 3 AND 20),
	password_hash TEXT NOT NULL,
	register_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
	avatar_id BIGINT NULL,
	FOREIGN KEY(avatar_id) REFERENCES public.image(image_id) ON DELETE SET NULL
);

DROP TABLE IF EXISTS public.pin;
CREATE TABLE public.pin (
	pin_id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	author_id BIGINT NOT NULL,
	created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
	title TEXT NOT NULL DEFAULT '',
	"description" TEXT NOT NULL DEFAULT '',
	click_url TEXT NOT NULL DEFAULT '',
	content_id BIGINT NOT NULL,
	allow_comments BOOLEAN NOT NULL DEFAULT TRUE,
	FOREIGN KEY(author_id) REFERENCES public.user(user_id) ON DELETE CASCADE,
	FOREIGN KEY(content_id) REFERENCES public.image(image_id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS public.board;
CREATE TABLE public.board (
	board_id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	title TEXT NOT NULL,
	created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
	"description" TEXT NOT NULL DEFAULT '',
	cover_id BIGINT NULL,
	visibility TEXT NOT NULL DEFAULT 'public' CHECK(visibility='public' OR visibility='private'),
	FOREIGN KEY(cover_id) REFERENCES public.image(image_id) ON DELETE SET NULL
);

DROP TABLE IF EXISTS public.board_pin;
CREATE TABLE public.board_pin (
	board_id BIGINT NOT NULL,
	pin_id BIGINT NOT NULL,
	PRIMARY KEY (board_id, pin_id),
	FOREIGN KEY(board_id) REFERENCES public.board(board_id) ON DELETE CASCADE,
	FOREIGN KEY(pin_id) REFERENCES public.pin(pin_id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS public.board_author;
CREATE TABLE public.board_author (
	board_id BIGINT NOT NULL,
	author_id BIGINT NOT NULL,
	PRIMARY KEY (board_id, author_id),
	FOREIGN KEY(board_id) REFERENCES public.board(board_id) ON DELETE CASCADE,
	FOREIGN KEY(author_id) REFERENCES public.user(user_id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS public.like;
CREATE TABLE public.like (
	pin_id BIGINT NOT NULL,
	user_id BIGINT NOT NULL,
	created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (pin_id, user_id),
	FOREIGN KEY(pin_id) REFERENCES public.pin(pin_id) ON DELETE CASCADE,
	FOREIGN KEY(user_id) REFERENCES public.user(user_id) ON DELETE CASCADE
);

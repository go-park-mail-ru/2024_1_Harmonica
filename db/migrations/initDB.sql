DROP TABLE IF EXISTS public.user;
CREATE TABLE public.user (
	user_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	nickname TEXT NOT NULL UNIQUE CHECK(length(nickname)<=20 AND length(nickname) >= 3),
	"password" TEXT NOT NULL,
	register_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	avatar_url TEXT NOT NULL DEFAULT ''
);

DROP TABLE IF EXISTS public.pin;
CREATE TABLE public.pin (
	pin_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	author_id bigint NOT NULL,
	created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	title TEXT NOT NULL DEFAULT '',
	"description" TEXT NOT NULL DEFAULT '',
	click_url TEXT NOT NULL DEFAULT '',
	content_url TEXT NOT NULL,
	allow_comments BOOLEAN NOT NULL DEFAULT TRUE,
	FOREIGN KEY(author_id) REFERENCES public.user(user_id) 
);

CREATE TYPE VISIBILITY AS ENUM('private', 'public');

DROP TABLE IF EXISTS public.board;
CREATE TABLE public.board (
	board_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	title TEXT NOT NULL,
	created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	description TEXT NOT NULL DEFAULT '',
	cover_url TEXT NULL DEFAULT '',
	visibility_type VISIBILITY NOT NULL DEFAULT 'public'
);

DROP TABLE IF EXISTS public.board_pin;
CREATE TABLE public.board_pin (
	board_id bigint NOT NULL,
	pin_id bigint NOT NULL,
	PRIMARY KEY (board_id, pin_id),
	FOREIGN KEY(board_id) REFERENCES public.board(board_id) ON DELETE CASCADE,
	FOREIGN KEY(pin_id) REFERENCES public.pin(pin_id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS public.board_author;
CREATE TABLE public.board_author (
	board_id bigint NOT NULL,
	author_id bigint NOT NULL,
	PRIMARY KEY (board_id, author_id),
	FOREIGN KEY(board_id) REFERENCES public.board(board_id) ON DELETE CASCADE,
	FOREIGN KEY(author_id) REFERENCES public.user(user_id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS public.like;
CREATE TABLE public.like (
	pin_id bigint NOT NULL,
	user_id bigint NOT NULL,
	created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (pin_id, user_id),
	FOREIGN KEY(pin_id) REFERENCES public.pin(pin_id) ON DELETE CASCADE,
	FOREIGN KEY(user_id) REFERENCES public.user(user_id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS public.subscribe_on_board;
CREATE TABLE public.subscribe_on_board (
	user_id bigint NOT NULL,
	followed_board_id bigint NOT NULL,
	PRIMARY KEY (user_id, followed_board_id),
	FOREIGN KEY(user_id) REFERENCES public.user(user_id) ON DELETE CASCADE,
	FOREIGN KEY(followed_board_id) REFERENCES public.board(board_id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS public.subscribe_on_person;
CREATE TABLE public.subscribe_on_person (
	user_id bigint NOT NULL,
	followed_user_id bigint NOT NULL,
	PRIMARY KEY (user_id, followed_user_id),
	FOREIGN KEY(user_id) REFERENCES public.user(user_id) ON DELETE CASCADE,
	FOREIGN KEY(followed_user_id) REFERENCES public.user(user_id) ON DELETE CASCADE
);
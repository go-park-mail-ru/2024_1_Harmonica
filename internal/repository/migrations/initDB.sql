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

CREATE TYPE MESSAGE_STATUS AS ENUM('read', 'unread');

DROP TABLE IF EXISTS public.message;
CREATE TABLE public.message (
    message_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    sender_id bigint NOT NULL,
    receiver_id bigint NOT NULL,
    text TEXT NOT NULL,
    status MESSAGE_STATUS NOT NULL DEFAULT 'unread',
    sent_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(sender_id) REFERENCES public.user(user_id) ON DELETE CASCADE,
    FOREIGN KEY(receiver_id) REFERENCES public.user(user_id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS public.comment;
CREATE TABLE public.comment (
    comment_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id bigint NOT NULL,
    pin_id bigint NOT NULL,
    text TEXT NOT NULL,
    created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES public.user(user_id) ON DELETE CASCADE,
    FOREIGN KEY (pin_id) REFERENCES public.pin(pin_id) ON DELETE CASCADE
);

CREATE TYPE NOTIFICATION_TYPE AS ENUM('subscription', 'new_pin', 'comment', 'message');
CREATE TYPE NOTIFICATION_STATUS AS ENUM ('read', 'unread');

DROP TABLE IF EXISTS public.notification;
CREATE TABLE public.notification (
    notification_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id bigint NOT NULL,  -- пользователь, получающий уведомление
    type NOTIFICATION_TYPE NOT NULL,  -- тип уведомления
    triggered_by_user_id bigint NOT NULL,  -- пользователь, вызвавший уведомление (например, подписчик или комментатор)
    pin_id bigint,  -- идентификатор пина, если уведомление связано с пином
    comment_id bigint,  -- идентификатор комментария, если уведомление связано с комментарием
    message_id bigint,
    status NOTIFICATION_STATUS NOT NULL DEFAULT 'unread',
    created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES public.user(user_id) ON DELETE CASCADE,
    FOREIGN KEY (triggered_by_user_id) REFERENCES public.user(user_id) ON DELETE CASCADE,
    FOREIGN KEY (pin_id) REFERENCES public.pin(pin_id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES public.comment(comment_id) ON DELETE CASCADE,
    FOREIGN KEY (message_id) REFERENCES public.message(message_id) ON DELETE CASCADE
);

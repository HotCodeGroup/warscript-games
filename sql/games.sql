CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS "games" CASCADE;
CREATE TABLE "games"
(
	id bigserial not null
		constraint game_pk
			primary key,
	slug CITEXT UNIQUE CONSTRAINT games_slug_check CHECK ( slug ~ '^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$' ),
	title CITEXT UNIQUE CONSTRAINT title_empty not null check ( title <> '' ),
	description TEXT NOT NULL,
	rules TEXT NOT NULL,
	code_example TEXT NOT NULL,
	bot_code TEXT NOT NULL,
	logo_uuid UUID NOT NULL,
	background_uuid UUID NOT NULL
);
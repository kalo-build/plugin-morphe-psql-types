-- Table definition for comments

CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.comments (
	content TEXT NOT NULL,
	created_at TEXT NOT NULL,
	id SERIAL PRIMARY KEY,
	commentable_type TEXT NOT NULL,
	commentable_id TEXT NOT NULL
);


-- View definition for comment_entities

CREATE SCHEMA IF NOT EXISTS public;

CREATE OR REPLACE VIEW public.comment_entities AS
SELECT
	comments.content,
	comments.created_at,
	comments.id
FROM public.comments;


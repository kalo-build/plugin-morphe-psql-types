-- View definition for tag_entities

CREATE SCHEMA IF NOT EXISTS public;

CREATE OR REPLACE VIEW public.tag_entities AS
SELECT
	tags.color,
	tags.id,
	tags.name
FROM public.tags;


-- View definition for person_entities

CREATE SCHEMA IF NOT EXISTS public;

CREATE OR REPLACE VIEW public.person_entities AS
SELECT
	contact_infos.email,
	people.id,
	people.last_name,
	people.nationality
FROM public.people
LEFT JOIN public.contact_infos
	ON people.id = contact_infos.id;


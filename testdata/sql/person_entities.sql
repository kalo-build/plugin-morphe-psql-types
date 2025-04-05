-- View definition for person_entities

CREATE SCHEMA IF NOT EXISTS public;

CREATE OR REPLACE VIEW public.person_entities AS
SELECT
	persons.id,
	persons.last_name,
	persons.nationality,
	contact_infos.email
FROM public.persons
LEFT JOIN public.contact_infos
	ON persons.id = contact_infos.person_id
LEFT JOIN public.companies
	ON persons.company_id = companies.id; 
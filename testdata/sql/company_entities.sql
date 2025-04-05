-- View definition for company_entities

CREATE SCHEMA IF NOT EXISTS public;

CREATE OR REPLACE VIEW public.company_entities AS
SELECT
	companies.id,
	companies.name,
	companies.tax_id
FROM public.companies; 
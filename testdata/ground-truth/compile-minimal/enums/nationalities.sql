-- Table definition for nationalities

CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.nationalities (
	id SERIAL PRIMARY KEY,
	key TEXT NOT NULL,
	value TEXT NOT NULL,
	value_type TEXT NOT NULL,
	UNIQUE (key)
);

-- Seed Data
INSERT INTO public.nationalities (key, value, value_type) VALUES ('DE', 'German', 'String');
INSERT INTO public.nationalities (key, value, value_type) VALUES ('FR', 'French', 'String');
INSERT INTO public.nationalities (key, value, value_type) VALUES ('US', 'American', 'String');


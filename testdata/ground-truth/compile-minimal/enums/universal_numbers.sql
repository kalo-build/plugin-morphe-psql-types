-- Table definition for universal_numbers

CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.universal_numbers (
	id SERIAL PRIMARY KEY,
	key TEXT NOT NULL,
	value TEXT NOT NULL,
	value_type TEXT NOT NULL,
	UNIQUE (key)
);

-- Seed Data
INSERT INTO public.universal_numbers (key, value, value_type) VALUES ('Euler', '2.7182818285', 'Float');
INSERT INTO public.universal_numbers (key, value, value_type) VALUES ('Pi', '3.1415926535', 'Float');


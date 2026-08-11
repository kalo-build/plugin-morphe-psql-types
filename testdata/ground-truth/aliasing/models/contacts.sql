-- Table definition for contacts

CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.contacts (
	address TEXT,
	email TEXT,
	id SERIAL PRIMARY KEY,
	phone TEXT
);


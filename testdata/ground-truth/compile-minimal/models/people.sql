-- Table definition for people

CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.people (
	first_name TEXT,
	id SERIAL PRIMARY KEY,
	last_name TEXT,
	nationality_id INTEGER NOT NULL,
	company_id INTEGER NOT NULL,
	CONSTRAINT fk_people_nationality_id FOREIGN KEY (nationality_id)
		REFERENCES public.nationalities (id)
		ON DELETE CASCADE,
	CONSTRAINT fk_people_company_id FOREIGN KEY (company_id)
		REFERENCES public.companies (id)
		ON DELETE CASCADE
);

-- Indices
CREATE INDEX IF NOT EXISTS idx_people_nationality_id ON public.people (nationality_id);
CREATE INDEX IF NOT EXISTS idx_people_company_id ON public.people (company_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_people_first_name_last_name ON public.people (first_name, last_name);


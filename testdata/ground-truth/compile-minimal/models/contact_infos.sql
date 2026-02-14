-- Table definition for contact_infos

CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.contact_infos (
	email TEXT NOT NULL,
	id SERIAL PRIMARY KEY,
	person_id INTEGER NOT NULL,
	CONSTRAINT fk_contact_infos_person_id FOREIGN KEY (person_id)
		REFERENCES public.people (id)
		ON DELETE CASCADE
);

-- Indices
CREATE INDEX IF NOT EXISTS idx_contact_infos_person_id ON public.contact_infos (person_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_contact_infos_email ON public.contact_infos (email);


-- Table definition for people

CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.people (
	email TEXT,
	id SERIAL PRIMARY KEY,
	name TEXT,
	personal_contact_id INTEGER NOT NULL,
	work_contact_id INTEGER NOT NULL,
	CONSTRAINT fk_people_personal_contact_id FOREIGN KEY (personal_contact_id)
		REFERENCES public.contacts (id)
		ON DELETE CASCADE,
	CONSTRAINT fk_people_work_contact_id FOREIGN KEY (work_contact_id)
		REFERENCES public.contacts (id)
		ON DELETE CASCADE
);

-- Indices
CREATE INDEX IF NOT EXISTS idx_people_personal_contact_id ON public.people (personal_contact_id);
CREATE INDEX IF NOT EXISTS idx_people_work_contact_id ON public.people (work_contact_id);


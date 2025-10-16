CREATE TABLE IF NOT EXISTS public.people (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    personal_contact_id INTEGER NOT NULL,
    work_contact_id INTEGER NOT NULL
);

ALTER TABLE public.people
    ADD CONSTRAINT fk_people_personal_contact_id FOREIGN KEY (personal_contact_id) REFERENCES public.contacts(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_people_work_contact_id FOREIGN KEY (work_contact_id) REFERENCES public.contacts(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_people_personal_contact_id ON public.people (personal_contact_id);
CREATE INDEX IF NOT EXISTS idx_people_work_contact_id ON public.people (work_contact_id);
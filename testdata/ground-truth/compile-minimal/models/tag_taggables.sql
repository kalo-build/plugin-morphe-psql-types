-- Table definition for tag_taggables

CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.tag_taggables (
	id SERIAL PRIMARY KEY,
	tag_id INTEGER,
	taggable_type TEXT,
	taggable_id TEXT,
	UNIQUE (tag_id, taggable_type, taggable_id),
	CONSTRAINT fk_tag_taggables_tag_id FOREIGN KEY (tag_id)
		REFERENCES public.tags (id)
		ON DELETE CASCADE
);

-- Indices
CREATE INDEX IF NOT EXISTS idx_tag_taggables_tag_id ON public.tag_taggables (tag_id);


CREATE TABLE IF NOT EXISTS public.contacts (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    address TEXT NOT NULL
);

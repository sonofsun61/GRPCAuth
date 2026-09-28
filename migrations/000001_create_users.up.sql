create table users(
    id uuid primary key default gen_random_uuid(),
    email text not null unique,
    name text not null,
    surname text not null,
    phone_number text not null unique,
    password_hash text not null,
    password_salt text not null,
    created_at timestamptz not null default now()
);
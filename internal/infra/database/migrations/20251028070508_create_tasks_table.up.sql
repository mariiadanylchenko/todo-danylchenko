CREATE TABLE IF NOT EXISTS public.tasks
(
    id              serial PRIMARY KEY,
    user_id         integer references public.users(id),
    title           varchar(100) NOT NULL,
    description     text,
    status          varchar(20) NOT NULL,
    deadline        timestamptz,
    created_date    timestamptz NOT NULL,
    updated_date    timestamptz NOT NULL,
    deleted_date    timestamptz
);

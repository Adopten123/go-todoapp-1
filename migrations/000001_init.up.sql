CREATE SCHEMA IF NOT EXISTS todoapp;

CREATE TABLE todoapp.users
(
    id           UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    version      BIGINT       NOT NULL DEFAULT 1,
    full_name    VARCHAR(100) NOT NULL CHECK (char_length(full_name) BETWEEN 3 AND 100),
    phone_number VARCHAR(15) UNIQUE CHECK (
        phone_number ~ '^\+[0-9]+$'
        AND
        char_length (phone_number) BETWEEN 10 AND 15
    )
);

CREATE TABLE todoapp.tasks
(
    id             UUID                 PRIMARY KEY DEFAULT gen_random_uuid(),
    version        BIGINT       NOT NULL DEFAULT 1,
    title          VARCHAR(100) NOT NULL CHECK (char_length(title) BETWEEN 1 AND 100),
    description    VARCHAR(1000) CHECK (char_length(description) BETWEEN 1 AND 1000),
    completed      BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
    completed_at   TIMESTAMPTZ,

    CHECK (
        (completed = FALSE AND completed_at IS NULL)
            OR
        (completed = TRUE AND completed_at IS NOT NULL AND completed_at >= created_at)
        ),

    author_user_id UUID         NOT NULL REFERENCES todoapp.users (id) ON DELETE CASCADE
);

CREATE INDEX idx_tasks_author_id ON todoapp.tasks (author_user_id);
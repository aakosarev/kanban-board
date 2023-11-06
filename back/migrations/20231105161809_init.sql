-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL CHECK ( email <> '' ),
    password VARCHAR(255) NOT NULL CHECK ( octet_length(password) <> 0 )
);

CREATE TABLE IF NOT EXISTS "column" (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES "user"(id) ON DELETE CASCADE,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS "task" (
    id SERIAL PRIMARY KEY,
    column_id INT REFERENCES "column"(id) ON DELETE CASCADE,
    description TEXT
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS "column";
DROP TABLE IF EXISTS "task";
-- +goose StatementEnd

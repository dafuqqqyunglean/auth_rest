CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    guid          VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL
);
CREATE TABLE IF NOT EXISTS friends
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    email         VARCHAR(255)  NOT NULL,
    description   TEXT,
    requirement   TEXT,
    select_friend INT                   DEFAULT NULL,
    validate_code VARCHAR(100) NOT NULL,
    is_valid      BOOLEAN      NOT NULL DEFAULT FALSE
);

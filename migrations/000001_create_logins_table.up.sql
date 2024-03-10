CREATE TABLE IF NOT EXISTS logins(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    password_hash varchar(65) NOT NULL
);
CREATE TABLE IF NOT EXISTS logins(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    password_hash varchar(72) NOT NULL
);
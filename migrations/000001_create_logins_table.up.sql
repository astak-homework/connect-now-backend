CREATE TABLE IF NOT EXISTS logins(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_name varchar(50) UNIQUE NOT NULL,
    password_hash varchar(65) NOT NULL
);
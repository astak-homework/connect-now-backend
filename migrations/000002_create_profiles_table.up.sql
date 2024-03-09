CREATE TABLE IF NOT EXISTS profiles(
    id uuid PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    birth_date DATE,
    gender VARCHAR(6),
    biography TEXT,
    city VARCHAR(50)
);
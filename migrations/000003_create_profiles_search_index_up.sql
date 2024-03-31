CREATE INDEX IF NOT EXISTS profiles_first_name_last_name_idx
ON profiles (first_name varchar_pattern_ops, last_name varchar_pattern_ops);
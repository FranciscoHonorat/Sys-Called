ALTER TABLE employees ADD COLUMN IF NOT EXISTS username TEXT;
ALTER TABLE employees ADD COLUMN IF NOT EXISTS role TEXT;
ALTER TABLE employees ADD COLUMN IF NOT EXISTS password_hash TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_employees_username ON employees (username);

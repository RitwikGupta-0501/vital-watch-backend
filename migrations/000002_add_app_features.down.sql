-- Drop tables in reverse order
DROP TABLE IF EXISTS prescriptions;
DROP TABLE IF EXISTS appointments;

-- Remove columns from doctors
ALTER TABLE doctors
DROP COLUMN specialty,
DROP COLUMN experience,
DROP COLUMN available;
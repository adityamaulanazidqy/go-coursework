ALTER TABLE users
    DROP COLUMN IF EXISTS batch;

DELETE FROM status_submissions WHERE id IN (1, 2, 3);
DELETE FROM study_programs WHERE name IN ('Informatics Engineering', 'Information Systems', 'Informatics Management');
DELETE FROM roles WHERE name IN ('Admin', 'Lecturer', 'Student');
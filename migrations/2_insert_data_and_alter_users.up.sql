INSERT INTO roles (name) VALUES
('Admin'),
('Lecturer'),
('Student')
ON CONFLICT (name) DO NOTHING;

INSERT INTO study_programs (name) VALUES
('Informatics Engineering'),
('Information Systems'),
('Informatics Management')
ON CONFLICT (name) DO NOTHING;

INSERT INTO status_submissions (id, name) VALUES
(1, 'Not Yet Collected'),
(2, 'Pending'),
(3, 'Revision Required'),
(4, 'Accepted')
ON CONFLICT (id) DO NOTHING;

ALTER TABLE users
ADD COLUMN IF NOT EXISTS batch INT NOT NULL;

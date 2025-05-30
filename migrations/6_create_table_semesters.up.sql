CREATE TABLE IF NOT EXISTS semesters (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

INSERT INTO semesters (name) VALUES
    ('One'),
    ('Two'),
    ('Tree'),
    ('Four'),
    ('Five'),
    ('Six'),
    ('Seven'),
    ('Eighth'),
    ('Special')
ON CONFLICT (name) DO NOTHING;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS semester_id INT DEFAULT 1,
    ADD CONSTRAINT fk_semester FOREIGN KEY (semester_id) REFERENCES semesters(id) ON DELETE SET NULL;

ALTER TABLE assignments
    ADD COLUMN IF NOT EXISTS semester_id INT,
    ADD COLUMN IF NOT EXISTS study_program_id INT,
    ADD CONSTRAINT fk_asgn_semester FOREIGN KEY (semester_id) REFERENCES semesters(id) ON DELETE SET NULL,
    ADD CONSTRAINT fk_asgn_study_program FOREIGN KEY (study_program_id) REFERENCES study_programs(id) ON DELETE CASCADE;
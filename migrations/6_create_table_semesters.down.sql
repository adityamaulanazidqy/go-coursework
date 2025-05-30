ALTER TABLE assignments
    DROP CONSTRAINT IF EXISTS fk_asgn_study_program,
    DROP CONSTRAINT IF EXISTS fk_asgn_semester;

ALTER TABLE assignments
    DROP COLUMN IF EXISTS semester_id,
    DROP COLUMN IF EXISTS study_program_id;

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS fk_semester;

ALTER TABLE users
    DROP COLUMN IF EXISTS semester_id;

DROP TABLE IF EXISTS semesters;

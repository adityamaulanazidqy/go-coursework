ALTER TABLE submission_grades
    DROP CONSTRAINT IF EXISTS fk_submission_grade,
    DROP CONSTRAINT IF EXISTS fk_lecturer_grade;

DROP TABLE IF EXISTS submission_grades;

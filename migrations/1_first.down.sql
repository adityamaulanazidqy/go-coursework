DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP TRIGGER IF EXISTS trg_assignments_updated_at ON assignments;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS submission_histories;
DROP TABLE IF EXISTS submissions;
DROP TABLE IF EXISTS status_submissions;
DROP TABLE IF EXISTS assignment_files;
DROP TABLE IF EXISTS assignment_comments;
DROP TABLE IF EXISTS assignments;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS activity_logs;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS study_programs;
CREATE TABLE roles (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE study_programs (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL UNIQUE,
  email VARCHAR(255) NOT NULL UNIQUE,
  telephone VARCHAR(20) DEFAULT '0',
  study_program_id int NOT NULL,
  password VARCHAR(100) NOT NULL,
  role_id INT NOT NULL,
  profile TEXT DEFAULT 'icon_default.jpg',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  CONSTRAINT fk_study_program FOREIGN KEY (study_program_id) REFERENCES study_programs(id) ON DELETE CASCADE
);

CREATE TABLE assignments (
  id SERIAL PRIMARY KEY,
  lecturer_id INT NOT NULL,
  title VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  file_url TEXT DEFAULT NULL,
  deadline TIMESTAMP NOT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_lecturer FOREIGN KEY (lecturer_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE assignment_comments (
  id SERIAL PRIMARY KEY,
  assignment_id INT NOT NULL,
  user_id INT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_assignment_comment FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_comment FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE assignment_files (
  id SERIAL PRIMARY KEY,
  assignment_id INT NOT NULL,
  file_url TEXT NOT NULL,
  uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_assignment_file FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE CASCADE
);

CREATE TABLE status_submissions (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE submissions (
  id SERIAL PRIMARY KEY,
  assignment_id INT NOT NULL,
  student_id INT NOT NULL,
  file_url TEXT NOT NULL,
  submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  status_submissions_id INT DEFAULT 1,
  CONSTRAINT fk_assignment FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE CASCADE,
  CONSTRAINT fk_student FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_status FOREIGN KEY (status_submissions_id) REFERENCES status_submissions(id) ON DELETE CASCADE
);

CREATE TABLE submission_histories (
  id SERIAL PRIMARY KEY,
  submission_id INT NOT NULL,
  file_url TEXT NOT NULL,
  changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  notes TEXT DEFAULT NULL,
  CONSTRAINT fk_submission FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE
);

CREATE TABLE notifications (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  message TEXT NOT NULL,
  is_read BOOLEAN DEFAULT FALSE,
  CONSTRAINT fk_user_notification FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE activity_logs (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  activity TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_user_log FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_assignments_updated_at
    BEFORE UPDATE ON assignments
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


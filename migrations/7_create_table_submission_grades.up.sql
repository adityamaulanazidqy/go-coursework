CREATE TABLE IF NOT EXISTS submission_grades (
    id SERIAL PRIMARY KEY,
    submission_id INT NOT NULL,
    lecturer_id INT NOT NULL,
    grade INT NOT NULL,
    notes TEXT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_submission_grade FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE,
    CONSTRAINT fk_lecturer_grade FOREIGN KEY (lecturer_id) REFERENCES users(id) ON DELETE CASCADE
);
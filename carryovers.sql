CREATE TABLE IF NOT EXISTS carryovers_2024_2025 (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES students_2024_2025(student_id) NOT NULL,
    subject_id INTEGER REFERENCES subjects_2024_2025(subject_id) NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (student_id, subject_id)
);
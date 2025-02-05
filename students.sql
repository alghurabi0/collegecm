CREATE TABLE IF NOT EXISTS students_2024_2025 (
    seq_in_college SERIAL,
    student_name VARCHAR(255) NOT NULL,
    stage VARCHAR(100) NOT NULL,
    student_id INTEGER NOT NULL PRIMARY KEY,
    state VARCHAR(100) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);
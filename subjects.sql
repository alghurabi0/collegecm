CREATE TABLE IF NOT EXISTS subjects_2024_2025(
   subject_id           INTEGER NOT NULL PRIMARY KEY 
  ,subject_name         VARCHAR(100) NOT NULL
  ,subject_name_english VARCHAR(100) NOT NULL
  ,stage                VARCHAR(30) NOT NULL
  ,semester             VARCHAR(30) NOT NULL
  ,department           VARCHAR(100) NOT NULL
  ,max_theory_mark      INTEGER  NOT NULL
  ,max_lab_mark         INTEGER  NOT NULL
  ,max_semester_mark    INTEGER  NOT NULL
  ,max_final_exam       INTEGER  NOT NULL
  ,credits              INTEGER  NOT NULL
  ,active               VARCHAR(10) NOT NULL
  ,ministerial          VARCHAR(10) NOT NULL
  ,created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
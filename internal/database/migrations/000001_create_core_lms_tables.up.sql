-- Migration: 000001_create_core_lms_tables
-- Created: YYYY-MM-DD HH:MM:SS
-- Description: Creates the initial tables for the Core LMS module.

BEGIN;

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ, -- For GORM soft delete
    clerk_user_id VARCHAR(255) UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL, -- Store hashed password
    role VARCHAR(50) NOT NULL CHECK (role IN ('student', 'teacher', 'admin', 'parent')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login TIMESTAMPTZ,
    pdpa_consent_timestamp TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_last_login ON users(last_login);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
COMMENT ON COLUMN users.password IS 'Stores the hashed password using bcrypt';
COMMENT ON COLUMN users.role IS 'User role: student, teacher, admin, parent';

-- Courses Table
CREATE TABLE IF NOT EXISTS courses (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    teacher_id BIGINT NOT NULL,
    subject VARCHAR(100),
    grade_level VARCHAR(50),
    is_published BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_courses_teacher_id ON courses(teacher_id);
CREATE INDEX IF NOT EXISTS idx_courses_subject ON courses(subject);
CREATE INDEX IF NOT EXISTS idx_courses_grade_level ON courses(grade_level);
CREATE INDEX IF NOT EXISTS idx_courses_is_published ON courses(is_published);
CREATE INDEX IF NOT EXISTS idx_courses_deleted_at ON courses(deleted_at);
ALTER TABLE courses
ADD CONSTRAINT fk_courses_teacher
FOREIGN KEY (teacher_id) REFERENCES users(id)
ON UPDATE CASCADE ON DELETE RESTRICT; -- Prevent deleting teacher with courses

-- Sections Table
CREATE TABLE IF NOT EXISTS sections (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    course_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    order_index INT NOT NULL DEFAULT 0,
    is_published BOOLEAN NOT NULL DEFAULT TRUE
);
CREATE INDEX IF NOT EXISTS idx_sections_course_id ON sections(course_id);
CREATE INDEX IF NOT EXISTS idx_sections_order_index ON sections(order_index);
CREATE INDEX IF NOT EXISTS idx_sections_is_published ON sections(is_published);
CREATE INDEX IF NOT EXISTS idx_sections_deleted_at ON sections(deleted_at);
-- Optional: Unique constraint for order within a course
-- ALTER TABLE sections ADD CONSTRAINT uq_section_order UNIQUE (course_id, order_index);
ALTER TABLE sections
ADD CONSTRAINT fk_sections_course
FOREIGN KEY (course_id) REFERENCES courses(id)
ON UPDATE CASCADE ON DELETE CASCADE; -- Delete sections if course is deleted

-- Lessons Table
CREATE TABLE IF NOT EXISTS lessons (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    section_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    content_type VARCHAR(50) NOT NULL, -- 'text', 'video', 'quiz', 'assignment', 'document'
    content_data JSONB NOT NULL DEFAULT '{}',
    order_index INT NOT NULL DEFAULT 0,
    is_published BOOLEAN NOT NULL DEFAULT TRUE,
    estimated_duration_minutes INT
);
CREATE INDEX IF NOT EXISTS idx_lessons_section_id ON lessons(section_id);
CREATE INDEX IF NOT EXISTS idx_lessons_content_type ON lessons(content_type);
CREATE INDEX IF NOT EXISTS idx_lessons_order_index ON lessons(order_index);
CREATE INDEX IF NOT EXISTS idx_lessons_is_published ON lessons(is_published);
CREATE INDEX IF NOT EXISTS idx_lessons_deleted_at ON lessons(deleted_at);
-- Optional: Unique constraint for order within a section
-- ALTER TABLE lessons ADD CONSTRAINT uq_lesson_order UNIQUE (section_id, order_index);
ALTER TABLE lessons
ADD CONSTRAINT fk_lessons_section
FOREIGN KEY (section_id) REFERENCES sections(id)
ON UPDATE CASCADE ON DELETE CASCADE; -- Delete lessons if section is deleted

-- Enrollments Table
CREATE TABLE IF NOT EXISTS enrollments (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    user_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,
    enrollment_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'withdrawn')),
    completed_at TIMESTAMPTZ,
    progress DECIMAL(5,4) NOT NULL DEFAULT 0.0 CHECK (progress >= 0.0 AND progress <= 1.0),
    UNIQUE (user_id, course_id, deleted_at) -- Allow re-enrollment after soft delete? Or just (user_id, course_id)? Check requirements.
);
CREATE INDEX IF NOT EXISTS idx_enrollments_user_id ON enrollments(user_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_course_id ON enrollments(course_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_status ON enrollments(status);
CREATE INDEX IF NOT EXISTS idx_enrollments_deleted_at ON enrollments(deleted_at);
ALTER TABLE enrollments
ADD CONSTRAINT fk_enrollments_user
FOREIGN KEY (user_id) REFERENCES users(id)
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE enrollments
ADD CONSTRAINT fk_enrollments_course
FOREIGN KEY (course_id) REFERENCES courses(id)
ON UPDATE CASCADE ON DELETE CASCADE;

-- Submissions Table
CREATE TABLE IF NOT EXISTS submissions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    lesson_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    enrollment_id BIGINT NOT NULL,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    submission_data JSONB NOT NULL DEFAULT '{}',
    is_late BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_submissions_lesson_id ON submissions(lesson_id);
CREATE INDEX IF NOT EXISTS idx_submissions_user_id ON submissions(user_id);
CREATE INDEX IF NOT EXISTS idx_submissions_enrollment_id ON submissions(enrollment_id);
CREATE INDEX IF NOT EXISTS idx_submissions_deleted_at ON submissions(deleted_at);
-- Optional: Unique constraint (lesson_id, user_id)? Allow resubmission? Check requirements.
-- ALTER TABLE submissions ADD CONSTRAINT uq_submission_lesson_user UNIQUE (lesson_id, user_id);
ALTER TABLE submissions
ADD CONSTRAINT fk_submissions_lesson
FOREIGN KEY (lesson_id) REFERENCES lessons(id)
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE submissions
ADD CONSTRAINT fk_submissions_user
FOREIGN KEY (user_id) REFERENCES users(id)
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE submissions
ADD CONSTRAINT fk_submissions_enrollment
FOREIGN KEY (enrollment_id) REFERENCES enrollments(id)
ON UPDATE CASCADE ON DELETE CASCADE; -- Delete submission if enrollment deleted

-- Grades Table
CREATE TABLE IF NOT EXISTS grades (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    enrollment_id BIGINT NOT NULL,
    lesson_id BIGINT NOT NULL,
    submission_id BIGINT UNIQUE, -- Can be NULL if grading lesson directly, UNIQUE ensures one grade per submission
    grader_id BIGINT NOT NULL, -- User ID of the teacher or system (for AI grading)
    score DECIMAL(10,2) NOT NULL,
    max_score DECIMAL(10,2) NOT NULL,
    feedback TEXT,
    graded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_grades_enrollment_id ON grades(enrollment_id);
CREATE INDEX IF NOT EXISTS idx_grades_lesson_id ON grades(lesson_id);
CREATE INDEX IF NOT EXISTS idx_grades_submission_id ON grades(submission_id);
CREATE INDEX IF NOT EXISTS idx_grades_grader_id ON grades(grader_id);
CREATE INDEX IF NOT EXISTS idx_grades_deleted_at ON grades(deleted_at);
-- Optional: Unique constraint on (enrollment_id, lesson_id) if only one grade per lesson per enrollment allowed?
-- ALTER TABLE grades ADD CONSTRAINT uq_grade_enrollment_lesson UNIQUE (enrollment_id, lesson_id);
ALTER TABLE grades
ADD CONSTRAINT fk_grades_enrollment
FOREIGN KEY (enrollment_id) REFERENCES enrollments(id)
ON UPDATE CASCADE ON DELETE CASCADE; -- Delete grade if enrollment deleted
ALTER TABLE grades
ADD CONSTRAINT fk_grades_lesson
FOREIGN KEY (lesson_id) REFERENCES lessons(id)
ON UPDATE CASCADE ON DELETE CASCADE; -- Delete grade if lesson deleted
ALTER TABLE grades
ADD CONSTRAINT fk_grades_submission
FOREIGN KEY (submission_id) REFERENCES submissions(id)
ON UPDATE CASCADE ON DELETE CASCADE; -- Delete grade if submission deleted (matches Submission model constraint)
ALTER TABLE grades
ADD CONSTRAINT fk_grades_grader
FOREIGN KEY (grader_id) REFERENCES users(id)
ON UPDATE CASCADE ON DELETE RESTRICT; -- Prevent deleting grader user

-- ParentChildLinks Table
CREATE TABLE IF NOT EXISTS parent_child_links (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    parent_user_id BIGINT NOT NULL,
    child_user_id BIGINT NOT NULL,
    UNIQUE (parent_user_id, child_user_id)
);
CREATE INDEX IF NOT EXISTS idx_parent_child_links_parent_id ON parent_child_links(parent_user_id);
CREATE INDEX IF NOT EXISTS idx_parent_child_links_child_id ON parent_child_links(child_user_id);
CREATE INDEX IF NOT EXISTS idx_parent_child_links_deleted_at ON parent_child_links(deleted_at);
ALTER TABLE parent_child_links
ADD CONSTRAINT fk_parent_child_links_parent
FOREIGN KEY (parent_user_id) REFERENCES users(id)
ON UPDATE CASCADE ON DELETE CASCADE; -- Delete link if parent deleted
ALTER TABLE parent_child_links
ADD CONSTRAINT fk_parent_child_links_child
FOREIGN KEY (child_user_id) REFERENCES users(id)
ON UPDATE CASCADE ON DELETE CASCADE; -- Delete link if child deleted

-- PDPAConsentLogs Table (Optional but Recommended)
CREATE TABLE IF NOT EXISTS pdpa_consent_logs (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ, -- Usually not soft deleted, but GORM model has it
    user_id BIGINT NOT NULL,
    consent_type VARCHAR(100) NOT NULL,
    consent_action VARCHAR(20) NOT NULL CHECK (consent_action IN ('GRANTED', 'WITHDRAWN')),
    "timestamp" TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- Use quotes for reserved keyword
    version VARCHAR(50) NOT NULL,
    details TEXT,
    source VARCHAR(100)
);
CREATE INDEX IF NOT EXISTS idx_pdpa_log_user_type ON pdpa_consent_logs(user_id, consent_type);
CREATE INDEX IF NOT EXISTS idx_pdpa_log_action ON pdpa_consent_logs(consent_action);
CREATE INDEX IF NOT EXISTS idx_pdpa_log_source ON pdpa_consent_logs(source);
CREATE INDEX IF NOT EXISTS idx_pdpa_log_deleted_at ON pdpa_consent_logs(deleted_at);
ALTER TABLE pdpa_consent_logs
ADD CONSTRAINT fk_pdpa_consent_logs_user
FOREIGN KEY (user_id) REFERENCES users(id)
ON UPDATE CASCADE ON DELETE CASCADE; -- Delete consent log if user deleted


COMMIT;
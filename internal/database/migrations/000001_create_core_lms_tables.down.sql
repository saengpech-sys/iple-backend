-- Migration: 000001_create_core_lms_tables
-- Created: YYYY-MM-DD HH:MM:SS
-- Reverts the changes from the corresponding .up.sql file.

BEGIN;

-- Drop tables in reverse order of dependency creation, removing FKs first.

-- Drop PDPA Consent Logs (Depends on Users)
ALTER TABLE pdpa_consent_logs DROP CONSTRAINT IF EXISTS fk_pdpa_consent_logs_user;
DROP INDEX IF EXISTS idx_pdpa_log_user_type;
DROP INDEX IF EXISTS idx_pdpa_log_action;
DROP INDEX IF EXISTS idx_pdpa_log_source;
DROP INDEX IF EXISTS idx_pdpa_log_deleted_at;
DROP TABLE IF EXISTS pdpa_consent_logs;

-- Drop Parent Child Links (Depends on Users)
ALTER TABLE parent_child_links DROP CONSTRAINT IF EXISTS fk_parent_child_links_child;
ALTER TABLE parent_child_links DROP CONSTRAINT IF EXISTS fk_parent_child_links_parent;
DROP INDEX IF EXISTS idx_parent_child_links_parent_id;
DROP INDEX IF EXISTS idx_parent_child_links_child_id;
DROP INDEX IF EXISTS idx_parent_child_links_deleted_at;
DROP TABLE IF EXISTS parent_child_links;

-- Drop Grades (Depends on Enrollments, Lessons, Submissions, Users)
ALTER TABLE grades DROP CONSTRAINT IF EXISTS fk_grades_grader;
ALTER TABLE grades DROP CONSTRAINT IF EXISTS fk_grades_submission;
ALTER TABLE grades DROP CONSTRAINT IF EXISTS fk_grades_lesson;
ALTER TABLE grades DROP CONSTRAINT IF EXISTS fk_grades_enrollment;
DROP INDEX IF EXISTS idx_grades_enrollment_id;
DROP INDEX IF EXISTS idx_grades_lesson_id;
DROP INDEX IF EXISTS idx_grades_submission_id;
DROP INDEX IF EXISTS idx_grades_grader_id;
DROP INDEX IF EXISTS idx_grades_deleted_at;
DROP TABLE IF EXISTS grades;

-- Drop Submissions (Depends on Lessons, Users, Enrollments)
ALTER TABLE submissions DROP CONSTRAINT IF EXISTS fk_submissions_enrollment;
ALTER TABLE submissions DROP CONSTRAINT IF EXISTS fk_submissions_user;
ALTER TABLE submissions DROP CONSTRAINT IF EXISTS fk_submissions_lesson;
-- DROP INDEX IF EXISTS uq_submission_lesson_user; -- If unique constraint was added
DROP INDEX IF EXISTS idx_submissions_lesson_id;
DROP INDEX IF EXISTS idx_submissions_user_id;
DROP INDEX IF EXISTS idx_submissions_enrollment_id;
DROP INDEX IF EXISTS idx_submissions_deleted_at;
DROP TABLE IF EXISTS submissions;

-- Drop Enrollments (Depends on Users, Courses)
ALTER TABLE enrollments DROP CONSTRAINT IF EXISTS fk_enrollments_course;
ALTER TABLE enrollments DROP CONSTRAINT IF EXISTS fk_enrollments_user;
DROP INDEX IF EXISTS idx_enrollments_user_id;
DROP INDEX IF EXISTS idx_enrollments_course_id;
DROP INDEX IF EXISTS idx_enrollments_status;
DROP INDEX IF EXISTS idx_enrollments_deleted_at;
-- Drop unique index if exists
-- ALTER TABLE enrollments DROP CONSTRAINT IF EXISTS enrollments_user_id_course_id_key; -- Name might differ based on DB/GORM
DROP TABLE IF EXISTS enrollments;

-- Drop Lessons (Depends on Sections)
ALTER TABLE lessons DROP CONSTRAINT IF EXISTS fk_lessons_section;
-- DROP INDEX IF EXISTS uq_lesson_order; -- If unique constraint was added
DROP INDEX IF EXISTS idx_lessons_section_id;
DROP INDEX IF EXISTS idx_lessons_content_type;
DROP INDEX IF EXISTS idx_lessons_order_index;
DROP INDEX IF EXISTS idx_lessons_is_published;
DROP INDEX IF EXISTS idx_lessons_deleted_at;
DROP TABLE IF EXISTS lessons;

-- Drop Sections (Depends on Courses)
ALTER TABLE sections DROP CONSTRAINT IF EXISTS fk_sections_course;
-- DROP INDEX IF EXISTS uq_section_order; -- If unique constraint was added
DROP INDEX IF EXISTS idx_sections_course_id;
DROP INDEX IF EXISTS idx_sections_order_index;
DROP INDEX IF EXISTS idx_sections_is_published;
DROP INDEX IF EXISTS idx_sections_deleted_at;
DROP TABLE IF EXISTS sections;

-- Drop Courses (Depends on Users)
ALTER TABLE courses DROP CONSTRAINT IF EXISTS fk_courses_teacher;
DROP INDEX IF EXISTS idx_courses_teacher_id;
DROP INDEX IF EXISTS idx_courses_subject;
DROP INDEX IF EXISTS idx_courses_grade_level;
DROP INDEX IF EXISTS idx_courses_is_published;
DROP INDEX IF EXISTS idx_courses_deleted_at;
DROP TABLE IF EXISTS courses;

-- Drop Users (Base table)
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_last_login;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_email;
-- Drop unique constraint if exists on clerk_user_id
-- ALTER TABLE users DROP CONSTRAINT IF EXISTS users_clerk_user_id_key; -- Name might differ
DROP TABLE IF EXISTS users;

COMMIT;
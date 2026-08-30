CREATE TABLE quizzes (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  section_id UUID REFERENCES course_sections(id) ON DELETE SET NULL,
  title VARCHAR(200) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  time_limit_secs INT,
  max_attempts INT NOT NULL DEFAULT 1 CHECK (max_attempts >= 0),
  scheduled_at TIMESTAMPTZ,
  available_until TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_quizzes_course ON quizzes(course_id);

CREATE TABLE quiz_questions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
  type VARCHAR(20) NOT NULL,
  statement TEXT NOT NULL,
  options JSONB NOT NULL DEFAULT '[]'::jsonb,
  correct_answer TEXT NOT NULL,
  points INT NOT NULL DEFAULT 1 CHECK (points > 0),
  position INT NOT NULL DEFAULT 0
);
CREATE INDEX idx_quiz_questions_quiz ON quiz_questions(quiz_id, position);

CREATE TABLE quiz_attempts (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
  user_id VARCHAR(255) NOT NULL,
  attempt_number INT NOT NULL,
  started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  submitted_at TIMESTAMPTZ,
  score INT NOT NULL DEFAULT 0,
  max_score INT NOT NULL DEFAULT 0,
  UNIQUE(quiz_id, user_id, attempt_number)
);
CREATE INDEX idx_quiz_attempts_quiz_user ON quiz_attempts(quiz_id, user_id);

CREATE TABLE quiz_answers (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  attempt_id UUID NOT NULL REFERENCES quiz_attempts(id) ON DELETE CASCADE,
  question_id UUID NOT NULL REFERENCES quiz_questions(id) ON DELETE CASCADE,
  answer_text TEXT NOT NULL DEFAULT '',
  is_correct BOOLEAN NOT NULL DEFAULT FALSE,
  UNIQUE(attempt_id, question_id)
);

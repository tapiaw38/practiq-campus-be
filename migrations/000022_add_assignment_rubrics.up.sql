CREATE TABLE assignment_rubric_criteria (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  assignment_id UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
  title VARCHAR(200) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  max_score INT NOT NULL CHECK (max_score > 0),
  position INT NOT NULL DEFAULT 0
);
CREATE INDEX idx_assignment_rubric_criteria_assignment ON assignment_rubric_criteria(assignment_id, position);

CREATE TABLE submission_rubric_scores (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  submission_id UUID NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
  criterion_id UUID NOT NULL REFERENCES assignment_rubric_criteria(id) ON DELETE CASCADE,
  score INT NOT NULL CHECK (score >= 0),
  feedback TEXT NOT NULL DEFAULT '',
  UNIQUE(submission_id, criterion_id)
);

package domain

type RubricCriterion struct { ID, AssignmentID, Title, Description string; MaxScore, Position int }
type RubricScore struct { ID, SubmissionID, CriterionID, Feedback string; Score int }

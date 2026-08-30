package rubric

import (
	"context"
	"database/sql"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	List(context.Context, string) ([]domain.RubricCriterion, error)
	Replace(context.Context, string, []domain.RubricCriterion) error
	Grade(context.Context, string, int, string, []domain.RubricScore) error
	Scores(context.Context, string) ([]domain.RubricScore, error)
}
type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db} }
func (r *repository) List(c context.Context, a string) (out []domain.RubricCriterion, e error) {
	out = make([]domain.RubricCriterion, 0)
	rows, e := r.db.QueryContext(c, `SELECT id,assignment_id,title,description,max_score,position FROM assignment_rubric_criteria WHERE assignment_id=$1 ORDER BY position`, a)
	if e != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var x domain.RubricCriterion
		e = rows.Scan(&x.ID, &x.AssignmentID, &x.Title, &x.Description, &x.MaxScore, &x.Position)
		if e != nil {
			return
		}
		out = append(out, x)
	}
	e = rows.Err()
	return
}
func (r *repository) Replace(c context.Context, a string, in []domain.RubricCriterion) (e error) {
	tx, e := r.db.BeginTx(c, nil)
	if e != nil {
		return
	}
	defer tx.Rollback()
	if _, e = tx.ExecContext(c, `DELETE FROM assignment_rubric_criteria WHERE assignment_id=$1`, a); e != nil {
		return
	}
	for i, x := range in {
		if _, e = tx.ExecContext(c, `INSERT INTO assignment_rubric_criteria(assignment_id,title,description,max_score,position) VALUES($1,$2,$3,$4,$5)`, a, x.Title, x.Description, x.MaxScore, i); e != nil {
			return
		}
	}
	return tx.Commit()
}
func (r *repository) Grade(c context.Context, s string, total int, feedback string, in []domain.RubricScore) (e error) {
	tx, e := r.db.BeginTx(c, nil)
	if e != nil {
		return
	}
	defer tx.Rollback()
	for _, x := range in {
		_, e = tx.ExecContext(c, `INSERT INTO submission_rubric_scores(submission_id,criterion_id,score,feedback) VALUES($1,$2,$3,$4) ON CONFLICT(submission_id,criterion_id) DO UPDATE SET score=EXCLUDED.score,feedback=EXCLUDED.feedback`, s, x.CriterionID, x.Score, x.Feedback)
		if e != nil {
			return
		}
	}
	_, e = tx.ExecContext(c, `UPDATE submissions SET score=$1,feedback=$2,status='graded',graded_at=NOW() WHERE id=$3`, total, feedback, s)
	if e != nil {
		return
	}
	return tx.Commit()
}
func (r *repository) Scores(c context.Context, s string) (out []domain.RubricScore, e error) {
	rows, e := r.db.QueryContext(c, `SELECT id,submission_id,criterion_id,score,feedback FROM submission_rubric_scores WHERE submission_id=$1`, s)
	if e != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var x domain.RubricScore
		e = rows.Scan(&x.ID, &x.SubmissionID, &x.CriterionID, &x.Score, &x.Feedback)
		if e != nil {
			return
		}
		out = append(out, x)
	}
	e = rows.Err()
	return
}

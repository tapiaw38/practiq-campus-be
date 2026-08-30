package notification

import (
	"context"
	"database/sql"
	"log"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Notification) error
	List(context.Context, string) ([]domain.Notification, error)
	MarkRead(context.Context, string, string) error
}
type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db} }
func (r *repository) Create(c context.Context, n domain.Notification) error {
	_, e := r.db.ExecContext(c, `INSERT INTO notifications(user_id,type,title,body,data) VALUES($1,$2,$3,$4,$5::jsonb)`, n.UserID, n.Type, n.Title, n.Body, n.Data)
	if e != nil {
		log.Printf("[notification] create failed user_id=%s type=%s err=%v", n.UserID, n.Type, e)
	}
	return e
}
func (r *repository) List(c context.Context, u string) (out []domain.Notification, e error) {
	rows, e := r.db.QueryContext(c, `SELECT id,user_id,type,title,body,data::text,read_at,created_at FROM notifications WHERE user_id=$1 ORDER BY created_at DESC LIMIT 100`, u)
	if e != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var n domain.Notification
		if e = rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.Data, &n.ReadAt, &n.CreatedAt); e != nil {
			return
		}
		out = append(out, n)
	}
	e = rows.Err()
	return
}
func (r *repository) MarkRead(c context.Context, id, u string) error {
	_, e := r.db.ExecContext(c, `UPDATE notifications SET read_at=COALESCE(read_at,NOW()) WHERE id=$1 AND user_id=$2`, id, u)
	return e
}

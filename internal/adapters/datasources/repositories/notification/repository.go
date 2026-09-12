package notification

import (
	"context"
	"database/sql"
	"log"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

type Repository interface {
	Create(context.Context, domain.Notification) error
	List(context.Context, string) ([]domain.Notification, error)
	MarkRead(context.Context, string, string) error
}
type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db} }

// Notifications are raised inside the institution whose activity caused them,
// so a person in two of them sees each inbox in its own context.
func (r *repository) Create(c context.Context, n domain.Notification) error {
	tenantID := tenantcontext.ID(c)
	if tenantID == "" {
		return tenantcontext.ErrNoTenant
	}
	_, e := r.db.ExecContext(c,
		`INSERT INTO notifications(tenant_id,user_id,type,title,body,data) VALUES($1,$2,$3,$4,$5,$6::jsonb)`,
		tenantID, n.UserID, n.Type, n.Title, n.Body, n.Data)
	if e != nil {
		log.Printf("[notification] create failed user_id=%s type=%s err=%v", n.UserID, n.Type, e)
	}
	return e
}

func (r *repository) List(c context.Context, u string) (out []domain.Notification, e error) {
	tenantID := tenantcontext.ID(c)
	if tenantID == "" {
		return nil, tenantcontext.ErrNoTenant
	}
	rows, e := r.db.QueryContext(c,
		`SELECT id,user_id,type,title,body,data::text,read_at,created_at FROM notifications
		 WHERE tenant_id=$1 AND user_id=$2 ORDER BY created_at DESC LIMIT 100`, tenantID, u)
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
	tenantID := tenantcontext.ID(c)
	if tenantID == "" {
		return tenantcontext.ErrNoTenant
	}
	_, e := r.db.ExecContext(c,
		`UPDATE notifications SET read_at=COALESCE(read_at,NOW()) WHERE id=$1 AND user_id=$2 AND tenant_id=$3`,
		id, u, tenantID)
	return e
}

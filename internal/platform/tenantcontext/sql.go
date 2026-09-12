package tenantcontext

import (
	"context"
	"strconv"
	"strings"
)

// Query builds a tenant-scoped statement.
//
// Isolation lives in the JOIN rather than in a check the caller runs first.
// A row that does not belong to the tenant is not rejected after being read,
// it is never selected: there is no version of the query that returns it, so
// no caller can forget to ask.
//
// Rows that carry tenant_id themselves use Own. Rows that hang off one — an
// assignment off a course, a message off a conversation — use Through, which
// walks to the ancestor that does.
type Query struct {
	joins  []string
	where  []string
	args   []any
	tenant string
	err    error
}

// NewQuery starts a scoped statement. args are the parameters already bound by
// the caller, so placeholders continue from there.
func NewQuery(ctx context.Context, args ...any) *Query {
	q := &Query{args: args}
	if ID(ctx) == "" {
		q.err = ErrNoTenant
	}
	q.tenant = ID(ctx)
	return q
}

// Own constrains a table that carries tenant_id.
func (q *Query) Own(column string) *Query {
	if q.err != nil {
		return q
	}
	q.where = append(q.where, column+" = "+q.bind(q.tenant))
	return q
}

// Through joins the ancestor that carries tenant_id and constrains it there.
//
//	Through("courses", "c", "a.course_id")
//	→ JOIN courses c ON c.id = a.course_id AND c.tenant_id = $N
func (q *Query) Through(table, alias, foreignKey string) *Query {
	if q.err != nil {
		return q
	}
	q.joins = append(q.joins,
		"JOIN "+table+" "+alias+" ON "+alias+".id = "+foreignKey+
			" AND "+alias+".tenant_id = "+q.bind(q.tenant))
	return q
}

// Join adds a plain join, for walking a chain to the table that holds the
// tenant: a submission reaches courses through its assignment.
func (q *Query) Join(table, alias, on string) *Query {
	if q.err != nil {
		return q
	}
	q.joins = append(q.joins, "JOIN "+table+" "+alias+" ON "+alias+".id = "+on)
	return q
}

// Where adds a condition, binding its values as parameters.
func (q *Query) Where(condition string, values ...any) *Query {
	if q.err != nil {
		return q
	}
	for _, value := range values {
		condition = strings.Replace(condition, "?", q.bind(value), 1)
	}
	q.where = append(q.where, condition)
	return q
}

// SQL assembles the statement around a SELECT and FROM the caller supplies.
func (q *Query) SQL(selectClause, from string) (string, []any, error) {
	if q.err != nil {
		return "", nil, q.err
	}
	statement := "SELECT " + selectClause + " FROM " + from
	if len(q.joins) > 0 {
		statement += " " + strings.Join(q.joins, " ")
	}
	if len(q.where) > 0 {
		statement += " WHERE " + strings.Join(q.where, " AND ")
	}
	return statement, q.args, nil
}

// Args exposes the bound parameters for statements assembled by hand.
func (q *Query) Args() []any { return q.args }

// Err reports a missing tenant, so a caller building a statement in pieces can
// stop before running it.
func (q *Query) Err() error { return q.err }

func (q *Query) bind(value any) string {
	q.args = append(q.args, value)
	return "$" + strconv.Itoa(len(q.args))
}

// ExistsIn builds the tenant guard for an UPDATE or DELETE, which cannot join
// the way a SELECT does.
//
// The guard goes in the same WHERE as the row's own id, so a statement that
// names a row in another institution updates nothing instead of updating it
// and being corrected afterwards.
//
//	ExistsIn(ctx, "courses c", "c.id = assignments.course_id", "c.tenant_id", 11)
//	→ EXISTS (SELECT 1 FROM courses c WHERE c.id = assignments.course_id AND c.tenant_id = $11)
//
// A chain passes the whole path as from:
//
//	ExistsIn(ctx, "assignments a JOIN courses c ON c.id = a.course_id",
//	         "a.id = submissions.assignment_id", "c.tenant_id", 5)
func ExistsIn(ctx context.Context, from, on, tenantColumn string, placeholder int) (string, any, error) {
	id := ID(ctx)
	if id == "" {
		return "", nil, ErrNoTenant
	}
	return "EXISTS (SELECT 1 FROM " + from + " WHERE " + on +
		" AND " + tenantColumn + " = $" + strconv.Itoa(placeholder) + ")", id, nil
}

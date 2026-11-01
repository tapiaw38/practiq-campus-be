package tenantcontext

import (
	"context"
	"strconv"
	"strings"
)

type Query struct {
	joins  []string
	where  []string
	args   []any
	tenant string
	err    error
}

func NewQuery(ctx context.Context, args ...any) *Query {
	q := &Query{args: args}
	if ID(ctx) == "" {
		q.err = ErrNoTenant
	}
	q.tenant = ID(ctx)
	return q
}

func (q *Query) Own(column string) *Query {
	if q.err != nil {
		return q
	}
	q.where = append(q.where, column+" = "+q.bind(q.tenant))
	return q
}

func (q *Query) Through(table, alias, foreignKey string) *Query {
	if q.err != nil {
		return q
	}
	q.joins = append(q.joins,
		"JOIN "+table+" "+alias+" ON "+alias+".id = "+foreignKey+
			" AND "+alias+".tenant_id = "+q.bind(q.tenant))
	return q
}

func (q *Query) Join(table, alias, on string) *Query {
	if q.err != nil {
		return q
	}
	q.joins = append(q.joins, "JOIN "+table+" "+alias+" ON "+alias+".id = "+on)
	return q
}

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

func (q *Query) Args() []any { return q.args }

func (q *Query) Err() error { return q.err }

func (q *Query) bind(value any) string {
	q.args = append(q.args, value)
	return "$" + strconv.Itoa(len(q.args))
}

func ExistsIn(ctx context.Context, from, on, tenantColumn string, placeholder int) (string, any, error) {
	id := ID(ctx)
	if id == "" {
		return "", nil, ErrNoTenant
	}
	return "EXISTS (SELECT 1 FROM " + from + " WHERE " + on +
		" AND " + tenantColumn + " = $" + strconv.Itoa(placeholder) + ")", id, nil
}

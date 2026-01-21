package query

import (
	"fmt"
	"strings"
)

// QueryBuilder helps build dynamic SQL queries
type QueryBuilder struct {
	baseQuery  string
	conditions []string
	orderBy    string
	limit      int
	offset     int
	args       []any
	argIndex   int
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder(baseQuery string) *QueryBuilder {
	return &QueryBuilder{
		baseQuery:  baseQuery,
		conditions: []string{},
		args:       []any{},
		argIndex:   1,
	}
}

// AddCondition adds a WHERE condition
func (qb *QueryBuilder) AddCondition(condition string, value any) *QueryBuilder {
	qb.conditions = append(qb.conditions, fmt.Sprintf(condition, qb.argIndex))
	qb.args = append(qb.args, value)
	qb.argIndex++
	return qb
}

// AddLikeCondition adds a LIKE condition for search
func (qb *QueryBuilder) AddLikeCondition(columns []string, searchValue string) *QueryBuilder {
	if searchValue == "" {
		return qb
	}

	var likeConds []string
	searchPattern := "%" + searchValue + "%"

	for _, col := range columns {
		likeConds = append(likeConds, fmt.Sprintf("%s ILIKE $%d", col, qb.argIndex))
		qb.args = append(qb.args, searchPattern)
		qb.argIndex++
	}

	if len(likeConds) > 0 {
		qb.conditions = append(qb.conditions, "("+strings.Join(likeConds, " OR ")+")")
	}

	return qb
}

// AddOrderBy adds ORDER BY clause
func (qb *QueryBuilder) AddOrderBy(orderBy string) *QueryBuilder {
	qb.orderBy = orderBy
	return qb
}

// AddLimit adds LIMIT clause
func (qb *QueryBuilder) AddLimit(limit int) *QueryBuilder {
	if limit > 0 {
		qb.limit = limit
	}
	return qb
}

// AddOffset adds OFFSET clause
func (qb *QueryBuilder) AddOffset(offset int) *QueryBuilder {
	if offset > 0 {
		qb.offset = offset
	}
	return qb
}

// Build returns the final query and arguments
func (qb *QueryBuilder) Build() (string, []any) {
	query := qb.baseQuery

	// Add WHERE clause
	if len(qb.conditions) > 0 {
		query += " WHERE " + strings.Join(qb.conditions, " AND ")
	}

	// Add ORDER BY clause
	if qb.orderBy != "" {
		query += " ORDER BY " + qb.orderBy
	}

	// Add LIMIT clause
	if qb.limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", qb.argIndex)
		qb.args = append(qb.args, qb.limit)
		qb.argIndex++
	}

	// Add OFFSET clause
	if qb.offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", qb.argIndex)
		qb.args = append(qb.args, qb.offset)
		qb.argIndex++
	}

	return query, qb.args
}

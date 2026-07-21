package query

import (
	"fmt"
	"strings"
)

type Builder struct {
	BaseQuery string
	Clauses   []string
	Args      []any
}

type Filterable interface {
	ToQuery() Builder
}

func (b *Builder) Equal(column string, val any) {
	if val != nil {
		b.Args = append(b.Args, val)
		b.Clauses = append(b.Clauses, fmt.Sprintf("%s = $%d", column, len(b.Args)))
	}
}

func (b *Builder) Like(column string, val any) {
	if val != nil {
		b.Args = append(b.Args, val)
		b.Clauses = append(b.Clauses, fmt.Sprintf("%s LIKE $%d", column, len(b.Args)))
	}
}

func (b *Builder) Raw(subquery string, arg any) {
	const onlyFirstOccurrence = 1
	b.Args = append(b.Args, arg)
	placeholder := fmt.Sprintf("$%d", len(b.Args))
	finalClause := strings.Replace(subquery, "?", placeholder, onlyFirstOccurrence)
	b.Clauses = append(b.Clauses, finalClause)
}

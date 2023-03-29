package clauses

import (
	"gorm.io/gorm/clause"
)

type SelectUnion struct {
	Columns []clause.Column
	Values  [][]interface{}
}

// Name from clause name
func (SelectUnion) Name() string {
	return "SELECT_UNION"
}

// Build build from clause
func (values SelectUnion) Build(builder clause.Builder) {
	if len(values.Columns) > 0 {
		builder.WriteByte('(')
		for idx, column := range values.Columns {
			if idx > 0 {
				builder.WriteByte(',')
			}
			builder.WriteQuoted(column)
		}
		builder.WriteByte(')')

		for idx, value := range values.Values {
			builder.WriteString(" SELECT ")
			builder.AddVar(builder, value...)
			builder.WriteString(" FROM DUAL ")
			if idx < len(values.Values)-1 {
				builder.WriteString(" UNION ")
			}
		}
	}
}

// MergeClause merge values clauses
func (values SelectUnion) MergeClause(clause *clause.Clause) {
	clause.Name = ""
	clause.Expression = values
}

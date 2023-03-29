package altibase

import (
	"github.com/bulenttokuzlu/altibase/clauses"
	"gorm.io/gorm/clause"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

func Create(db *gorm.DB) {
	stmt := db.Statement
	schema := stmt.Schema

	if stmt == nil || schema == nil {
		return
	}

	values := callbacks.ConvertToCreateValues(stmt)
	stmt.AddClauseIfNotExists(clause.Insert{Table: clause.Table{Name: stmt.Table}})
	stmt.AddClauseIfNotExists(clauses.SelectUnion{Columns: values.Columns, Values: values.Values})
	//	stmt.AddClause(clause.Values{Columns: values.Columns, Values: [][]interface{}{values.Values[0]}})
	stmt.Build("INSERT", "SELECT_UNION")

	if !db.DryRun {
		for idx, vals := range values.Values {
			// HACK HACK: replace values one by one, assuming its value layout will be the same all the time, i.e. aligned
			for idx, val := range vals {
				switch v := val.(type) {
				case bool:
					if v {
						val = 1
					} else {
						val = 0
					}
				}
				stmt.Vars[idx] = val
			}
			switch result, err := stmt.ConnPool.ExecContext(stmt.Context, stmt.SQL.String(), stmt.Vars...); err {
			case nil: // success
				db.RowsAffected, _ = result.RowsAffected()

				insertTo := stmt.ReflectValue
				switch insertTo.Kind() {
				case reflect.Slice, reflect.Array:
					insertTo = insertTo.Index(idx)
				}
			default: // failure
				db.AddError(err)
			}
		}

	}
}

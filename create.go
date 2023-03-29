package altibase

import (
	"fmt"
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
	fmt.Println("stmt.Build - stmt.SQL.String() - ", stmt.SQL.String())

	if !db.DryRun {
		fmt.Println("len(values.Values) = ", len(values.Values))
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
			fmt.Println("idx = ", idx)
			//fmt.Println("stmt.SQL.String() = ", stmt.SQL.String())
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

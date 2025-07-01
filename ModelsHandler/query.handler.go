package models_handler

import (
	"fmt"
	"strings"

	"github.com/vrianta/Server/DatabaseHandler"
)

func (m *Struct) Get() *Query {
	return &Query{
		model: m,
	}
}

func (q *Query) Where(column string) *Query {
	q.lastColumn = column
	return q
}

func (q *Query) Is(value any) *Query {
	q.conditions = append(q.conditions, fmt.Sprintf("`%s` = ?", q.lastColumn))
	q.args = append(q.args, value)
	q.lastColumn = ""
	return q
}

func (q *Query) Like(value string) *Query {
	q.conditions = append(q.conditions, fmt.Sprintf("`%s` LIKE ?", q.lastColumn))
	q.args = append(q.args, value)
	q.lastColumn = ""
	return q
}

func (q *Query) And() *Query {
	q.conditions = append(q.conditions, "AND")
	return q
}

func (q *Query) Or() *Query {
	q.conditions = append(q.conditions, "OR")
	return q
}

func (q *Query) Limit(n int) *Query {
	q.limit = n
	return q
}

func (q *Query) Fetch() ([]*Struct, error) {
	db, err := DatabaseHandler.GetDatabase()
	if err != nil {
		return nil, err
	}

	whereClause := ""
	if len(q.conditions) > 0 {
		whereClause = "WHERE " + strings.Join(q.conditions, " ")
	}

	limitClause := ""
	if q.limit > 0 {
		limitClause = fmt.Sprintf("LIMIT %d", q.limit)
	}

	query := fmt.Sprintf("SELECT * FROM %s %s %s", q.model.TableName, whereClause, limitClause)

	rows, err := db.Query(query, q.args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []*Struct

	for rows.Next() {
		fieldPointers := make([]any, len(columns))
		valueHolders := make([]any, len(columns))

		for i := range columns {
			fieldPointers[i] = &valueHolders[i]
		}

		if err := rows.Scan(fieldPointers...); err != nil {
			return nil, err
		}

		rowModel := &Struct{
			TableName: q.model.TableName,
			fields:    make(map[string]Field),
		}

		for k, v := range q.model.fields {
			rowModel.fields[k] = v
		}

		for i, col := range columns {
			val := valueHolders[i]
			if b, ok := val.([]byte); ok {
				val = string(b)
			}

			if field, ok := rowModel.fields[col]; ok {
				field.value = val
				rowModel.fields[col] = field
			}
		}

		results = append(results, rowModel)
	}

	return results, rows.Err()
}

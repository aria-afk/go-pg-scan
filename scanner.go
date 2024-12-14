package scanner

import (
	"database/sql"
	"reflect"
)

// target is expected to be of type []*struct/[]*primitive
func Query(conn *sql.DB, target any, query string, args ...any) error {
	// NOTE: we need to ensure dest type is a []*FOO type
	dest := reflect.ValueOf(target).Elem()
	destT := dest.Type().Elem().Elem()

	rows, err := conn.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		cols, err := rows.Columns()
		if err != nil {
			return err
		}
		colTypes, err := rows.ColumnTypes()
		if err != nil {
			return err
		}

		structP := reflect.New(destT)
		structVal := structP.Elem()

		for i := 0; i < structVal.NumField(); i++ {
			name := cols[i]
			field := structVal.FieldByName(name)
			/* TODO: field not found (should be a 0 value if not found)
			   if !name {}
			*/
			// TODO: colType should align with field.Type()
		}
	}

	return nil
}

/*
Notes on how to do mapping.

TLDR; iterate via value.NumField() and use colName to set props

func main() {
	var res []*Foo
	test(&res)
}

func test(target any) {
	slice := reflect.ValueOf(target).Elem()

	valuet := slice.Type().Elem().Elem()
	valuep := reflect.New(valuet)
	value := valuep.Elem()

	for i := 0; i < value.NumField(); i++ {
		name := value.Type().Field(i).Name
		field := value.FieldByName(name)
		fieldType := field.Type()
		fmt.Println(fieldType)
	}
}
*/

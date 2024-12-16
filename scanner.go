package scanner

import (
	"database/sql"
	"fmt"
	"reflect"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TODO: Struct field name problem:
// Exported values like Username should be captured when SQL/pSQL returns username
// ie; we need some kind of smart normalizing of this

// TODO: Nested struct
// TODO: Error handling and type confirmation
// TODO: (Maybe) context

func Query(conn *sql.DB, target any, query string, args ...any) error {
	dest := reflect.ValueOf(target).Elem()
	destT := dest.Type().Elem().Elem()

	rows, err := conn.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	colLength := len(cols)
	collectedVals := make([]interface{}, colLength)
	for i := range collectedVals {
		var placeholder interface{}
		collectedVals[i] = &placeholder
	}

	for rows.Next() {
		err := rows.Scan(collectedVals...)
		if err != nil {
			return err
		}
		structP := reflect.New(destT)
		structVal := structP.Elem()

		for i, colName := range cols {
			fmt.Println(FirstToUpper(colName))
			collectedValue := *(collectedVals[i]).(*interface{})
			valueType := reflect.TypeOf(collectedValue).Kind()

			structField := structVal.FieldByName(FirstToUpper(colName))

			switch valueType {
			case reflect.String:
				structField.SetString(collectedValue.(string))
			}
		}

		dest.Set(reflect.Append(dest, structP))
	}

	return nil
}

func FirstToUpper(str string) string {
	return cases.Title(language.English).String(str)
}

package bootstrap

import (
	"reflect"

	"github.com/DATA-DOG/godog/gherkin"
)

// SetMember sets a value of a struct member
func SetMember(s interface{}, member string, value string) {
	v := reflect.ValueOf(s).Elem().FieldByName(member)

	if v.IsValid() {
		v.SetString(value)
	}
}

// SerializeTableRow maps values of a gherkin table row to a struct
func SerializeTableRow(s interface{}, columns *gherkin.TableRow, values *gherkin.TableRow) {
	members := columns.Cells

	for index, member := range members {
		SetMember(s, member.Value, values.Cells[index].Value)
	}
}

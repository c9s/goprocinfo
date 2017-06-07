package linux

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSnmp(t *testing.T) {
	{
		var expected = Snmp{2, 64, 15379636127, 0, 216025, 0, 0, 0, 15378356252, 18361945377, 850, 90, 72, 100, 0, 72, 0, 0, 0, 160938, 192, 0, 12113, 63, 0, 0, 0, 148754, 0, 3, 0, 3, 0, 358528, 0, 209771, 0, 0, 0, 0, 0, 148754, 0, 3, 0, 0, 0, 12113, 0, 148754, 63, 3, 148754, 209771, 0, 3, 0, 1, 200, 120000, 0, 31508152, 20211455, 35049, 3917068, 41, 14849418023, 18216233673, 3100617, 336, 24740359, 13, 526064282, 220285, 2329341, 528804507, 2329341, 0, 0, 0, 0, 0, 0, 0, 0, 0}

		read, err := ReadSnmp("proc/snmp")
		if err != nil {
			t.Fatal("snmp read fail", err)
		}

		t.Logf("%+v", expected)
		t.Logf("%+v", read)

		if err := compareExpectedReadFieldsSnmp(&expected, read); err != nil {
			t.Error(err.Error())
		}

		if !reflect.DeepEqual(*read, expected) {
			t.Error("not equal to expected")
		}
	}
}

func compareExpectedReadFieldsSnmp(expected *Snmp, read *Snmp) error {
	elemExpected := reflect.ValueOf(*expected)
	typeOfElemExpected := elemExpected.Type()
	elemRead := reflect.ValueOf(*read)

	for i := 0; i < elemExpected.NumField(); i++ {
		fieldName := typeOfElemExpected.Field(i).Name

		if elemExpected.Field(i).Uint() != elemRead.Field(i).Uint() {
			return fmt.Errorf("Read value not equal to expected value for field %s. Got %d and expected %d.", fieldName, elemRead.Field(i).Uint(), elemExpected.Field(i).Uint())
		}
	}

	return nil
}

package testing

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// CompareUnexported allows cmp to compare private fields of internal structs
var CompareUnexported = cmp.Exporter(func(reflect.Type) bool { return true })

// TestHelper methods exist to provide handle mundane tasks which can result in error without clogging test files with unimportant error checking code
type TestHelper struct {
	t *testing.T
}

func NewTestHelper(t *testing.T) TestHelper {
	return TestHelper{t}
}

func (h TestHelper) ReadFile(path string) []byte {
	h.t.Helper()

	data, err := ioutil.ReadFile(path)
	if err != nil {
		h.t.Fatalf("Error reading file\n %+v", err)
		return nil
	}

	return data
}

func (h TestHelper) Unmarshal(data []byte, destination interface{}) {
	h.t.Helper()

	err := json.Unmarshal(data, destination)
	if err != nil {
		h.t.Fatalf("Error unmarshaling json %q \n", err)
	}
}

package testing

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
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

func (h TestHelper) XmlUnmarshal(data []byte, destination interface{}) {
	h.t.Helper()

	err := xml.Unmarshal(data, destination)
	if err != nil {
		h.t.Fatalf("Error unmarshaling json %q \n", err)
	}
}

// RoundTrip allows TestHelper to act as a HTTP RoundTripper that logs requests and responses.
// This can be used by overriding the HTTP client used by a PSP client to be the TestHelper instance.
//
// Example:
//
//	 helper := sleet_testing.NewTestHelper(t)
//	 httpClient := &http.Client{
//		 Transport: helper,
//		 Timeout:   common.DefaultTimeout,
//	 }
func (h TestHelper) RoundTrip(req *http.Request) (*http.Response, error) {
	h.t.Helper()

	resp, err := http.DefaultTransport.RoundTrip(req)

	reqBodyStream, _ := req.GetBody()
	defer reqBodyStream.Close()
	reqBody, _ := ioutil.ReadAll(reqBodyStream)

	respBodyStream := resp.Body
	defer respBodyStream.Close()
	respBody, _ := ioutil.ReadAll(respBodyStream)
	// we need to replace the resp body to be read again by the actual handler
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))

	h.t.Logf(
		"logTransport HTTP request\n"+
			"-> status %s\n"+
			"-v request\n"+
			string(reqBody)+"\n"+
			"-v response\n"+
			string(respBody)+"\n\n",
		resp.Status)
	return resp, err
}

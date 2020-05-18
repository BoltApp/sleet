package stripe_test

import (
	"fmt"
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/gateways/stripe"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

// TestSuccessfulAuth is a unit test mocking Stripe calls
// This was used for us to ensure that with successful responses, authorization performed correctly
// Integration tests hitting Stripe's service can be found in the integration-tests folder
func TestSuccessfulAuth(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodPost, "https://api.stripe.com/v1/tokens",
		httpmock.NewBytesResponder(200, readJson(t, "tokens_success.json")))
	httpmock.RegisterResponder(http.MethodPost, "https://api.stripe.com/v1/charges",
		httpmock.NewBytesResponder(200, readJson(t, "charges_success.json")))
	client := stripe.NewWithHTTPClient("apiKey", http.DefaultClient)
	authRequest := sleet_testing.BaseAuthorizationRequest()

	response, err := client.Authorize(authRequest)

	assert.Nil(t, err)
	expectedResponse := &sleet.AuthorizationResponse{
		Success:              true,
		TransactionReference: "ch_1FfpIZFSEDlaFyqYGbP2DpkI",
		AvsResult:            sleet.AVSresponseZipMatchAddressMatch, // TODO: Add translator
		CvvResult:            sleet.CVVResponseMatch,                // TODO: Add translator
		AvsResultRaw:         "unchecked",
		CvvResultRaw:         "unchecked",
		ErrorCode:            "200",
	}
	assert.Equal(t, expectedResponse, response)
}

func readJson(t *testing.T, fileName string) []byte {
	jsonFile, err := os.Open(fmt.Sprintf("testdata/%s", fileName))
	assert.Nil(t, err)
	bytes, err := ioutil.ReadAll(jsonFile)
	assert.Nil(t, err)
	return bytes
}

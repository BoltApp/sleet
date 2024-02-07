package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/checkoutcom"
	sleet_testing "github.com/BoltApp/sleet/testing"
)

type ClientNamePair struct {
	client *checkoutcom.CheckoutComClient
	name   string
}

func generateClients() []ClientNamePair {
	fmt.Println(getEnv("CHECKOUTCOM_TEST_KEY"))
	fmt.Println(getEnv("CHECKOUTCOM_TEST_KEY_WITH_PCID"))
	fmt.Println(getEnv("CHECKOUTCOM_TEST_PCID"))
	legacyClient := checkoutcom.NewClient(common.Sandbox, getEnv("CHECKOUTCOM_TEST_KEY"), nil)
	pcidClient := checkoutcom.NewClient(common.Sandbox, getEnv("CHECKOUTCOM_TEST_KEY_WITH_PCID"), common.SPtr(getEnv("CHECKOUTCOM_TEST_PCID")))

	clients := []ClientNamePair{
		{
			client: pcidClient,
			name:   "PCID Client",
		},
		{
			client: legacyClient,
			name:   "Legacy Client",
		},
	}
	return clients
}

// TestCheckoutComAuthorizeFailed
//
// checkout.com has test cards here: https://www.checkout.com/docs/four/testing/response-code-testing
// Using a rejected card number
func TestCheckoutComAuthorizeFailed(t *testing.T) {
	clients := generateClients()

	for _, clientNamePair := range clients {
		client := clientNamePair.client
		name := clientNamePair.name
		failedRequest := sleet_testing.BaseAuthorizationRequest()
		failedRequest.CreditCard.Number = "4544249167673670"
		response, err := client.Authorize(failedRequest)

		if err != nil {
			t.Errorf("%s: Authorize request should not have an error even if authorization failed- %s", name, err.Error())
		}

		if response.Success {
			t.Errorf("%s: Auth response should indicate a failure", name)
		}

		if response.Response != "20051" {
			t.Errorf("%s: Response should be 20051, code for insufficient funds- %s", name, response.Response)
		}
	}
}

// TestCheckoutComAuth
//
// This should successfully create an authorization
func TestCheckoutComAuth(t *testing.T) {
	clients := generateClients()

	for _, clientNamePair := range clients {
		client := clientNamePair.client
		name := clientNamePair.name
		request := sleet_testing.BaseAuthorizationRequest()
		auth, err := client.Authorize(request)
		if err != nil {
			t.Errorf("%s: Authorize request should not have failed", name)
		}

		if !auth.Success {
			t.Errorf("%s: Resulting auth should have been successful", name)
		}
	}
}

// TestCheckoutComAuthFullCapture
//
// This should successfully create an authorization on checkout.com then Capture for full amount
func TestCheckoutComAuthFullCapture(t *testing.T) {
	clients := generateClients()

	for _, clientNamePair := range clients {
		client := clientNamePair.client
		name := clientNamePair.name
		authRequest := sleet_testing.BaseAuthorizationRequest()
		auth, err := client.Authorize(authRequest)
		if err != nil {
			t.Errorf("%s: Authorize request should not have failed", name)
		}

		if !auth.Success {
			t.Errorf("%s: Resulting auth should have been successful", name)
		}

		captureRequest := &sleet.CaptureRequest{
			Amount:                     &authRequest.Amount,
			TransactionReference:       auth.TransactionReference,
			ClientTransactionReference: authRequest.ClientTransactionReference,
			MerchantOrderReference:     &authRequest.MerchantOrderReference,
		}
		capture, err := client.Capture(captureRequest)
		if err != nil {
			t.Errorf("%s: Capture request should not have failed", name)
		}

		if !capture.Success {
			t.Errorf("%s: Resulting capture should have been successful", name)
		}
	}
}

// TestCheckoutComAuthFullCapture
//
// This should successfully create an authorization on checkout.com then Capture for full amount
func TestCheckoutComAuthPartialCapture(t *testing.T) {
	clients := generateClients()

	for _, clientNamePair := range clients {
		client := clientNamePair.client
		name := clientNamePair.name
		authRequest := sleet_testing.BaseAuthorizationRequest()
		authRequest.Amount.Amount = 100
		auth, err := client.Authorize(authRequest)
		if err != nil {
			t.Errorf("%s: Authorize request should not have failed", name)
		}

		if !auth.Success {
			t.Errorf("%s: Resulting auth should have been successful", name)
		}

		// Partial capture request
		captureRequest := &sleet.CaptureRequest{
			Amount: &sleet.Amount{
				Amount:   50,
				Currency: "USD",
			},
			TransactionReference:       auth.TransactionReference,
			ClientTransactionReference: authRequest.ClientTransactionReference,
			MerchantOrderReference:     &authRequest.MerchantOrderReference,
		}
		capture, err := client.Capture(captureRequest)
		if err != nil {
			t.Errorf("%s: Partial capture request should not have failed", name)
		}

		if !capture.Success {
			t.Errorf("%s: Resulting partial capture should have been successful", name)
		}

		// Capture the rest
		captureRequest = &sleet.CaptureRequest{
			Amount: &sleet.Amount{
				Amount:   50,
				Currency: "USD",
			},
			TransactionReference:       auth.TransactionReference,
			ClientTransactionReference: authRequest.ClientTransactionReference,
			MerchantOrderReference:     &authRequest.MerchantOrderReference,
		}
		capture, err = client.Capture(captureRequest)
		if err != nil {
			t.Errorf("%s: Final capture request should not have failed", name)
		}

		if !capture.Success {
			t.Errorf("%s: Resulting final capture should have been successful", name)
		}
	}
}

// TestCheckoutComAuthVoid
//
// This should successfully create an authorization on checkout.com then Void/Cancel the Auth
func TestCheckoutComAuthVoid(t *testing.T) {
	clients := generateClients()

	for _, clientNamePair := range clients {
		client := clientNamePair.client
		name := clientNamePair.name
		authRequest := sleet_testing.BaseAuthorizationRequest()
		auth, err := client.Authorize(authRequest)
		if err != nil {
			t.Errorf("%s: Authorize request should not have failed", name)
		}

		if !auth.Success {
			t.Errorf("%s: Resulting auth should have been successful", name)
		}

		voidRequest := &sleet.VoidRequest{
			TransactionReference:   auth.TransactionReference,
			MerchantOrderReference: &authRequest.MerchantOrderReference,
		}
		void, err := client.Void(voidRequest)
		if err != nil {
			t.Errorf("%s: Void request should not have failed", name)
		}

		if !void.Success {
			t.Errorf("%s: Resulting void should have been successful", name)
		}
	}
}

// TestCheckoutComAuthCaptureRefund
//
// This should successfully create an authorization on checkout.com, then Capture for full amount, then refund for full amount
func TestCheckoutComAuthCaptureRefund(t *testing.T) {
	clients := generateClients()

	for _, clientNamePair := range clients {
		client := clientNamePair.client
		name := clientNamePair.name
		authRequest := sleet_testing.BaseAuthorizationRequest()
		auth, err := client.Authorize(authRequest)
		if err != nil {
			t.Errorf("%s: Authorize request should not have failed", name)
		}

		if !auth.Success {
			t.Errorf("%s: Resulting auth should have been successful", name)
		}

		captureRequest := &sleet.CaptureRequest{
			Amount:                     &authRequest.Amount,
			TransactionReference:       auth.TransactionReference,
			ClientTransactionReference: authRequest.ClientTransactionReference,
			MerchantOrderReference:     &authRequest.MerchantOrderReference,
		}
		capture, err := client.Capture(captureRequest)
		if err != nil {
			t.Errorf("%s: Capture request should not have failed", name)
		}

		if !capture.Success {
			t.Errorf("%s: Resulting capture should have been successful", name)
		}

		time.Sleep(4 * time.Second) // Delay to make sure capture has processed

		refundRequest := &sleet.RefundRequest{
			Amount:                     &authRequest.Amount,
			TransactionReference:       capture.TransactionReference,
			ClientTransactionReference: authRequest.ClientTransactionReference,
			MerchantOrderReference:     &authRequest.MerchantOrderReference,
		}

		refund, err := client.Refund(refundRequest)
		if err != nil {
			t.Errorf("%s: Refund request should not have failed", name)
		}

		if !refund.Success {
			t.Errorf("%s: Resulting refund should have been successful", name)
		}
	}
}

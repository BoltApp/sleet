package stripe

import (
	"github.com/BoltApp/sleet"
	"strconv"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest) (*TokenRequest, error) {
	request := &TokenRequest{
		ExpYear:  strconv.Itoa(authRequest.CreditCard.ExpirationYear),
		ExpMonth: strconv.Itoa(authRequest.CreditCard.ExpirationMonth),
		Number:   authRequest.CreditCard.Number,
		CVC:      authRequest.CreditCard.CVV,
		Name:     authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName,
	}

	if authRequest.BillingAddress.StreetAddress1 != nil {
		request.AddressLine1 = *authRequest.BillingAddress.StreetAddress1
	}
	if authRequest.BillingAddress.StreetAddress2 != nil {
		request.AddressLine2 = *authRequest.BillingAddress.StreetAddress2
	}
	if authRequest.BillingAddress.Locality != nil {
		request.AddressCity = *authRequest.BillingAddress.Locality
	}
	if authRequest.BillingAddress.RegionCode != nil {
		request.AddressState = *authRequest.BillingAddress.RegionCode
	}
	if authRequest.BillingAddress.CountryCode != nil {
		request.AddressCountry = *authRequest.BillingAddress.CountryCode
	}
	if authRequest.BillingAddress.PostalCode != nil {
		request.AddressZip = *authRequest.BillingAddress.PostalCode
	}

	return request, nil
}

func buildChargeRequest(authRequest *sleet.AuthorizationRequest, chargeID string) (*ChargeRequest, error) {
	request := &ChargeRequest{
		Amount:   strconv.FormatInt(authRequest.Amount.Amount, 10),
		Currency: authRequest.Amount.Currency,
		Source:   chargeID,
		Capture:  "false",
	}
	return request, nil
}

func buildRefundRequest(authRequest *sleet.RefundRequest) (*PostAuthRequest, error) {
	request := &PostAuthRequest{
		Amount: strconv.FormatInt(authRequest.Amount.Amount, 10),
		Charge: authRequest.TransactionReference,
	}
	return request, nil
}

func buildCaptureRequest(authRequest *sleet.CaptureRequest) (*PostAuthRequest, error) {
	request := &PostAuthRequest{
		Amount: strconv.FormatInt(authRequest.Amount.Amount, 10),
	}
	return request, nil
}

func buildVoidRequest(authRequest *sleet.VoidRequest) (*PostAuthRequest, error) {
	request := &PostAuthRequest{
		Charge: authRequest.TransactionReference,
	}
	return request, nil
}

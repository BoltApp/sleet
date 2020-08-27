package orbital

import (
	"encoding/xml"
	"strconv"

	"github.com/BoltApp/sleet"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest, credentials Credentials) Request {

	amount := authRequest.Amount.Amount
	exp := strconv.Itoa(authRequest.CreditCard.ExpirationYear) + strconv.Itoa(authRequest.CreditCard.ExpirationMonth)
	code := currencyMap[authRequest.Amount.Currency]

	body := RequestBody{
		IndustryType:     IndustryTypeEcomm,
		MessageType:      MessageTypeAuth,
		BIN:              BINStratus,
		TerminalID:       TerminalIDStratus,
		AccountNum:       authRequest.CreditCard.Number,
		Exp:              exp,
		CurrencyCode:     code,
		CurrencyExponent: CurrencyExponentDefault,
		CardSecVal:       authRequest.CreditCard.CVV,
		OrderID:          *authRequest.ClientTransactionReference,
		Amount:           amount,
		AVSzip:           *authRequest.BillingAddress.PostalCode,
		AVSaddress1:      *authRequest.BillingAddress.StreetAddress1,
		AVSaddress2:      authRequest.BillingAddress.StreetAddress2,
		AVSstate:         *authRequest.BillingAddress.RegionCode,
		AVScity:          *authRequest.BillingAddress.Locality,
		AVScountryCode:   *authRequest.BillingAddress.CountryCode,
	}

	if authRequest.CreditCard.Network == sleet.CreditCardNetworkVisa || authRequest.CreditCard.Network == sleet.CreditCardNetworkDiscover {
		body.CardSecValInd = CardSecPresent
	}

	body.XMLName = xml.Name{Local: RequestTypeNewOrder}

	body.OrbitalConnectionUsername = credentials.username
	body.OrbitalConnectionPassword = credentials.password
	body.MerchantID = credentials.merchantID
	return Request{Body: body}
}

func buildCaptureRequest(captureRequest *sleet.CaptureRequest, credentials Credentials) Request {
	body := RequestBody{
		Amount:     captureRequest.Amount.Amount,
		BIN:        BINStratus,
		TerminalID: TerminalIDStratus,
		TxRefNum:   captureRequest.TransactionReference,
		OrderID:    *captureRequest.ClientTransactionReference,
	}

	body.XMLName = xml.Name{Local: RequestTypeCapture}
	body.OrbitalConnectionUsername = credentials.username
	body.OrbitalConnectionPassword = credentials.password
	body.MerchantID = credentials.merchantID
	return Request{Body: body}
}

func buildVoidRequest(voidRequest *sleet.VoidRequest, credentials Credentials) Request {
	body := RequestBody{
		BIN:        BINStratus,
		TerminalID: TerminalIDStratus,
		TxRefNum:   voidRequest.TransactionReference,
		OrderID:    *voidRequest.ClientTransactionReference,
	}

	body.XMLName = xml.Name{Local: RequestTypeVoid}
	body.OrbitalConnectionUsername = credentials.username
	body.OrbitalConnectionPassword = credentials.password
	body.MerchantID = credentials.merchantID
	return Request{Body: body}
}

func buildRefundRequest(refundRequest *sleet.RefundRequest, credentials Credentials) Request {
	amount := refundRequest.Amount.Amount
	code := currencyMap[refundRequest.Amount.Currency]

	body := RequestBody{
		IndustryType:     IndustryTypeEcomm,
		MessageType:      MessageTypeRefund,
		BIN:              BINStratus,
		TerminalID:       TerminalIDStratus,
		CurrencyCode:     code,
		CurrencyExponent: CurrencyExponentDefault,
		OrderID:          *refundRequest.ClientTransactionReference,
		Amount:           amount,
		TxRefNum:         refundRequest.TransactionReference,
	}

	body.XMLName = xml.Name{Local: RequestTypeNewOrder}
	body.OrbitalConnectionUsername = credentials.username
	body.OrbitalConnectionPassword = credentials.password
	body.MerchantID = credentials.merchantID
	return Request{Body: body}
}

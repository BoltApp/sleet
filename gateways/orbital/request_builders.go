package orbital

import (
	"encoding/xml"
	"strconv"

	"github.com/BoltApp/sleet"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest, creds Credentials) Request {

	amount := authRequest.Amount.Amount
	exp := strconv.Itoa(authRequest.CreditCard.ExpirationYear) + strconv.Itoa(authRequest.CreditCard.ExpirationMonth)
	code := currencyMap[authRequest.Amount.Currency]

	body := RequestBody{
		OrbitalConnectionUsername: creds.Username,
		OrbitalConnectionPassword: creds.Password,
		IndustryType:              IndustryTypeEcomm,
		MessageType:               MessageTypeAuth,
		BIN:                       BINStratus,
		MerchantID:                creds.MerchantID,
		TerminalID:                TerminalIDStratus,
		CardBrand:                 "VI",
		AccountNum:                authRequest.CreditCard.Number,
		Exp:                       exp,
		CurrencyCode:              code,
		CurrencyExponent:          CurrencyExponentDefault,
		CardSecVal:                authRequest.CreditCard.CVV,
		OrderID:                   *authRequest.ClientTransactionReference,
		Amount:                    amount,
		AVSzip:                    *authRequest.BillingAddress.PostalCode,
		AVSaddress1:               *authRequest.BillingAddress.StreetAddress1,
		AVSaddress2:               authRequest.BillingAddress.StreetAddress2,
		AVScity:                   *authRequest.BillingAddress.Locality,
		AVSstate:                  *authRequest.BillingAddress.RegionCode,
		AVScountryCode:            *authRequest.BillingAddress.CountryCode,
	}

	if authRequest.CreditCard.Network == sleet.CreditCardNetworkVisa || authRequest.CreditCard.Network == sleet.CreditCardNetworkDiscover {
		body.CardSecValInd = CardSecPresent
	}

	body.XMLName = xml.Name{Local: RequestTypeNewOrder}
	return Request{Body: body}
}

func buildCaptureRequest(captureRequest *sleet.CaptureRequest) Request {
	body := RequestBody{
		Amount:     captureRequest.Amount.Amount,
		BIN:        BINStratus,
		TerminalID: TerminalIDStratus,
		TxRefNum:   captureRequest.TransactionReference,
		OrderID:    *captureRequest.ClientTransactionReference,
	}

	body.XMLName = xml.Name{Local: RequestTypeCapture}
	return Request{Body: body}
}

func buildVoidRequest(voidRequest *sleet.VoidRequest) Request {
	body := RequestBody{
		BIN:        BINStratus,
		TerminalID: TerminalIDStratus,
		TxRefNum:   voidRequest.TransactionReference,
		OrderID:    *voidRequest.ClientTransactionReference,
	}

	body.XMLName = xml.Name{Local: RequestTypeVoid}
	return Request{Body: body}
}

func buildRefundRequest(refundRequest *sleet.RefundRequest) Request {
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
	return Request{Body: body}
}

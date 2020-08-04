package orbital

import (
	"github.com/BoltApp/sleet"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest) Request {
	// TODO NOTE why is amount in sleet an int, already in no decimals format ?
	// amount := &authRequest.Amount.Amount * 100 //convert to no decimals
	amount := authRequest.Amount.Amount //convert to no decimals
	exp := authRequest.CreditCard.ExpirationYear + authRequest.CreditCard.ExpirationMonth
	// code := currencyCode(authRequest.Amount.Currency)
	code := 840 //TODO should create map of all values ?

	body := RequestBody{
		IndustryType:     IndustryTypeEcomm,
		MessageType:      MessageTypeAuth,
		BIN:              BINStratus,        // create consts for this
		TerminalID:       TerminalIDStratus, // create consts for this
		AccountNum:       authRequest.CreditCard.Number,
		Exp:              exp,
		CurrencyCode:     code,
		CurrencyExponent: "2",                                        //should we just deafault 2 or create a map for those with 0 ?
		AVSzip:           *authRequest.BillingAddress.PostalCode,     // create consts for this
		AVSaddress1:      *authRequest.BillingAddress.StreetAddress1, // create consts for this
		AVSaddress2:      *authRequest.BillingAddress.StreetAddress2, // create consts for this
		CardSecValInd:    CardSecPresent,                             // create consts for this
		CardSecVal:       authRequest.CreditCard.CVV,
		OrderID:          *authRequest.ClientTransactionReference,
		Amount:           amount,
		// AVSstate
		// AVScity
		// AVSname
		// AVScountryCode
	}
	return Request{Body: body}
}

func buildCaptureRequest(captureRequest *sleet.CaptureRequest, ref string) Request {
	body := RequestBody{
		Amount:     captureRequest.Amount.Amount,
		TerminalID: "001",    // create consts for this
		BIN:        "000001", // create consts for this
		TxRefNum:   ref,
		OrderID:    *captureRequest.ClientTransactionReference,
	}
	return Request{Body: body}
}

func buildVoidRequest(voidRequest *sleet.VoidRequest) Request {
	body := RequestBody{
		TerminalID: "001",    // create consts for this
		BIN:        "000001", // create consts for this
		// TxRefNum:    ref, TODO NOTE OrderID vs txrefNum ??
		OrderID: *voidRequest.ClientTransactionReference,
	}
	return Request{Body: body}
}

func buildRefundRequest(refundRequest *sleet.RefundRequest) Request {
	body := RequestBody{
		AdjustedAmt: refundRequest.Amount.Amount,
		TerminalID:  "001",    // create consts for this
		BIN:         "000001", // create consts for this
		// TxRefNum:    ref, TODO NOTE OrderID vs txrefNum ??
		OrderID: *refundRequest.ClientTransactionReference,
	}
	return Request{Body: body}
}

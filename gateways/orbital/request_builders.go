package orbital

import (
	"encoding/xml"
	"strconv"

	"github.com/BoltApp/sleet"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest) Request {
	// TODO NOTE why is amount in sleet an int, already in no decimals format ?
	// amount := &authRequest.Amount.Amount * 100 //convert to no decimals
	amount := authRequest.Amount.Amount //convert to no decimals
	exp := strconv.Itoa(authRequest.CreditCard.ExpirationYear) + strconv.Itoa(authRequest.CreditCard.ExpirationMonth)
	// code := currencyCode(authRequest.Amount.Currency)
	code := currencyMap[authRequest.Amount.Currency]

	body := RequestBody{
		IndustryType:     IndustryTypeEcomm,
		MessageType:      MessageTypeAuth,
		BIN:              BINStratus,
		TerminalID:       TerminalIDStratus,
		AccountNum:       authRequest.CreditCard.Number,
		Exp:              exp,
		CurrencyCode:     code,
		CurrencyExponent: "2", //should we just deafault 2 or create a map for those with 0 ?
		AVSzip:           *authRequest.BillingAddress.PostalCode,
		AVSaddress1:      *authRequest.BillingAddress.StreetAddress1,
		CardSecValInd:    CardSecPresent,
		CardSecVal:       authRequest.CreditCard.CVV,
		OrderID:          *authRequest.ClientTransactionReference,
		Amount:           amount,
		//AVSaddress2:      *authRequest.BillingAddress.StreetAddress2, // create consts for this
		// AVSstate
		// AVScity
		// AVSname
		// AVScountryCode
	}

	body.XMLName = xml.Name{Local: RequestTypeAuth}
	return Request{Body: body}
}

func buildCaptureRequest(captureRequest *sleet.CaptureRequest) Request {

	code := currencyMap[captureRequest.Amount.Currency]

	body := RequestBody{
		Amount:       captureRequest.Amount.Amount,
		CurrencyCode: code,
		TerminalID:   "001",    // create consts for this
		BIN:          "000001", // create consts for this
		TxRefNum:     captureRequest.TransactionReference,
		OrderID:      *captureRequest.ClientTransactionReference,
	}

	body.XMLName = xml.Name{Local: RequestTypeCapture}
	return Request{Body: body}
}

func buildVoidRequest(voidRequest *sleet.VoidRequest) Request {
	body := RequestBody{
		TerminalID: "001",    // create consts for this
		BIN:        "000001", // create consts for this
		TxRefNum:   voidRequest.TransactionReference,
		OrderID:    *voidRequest.ClientTransactionReference,
	}

	body.XMLName = xml.Name{Local: RequestTypeVoid}
	return Request{Body: body}
}

func buildRefundRequest(refundRequest *sleet.RefundRequest) Request {

	code := currencyMap[refundRequest.Amount.Currency]
	body := RequestBody{
		CurrencyCode: code, // should this be here ?
		AdjustedAmt:  refundRequest.Amount.Amount,
		TerminalID:   "001",    // create consts for this
		BIN:          "000001", // create consts for this
		TxRefNum:     refundRequest.TransactionReference,
		OrderID:      *refundRequest.ClientTransactionReference,
	}

	body.XMLName = xml.Name{Local: RequestTypeRefund}
	return Request{Body: body}
}

package rocketgate

import (
	"strconv"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/rocketgate-sdk/request"
)

// Cof specifies the transaction type under the Credential-on-File framework
const (
	cofCIT = "CIT"	// Customer Initiated Transaction
	cofMIT = "MIT"	// Merchant Initiated Transaction
)

var initiatorTypeToCofType = map[sleet.ProcessingInitiatorType]string{
	sleet.ProcessingInitiatorTypeInitialCardOnFile:         cofCIT,
	sleet.ProcessingInitiatorTypeInitialRecurring:          cofCIT,
	sleet.ProcessingInitiatorTypeStoredCardholderInitiated: cofCIT,
	sleet.ProcessingInitiatorTypeStoredMerchantInitiated:   cofMIT,
	sleet.ProcessingInitiatorTypeFollowingRecurring:        cofMIT,
}

func buildAuthRequest(merchantID string, merchantPassword string, merchantAccount *string, authRequest *sleet.AuthorizationRequest) *request.GatewayRequest {
	card := authRequest.CreditCard

	gatewayRequest := request.NewGatewayRequest()

	gatewayRequest.Set(request.MERCHANT_ID, merchantID)
	gatewayRequest.Set(request.MERCHANT_PASSWORD, merchantPassword)
	if merchantAccount != nil {
		gatewayRequest.Set(request.MERCHANT_ACCOUNT, *merchantAccount)
	}

	// Credit Card
	gatewayRequest.Set(request.CARDNO, card.Number)
	gatewayRequest.Set(request.EXPIRE_MONTH, strconv.Itoa(card.ExpirationMonth))
	gatewayRequest.Set(request.EXPIRE_YEAR, strconv.Itoa(card.ExpirationYear))
	gatewayRequest.Set(request.AMOUNT, sleet.AmountToDecimalString(&authRequest.Amount))
	gatewayRequest.Set(request.CURRENCY, authRequest.Amount.Currency)

	// overwrites the flag transactions
	if authRequest.ProcessingInitiator != nil {
		if cofType, ok := initiatorTypeToCofType[*authRequest.ProcessingInitiator]; ok {
			gatewayRequest.Set(request.COF_FRAMEWORK, cofType)
		}
	}

	// Ignore CVV and AVS check
	gatewayRequest.Set(request.CVV2_CHECK, "NO")
	gatewayRequest.Set(request.AVS_CHECK, "IGNORE")

	return gatewayRequest
}

func buildCaptureRequest(merchantID string, merchantPassword string, captureRequest *sleet.CaptureRequest) *request.GatewayRequest {
	gatewayRequest := request.NewGatewayRequest()

	gatewayRequest.Set(request.MERCHANT_ID, merchantID)
	gatewayRequest.Set(request.MERCHANT_PASSWORD, merchantPassword)
	gatewayRequest.Set(request.TRANSACT_ID, captureRequest.TransactionReference)

	// Optional if the amount is the same as the original purchase or auth-only transaction.
	gatewayRequest.Set(request.AMOUNT, strconv.Itoa(int(captureRequest.Amount.Amount)))
	gatewayRequest.Set(request.CURRENCY, captureRequest.Amount.Currency)

	return gatewayRequest
}

func buildVoidRequest(merchantID string, merchantPassword string, voidRequest *sleet.VoidRequest) *request.GatewayRequest {
	gatewayRequest := request.NewGatewayRequest()

	gatewayRequest.Set(request.MERCHANT_ID, merchantID)
	gatewayRequest.Set(request.MERCHANT_PASSWORD, merchantPassword)
	gatewayRequest.Set(request.TRANSACT_ID, voidRequest.TransactionReference)

	return gatewayRequest
}

func buildRefundRequest(merchantID string, merchantPassword string, refundRequest *sleet.RefundRequest) *request.GatewayRequest {
	gatewayRequest := request.NewGatewayRequest()

	gatewayRequest.Set(request.MERCHANT_ID, merchantID)
	gatewayRequest.Set(request.MERCHANT_PASSWORD, merchantPassword)
	gatewayRequest.Set(request.TRANSACT_ID, refundRequest.TransactionReference)

	// Optional if the amount is the same as the original purchase or auth-only transaction.
	gatewayRequest.Set(request.AMOUNT, strconv.Itoa(int(refundRequest.Amount.Amount)))
	gatewayRequest.Set(request.CURRENCY, refundRequest.Amount.Currency)

	return gatewayRequest
}

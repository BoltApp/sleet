package authorizenet

import (
	"errors"
	"fmt"

	"github.com/BoltApp/sleet"
)

func buildAuthRequest(merchantName string, transactionKey string, authRequest *sleet.AuthorizationRequest) *Request{
	amountStr := sleet.AmountToDecimalString(&authRequest.Amount)
	billingAddress := authRequest.BillingAddress
	authorizeRequest := CreateTransactionRequest{
		MerchantAuthentication: authentication(merchantName, transactionKey),
		TransactionRequest: TransactionRequest{
			TransactionType: TransactionTypeAuthOnly,
			Amount:          &amountStr,
			Payment: &Payment{
				CreditCard: CreditCard{
					CardNumber:     authRequest.CreditCard.Number,
					ExpirationDate: fmt.Sprintf("%d-%d", authRequest.CreditCard.ExpirationYear, authRequest.CreditCard.ExpirationMonth),
					CardCode:       &authRequest.CreditCard.CVV,
				},
			},
			BillingAddress: &BillingAddress{
				FirstName: authRequest.CreditCard.FirstName,
				LastName:  authRequest.CreditCard.LastName,
				Address:   billingAddress.StreetAddress1,
				City:      billingAddress.Locality,
				State:     billingAddress.RegionCode,
				Zip:       billingAddress.PostalCode,
				Country:   billingAddress.CountryCode,
			},
		},
	}
	return &Request{CreateTransactionRequest: authorizeRequest}
}

func buildVoidRequest(merchantName string, transactionKey string, voidRequest *sleet.VoidRequest) *Request {
	return &Request{
		CreateTransactionRequest: CreateTransactionRequest{
			MerchantAuthentication: authentication(merchantName, transactionKey),
			TransactionRequest: TransactionRequest{
				TransactionType:  TransactionTypeVoid,
				RefTransactionID: &voidRequest.TransactionReference,
			},
		},
	}
}

func buildCaptureRequest(merchantName string, transactionKey string, captureRequest *sleet.CaptureRequest) *Request {
	amountStr := sleet.AmountToDecimalString(captureRequest.Amount)
	request := &Request{
		CreateTransactionRequest: CreateTransactionRequest{
			MerchantAuthentication: authentication(merchantName, transactionKey),
			TransactionRequest: TransactionRequest{
				TransactionType:  TransactionTypePriorAuthCapture,
				Amount:           &amountStr,
				RefTransactionID: &captureRequest.TransactionReference,
			},
		},
	}
	return request
}

func buildRefundRequest(merchantName string, transactionKey string, refundRequest *sleet.RefundRequest) (*Request, error) {
	lastFour, ok := refundRequest.Options["credit_card"]
	if !ok {
		return nil, errors.New("missing credit card last four digits")
	}
	lastFourAsString := lastFour.(string)
	if len(lastFourAsString) != 4 {
		return nil, errors.New("incorrect credit card last four digits")
	}
	amountStr := sleet.AmountToDecimalString(refundRequest.Amount)
	request := &Request{
		CreateTransactionRequest: CreateTransactionRequest{
			MerchantAuthentication: authentication(merchantName, transactionKey),
			TransactionRequest: TransactionRequest{
				TransactionType:  TransactionTypeRefund,
				Amount:           &amountStr,
				RefTransactionID: &refundRequest.TransactionReference,
				Payment: &Payment{
					CreditCard: CreditCard{
						CardNumber:     lastFourAsString,
						ExpirationDate: expirationDateXXXX,
					},
				},
			},
		},
	}
	return request, nil
}

func authentication(merchantName string, transactionKey string) MerchantAuthentication {
	return MerchantAuthentication{
		Name:           merchantName,
		TransactionKey: transactionKey,
	}
}

package authorize_net

import (
	"errors"
	"fmt"
	"github.com/BoltApp/sleet"
)

func buildAuthRequest(merchantName string, transactionKey string, authRequest *sleet.AuthorizationRequest) (*Request, error) {
	amountStr := fmt.Sprintf("%.2f", float64(authRequest.Amount.Amount)/100.0)

	billingAddress := authRequest.BillingAddress
	authorizeRequest := CreateTransactionRequest{
		MerchantAuthentication: authentication(merchantName, transactionKey),
		TransactionRequest: TransactionRequest{
			TransactionType: transactionTypeAuthOnly,
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
	request := Request{CreateTransactionRequest: authorizeRequest}
	return &request, nil
}

func buildVoidRequest(merchantName string, transactionKey string, voidRequest *sleet.VoidRequest) (*Request, error) {
	request := &Request{
		CreateTransactionRequest: CreateTransactionRequest{
			MerchantAuthentication: authentication(merchantName, transactionKey),
			TransactionRequest: TransactionRequest{
				TransactionType:  transactionTypeVoid,
				RefTransactionID: &voidRequest.TransactionReference,
			},
		},
	}
	return request, nil
}

func buildCaptureRequest(merchantName string, transactionKey string, captureRequest *sleet.CaptureRequest) (*Request, error) {
	amount := fmt.Sprintf("%.2f", float64(captureRequest.Amount.Amount)/100.0)
	request := &Request{
		CreateTransactionRequest: CreateTransactionRequest{
			MerchantAuthentication: authentication(merchantName, transactionKey),
			TransactionRequest: TransactionRequest{
				TransactionType:  transactionTypePriorAuthCapture,
				Amount:           &amount,
				RefTransactionID: &captureRequest.TransactionReference,
			},
		},
	}
	return request, nil
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
	amount := fmt.Sprintf("%.2f", float64(refundRequest.Amount.Amount) / 100.0)
	request := &Request{
		CreateTransactionRequest: CreateTransactionRequest{
			MerchantAuthentication: authentication(merchantName, transactionKey),
			TransactionRequest: TransactionRequest{
				TransactionType:  transactionTypeRefund,
				Amount:           &amount,
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

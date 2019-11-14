package authorize_net

import (
	"fmt"
	"github.com/BoltApp/sleet"
)

func buildAuthRequest(merchantName string, transactionKey string, authRequest *sleet.AuthorizationRequest) (*Request, error) {
	amountStr := sleet.AmountToCentsString(authRequest.Amount.Amount)
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
					CardCode:       authRequest.CreditCard.CVV,
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
	amountStr := sleet.AmountToCentsString(captureRequest.Amount.Amount)
	request := &Request{
		CreateTransactionRequest: CreateTransactionRequest{
			MerchantAuthentication: authentication(merchantName, transactionKey),
			TransactionRequest: TransactionRequest{
				TransactionType:  transactionTypePriorAuthCapture,
				Amount:           &amountStr,
				RefTransactionID: &captureRequest.TransactionReference,
			},
		},
	}
	return request, nil
}

func buildRefundRequest(merchantName string, transactionKey string, refundRequest *sleet.RefundRequest) (*Request, error) {
	amountStr := sleet.AmountToCentsString(refundRequest.Amount.Amount)
	request := &Request{
		CreateTransactionRequest: CreateTransactionRequest{
			MerchantAuthentication: authentication(merchantName, transactionKey),
			TransactionRequest: TransactionRequest{
				TransactionType:  transactionTypeRefund,
				Amount:           &amountStr,
				RefTransactionID: &refundRequest.TransactionReference,
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

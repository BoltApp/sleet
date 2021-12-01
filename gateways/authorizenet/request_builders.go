package authorizenet

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
)

const (
	InvoiceNumberMaxLength = 20
)

func buildAuthRequest(merchantName string, transactionKey string, authRequest *sleet.AuthorizationRequest) *Request {
	amountStr := sleet.AmountToDecimalString(&authRequest.Amount)
	billingAddress := authRequest.BillingAddress

	creditCard := CreditCard{
		CardNumber:     authRequest.CreditCard.Number,
		ExpirationDate: fmt.Sprintf("%d-%d", authRequest.CreditCard.ExpirationYear, authRequest.CreditCard.ExpirationMonth),
	}
	if authRequest.Cryptogram != "" {
		// Apple Pay request
		creditCard.IsPaymentToken = common.BPtr(true)
		creditCard.Cryptogram = authRequest.Cryptogram
	} else {
		// Credit Card request
		creditCard.CardCode = authRequest.CreditCard.CVV
	}

	authorizeRequest := CreateTransactionRequest{
		MerchantAuthentication: authentication(merchantName, transactionKey),
		TransactionRequest: TransactionRequest{
			TransactionType: TransactionTypeAuthOnly,
			Amount:          &amountStr,
			Payment: &Payment{
				CreditCard: creditCard,
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

	if authRequest.MerchantOrderReference != "" {
		invoiceNumber := sleet.TruncateString(authRequest.MerchantOrderReference, InvoiceNumberMaxLength)
		authorizeRequest.TransactionRequest.Order = &Order{InvoiceNumber: invoiceNumber}
	}

	authorizeRequest = *addL2L3Data(authRequest, &authorizeRequest)

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
						CardNumber:     refundRequest.Last4,
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

func addL2L3Data(authRequest *sleet.AuthorizationRequest, authNetAuthRequest *CreateTransactionRequest) *CreateTransactionRequest {
	if authRequest.Level3Data != nil {
		authNetAuthRequest.TransactionRequest.LineItem = json.RawMessage(buildLineItemsString(authRequest))

		authNetAuthRequest.TransactionRequest.Tax = &Tax{
			Amount:  strconv.FormatInt(authRequest.Level3Data.TaxAmount.Amount, 10),
		}

		authNetAuthRequest.TransactionRequest.Duty = &Tax{
			Amount:  strconv.FormatInt(authRequest.Level3Data.DutyAmount.Amount, 10),
		}

		authNetAuthRequest.TransactionRequest.Shipping = &Tax{
			Amount:  strconv.FormatInt(authRequest.Level3Data.ShippingAmount.Amount, 10),
		}

		authNetAuthRequest.TransactionRequest.Customer = &Customer{
			Id: authRequest.Level3Data.CustomerReference,
		}
	}

	if authRequest.ShippingAddress != nil {
		authNetAuthRequest.TransactionRequest.ShippingAddress = &BillingAddress{
			FirstName: 	authRequest.CreditCard.FirstName,
			LastName: 	authRequest.CreditCard.LastName,
			Company: 	common.SafeStr(authRequest.ShippingAddress.Company),
			Address: 	authRequest.ShippingAddress.StreetAddress1,
			City:      	authRequest.ShippingAddress.Locality,
			State:    	authRequest.ShippingAddress.RegionCode,
			Zip:       	authRequest.ShippingAddress.PostalCode,
			Country:   	authRequest.ShippingAddress.CountryCode,
		}
	}

	return authNetAuthRequest
}

func buildLineItemsString(authRequest *sleet.AuthorizationRequest) string {
	// No more than 30 lineitems
	lineItemCount := len(authRequest.Level3Data.LineItems)
	if lineItemCount > 30 {
		lineItemCount = 30
	}

	lineItems := "{"
	for i := 0; i < lineItemCount; i++ {
		lineItem := &LineItem{
			ItemId: sleet.TruncateString(authRequest.Level3Data.LineItems[i].CommodityCode, 31),
			Name: sleet.TruncateString(authRequest.Level3Data.LineItems[i].ProductCode, 31),
			Description: sleet.TruncateString(authRequest.Level3Data.LineItems[i].Description, 255),
			Quantity: strconv.FormatInt(authRequest.Level3Data.LineItems[i].Quantity, 10),
			UnitPrice: strconv.FormatInt(authRequest.Level3Data.LineItems[i].UnitPrice.Amount, 10),
		}

		// This is absolutely terrible why is json -> xml a thing
		lineItemByte, err := json.Marshal(lineItem)
		if err == nil {
			// No error, add the string
			// We must use a string because authnet allows invalid JSON
			lineItems += "\"lineItem\":" + string(lineItemByte)
			if i < lineItemCount - 1 {
				lineItems += ","
			}
		}
	}
	lineItems += "}"

	if len(lineItems) > 2 {
		return lineItems
	}
	return ""
}
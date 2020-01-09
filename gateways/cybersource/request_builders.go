package cybersource

import (
	"strconv"

	. "github.com/BoltApp/sleet/common"

	"github.com/BoltApp/sleet"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest) (*Request, error) {
	amountStr := sleet.AmountToString(&authRequest.Amount)
	request := &Request{
		ProcessingInformation: &ProcessingInformation{
			Capture:           false, // no autocapture for now
			CommerceIndicator: "internet",
		},
		PaymentInformation: &PaymentInformation{
			Card: CardInformation{
				ExpYear:  strconv.Itoa(authRequest.CreditCard.ExpirationYear),
				ExpMonth: strconv.Itoa(authRequest.CreditCard.ExpirationMonth),
				Number:   authRequest.CreditCard.Number,
				CVV:      authRequest.CreditCard.CVV,
			},
		},
		OrderInformation: &OrderInformation{
			AmountDetails: AmountDetails{
				Amount:   amountStr,
				Currency: authRequest.Amount.Currency,
			},
			BillTo: BillingInformation{
				FirstName:  authRequest.CreditCard.FirstName,
				LastName:   authRequest.CreditCard.LastName,
				Address1:   *authRequest.BillingAddress.StreetAddress1,
				Address2:   SafeStr(authRequest.BillingAddress.StreetAddress2),
				PostalCode: *authRequest.BillingAddress.PostalCode,
				Locality:   *authRequest.BillingAddress.Locality,
				AdminArea:  *authRequest.BillingAddress.RegionCode,
				Country:    SafeStr(authRequest.BillingAddress.CountryCode),
				Email:      SafeStr(authRequest.BillingAddress.Email),
				Company:    SafeStr(authRequest.BillingAddress.Company),
			},
		},
	}
	if authRequest.ClientTransactionReference != nil {
		request.ClientReferenceInformation = &ClientReferenceInformation{
			Code: *authRequest.ClientTransactionReference,
		}
	}
	if authRequest.Level3Data != nil {
		level3 := authRequest.Level3Data
		request.ProcessingInformation.PurchaseLevel = "3" // Signify that this request contains level 3 data
		if request.ClientReferenceInformation == nil {
			request.ClientReferenceInformation = &ClientReferenceInformation{}
		}
		request.ClientReferenceInformation.Code = level3.CustomerReference

		request.OrderInformation.ShipTo = ShippingDetails{
			PostalCode: level3.DestinationPostalCode,
			Country:    level3.DestinationCountryCode,
			// TODO: Add administrative area
		}
		request.OrderInformation.AmountDetails.DiscountAmount = strconv.FormatInt(level3.DiscountAmount, 10)
		request.OrderInformation.AmountDetails.TaxAmount = strconv.FormatInt(level3.TaxAmount, 10)
		request.OrderInformation.AmountDetails.FreightAmount = strconv.FormatInt(level3.ShippingAmount, 10)
		request.OrderInformation.AmountDetails.DutyAmount = "0" // TODO: Add DutyAmount to level3 data
		for _, lineItem := range level3.LineItems {
			request.OrderInformation.LineItems = append(request.OrderInformation.LineItems, LineItem{
				ProductCode:    lineItem.ProductCode,
				ProductName:    lineItem.Description, // TODO: Check if this is correct, add ProductName to level3 data?
				Quantity:       strconv.FormatInt(lineItem.Quantity, 10),
				UnitPrice:      strconv.FormatInt(lineItem.UnitPrice, 10),
				TotalAmount:    strconv.FormatInt(lineItem.TotalAmount, 10),
				DiscountAmount: strconv.FormatInt(lineItem.ItemDiscountAmount, 10),
				UnitOfMeasure:  lineItem.UnitOfMeasure,
				CommodityCode:  lineItem.CommodityCode,
				TaxAmount:      strconv.FormatInt(lineItem.ItemTaxAmount, 10),
			})
		}
	}
	return request, nil
}

func buildCaptureRequest(captureRequest *sleet.CaptureRequest) (*Request, error) {
	amountStr := sleet.AmountToString(captureRequest.Amount)
	request := &Request{
		OrderInformation: &OrderInformation{
			AmountDetails: AmountDetails{
				Amount:   amountStr,
				Currency: captureRequest.Amount.Currency,
			},
		},
	}
	return request, nil
}

func buildVoidRequest(voidRequest *sleet.VoidRequest) (*Request, error) {
	// Maybe add reason / more details, but for now nothing
	request := &Request{}
	return request, nil
}

func buildRefundRequest(refundRequest *sleet.RefundRequest) (*Request, error) {
	amountStr := sleet.AmountToString(refundRequest.Amount)
	request := &Request{
		OrderInformation: &OrderInformation{
			AmountDetails: AmountDetails{
				Amount:   amountStr,
				Currency: refundRequest.Amount.Currency,
			},
		},
	}
	return request, nil
}

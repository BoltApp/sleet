package cybersource

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"strconv"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest) (*Request, error) {
	amountStr := sleet.AmountToString(&authRequest.Amount)
	request := &Request{
		ClientReferenceInformation: &ClientReferenceInformation{
			Partner: Partner{
				SolutionID: authRequest.Channel,
			},
		},
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
				Address2:   common.SafeStr(authRequest.BillingAddress.StreetAddress2),
				PostalCode: *authRequest.BillingAddress.PostalCode,
				Locality:   *authRequest.BillingAddress.Locality,
				AdminArea:  *authRequest.BillingAddress.RegionCode,
				Country:    common.SafeStr(authRequest.BillingAddress.CountryCode),
				Email:      common.SafeStr(authRequest.BillingAddress.Email),
				Company:    common.SafeStr(authRequest.BillingAddress.Company),
			},
		},
	}
	// If level 3 data is present, and ClientReferenceInformation in that data exists, it will override this.
	if authRequest.ClientTransactionReference != nil {
		request.ClientReferenceInformation.Code = *authRequest.ClientTransactionReference
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
			AdminArea:  level3.DestinationAdminArea,
		}
		request.OrderInformation.AmountDetails.DiscountAmount = sleet.AmountToString(&level3.DiscountAmount)
		request.OrderInformation.AmountDetails.TaxAmount = sleet.AmountToString(&level3.TaxAmount)
		request.OrderInformation.AmountDetails.FreightAmount = sleet.AmountToString(&level3.ShippingAmount)
		request.OrderInformation.AmountDetails.DutyAmount = sleet.AmountToString(&level3.DutyAmount)
		for _, lineItem := range level3.LineItems {
			request.OrderInformation.LineItems = append(request.OrderInformation.LineItems, LineItem{
				ProductCode:    lineItem.ProductCode,
				ProductName:    lineItem.Description,
				Quantity:       strconv.FormatInt(lineItem.Quantity, 10),
				UnitPrice:      sleet.AmountToString(&lineItem.UnitPrice),
				TotalAmount:    sleet.AmountToString(&lineItem.TotalAmount),
				DiscountAmount: sleet.AmountToString(&lineItem.ItemDiscountAmount),
				UnitOfMeasure:  lineItem.UnitOfMeasure,
				CommodityCode:  lineItem.CommodityCode,
				TaxAmount:      sleet.AmountToString(&lineItem.ItemTaxAmount),
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
	if captureRequest.ClientTransactionReference != nil {
		request.ClientReferenceInformation = &ClientReferenceInformation{
			Code: *captureRequest.ClientTransactionReference,
		}
	}
	return request, nil
}

func buildVoidRequest(voidRequest *sleet.VoidRequest) (*Request, error) {
	// Maybe add reason / more details, but for now nothing
	request := &Request{}
	if voidRequest.ClientTransactionReference != nil {
		request.ClientReferenceInformation = &ClientReferenceInformation{
			Code: *voidRequest.ClientTransactionReference,
		}
	}
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
	if refundRequest.ClientTransactionReference != nil {
		request.ClientReferenceInformation = &ClientReferenceInformation{
			Code: *refundRequest.ClientTransactionReference,
		}
	}
	return request, nil
}

package adyen

import (
	"fmt"
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/adyen/adyen-go-api-library/v2/src/checkout"
	"github.com/adyen/adyen-go-api-library/v2/src/payments"
	"strconv"
)

const (
	maxLineItemDescriptionLength = 26
	level3Default                = "NA"
)

func buildAuthRequest(authRequest *sleet.AuthorizationRequest, merchantAccount string) *checkout.PaymentRequest {
	request := &checkout.PaymentRequest{
		Amount: checkout.Amount{
			Value:    authRequest.Amount.Amount,
			Currency: authRequest.Amount.Currency,
		},
		// Adyen requires a reference in request so this will panic if client doesn't pass it. Assuming this is good for now
		Reference: *authRequest.ClientTransactionReference,
		PaymentMethod: map[string]interface{}{
			"number":      authRequest.CreditCard.Number,
			"expiryMonth": strconv.Itoa(authRequest.CreditCard.ExpirationMonth),
			"expiryYear":  strconv.Itoa(authRequest.CreditCard.ExpirationYear),
			"holderName":  authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName,
		},
		MerchantAccount: merchantAccount,
	}

	if authRequest.BillingAddress != nil {
		request.BillingAddress = &checkout.Address{
			City:              common.SafeStr(authRequest.BillingAddress.Locality),
			Country:           common.SafeStr(authRequest.BillingAddress.CountryCode),
			HouseNumberOrName: common.SafeStr(authRequest.BillingAddress.StreetAddress2),
			PostalCode:        common.SafeStr(authRequest.BillingAddress.PostalCode),
			StateOrProvince:   common.SafeStr(authRequest.BillingAddress.RegionCode),
			Street:            common.SafeStr(authRequest.BillingAddress.StreetAddress1),
		}
	}

	addPaymentSpecificFields(authRequest, request)

	level3 := authRequest.Level3Data
	if level3 != nil {
		request.AdditionalData = buildLevel3Data(level3)
	}

	return request
}

// addPaymentSpecificFields adds fields to the Adyen Payment request that are dependent on the payment method
func addPaymentSpecificFields(authRequest *sleet.AuthorizationRequest, request *checkout.PaymentRequest) {
	if authRequest.Cryptogram != "" && authRequest.ECI != "" {
		// Apple Pay request
		request.MpiData = &checkout.ThreeDSecureData{
			AuthenticationResponse: "Y",
			Cavv:                   authRequest.Cryptogram,
			DirectoryResponse:      "Y",
			Eci:                    authRequest.ECI,
		}
		request.PaymentMethod["brand"] = authRequest.CreditCard.Network.String()
		request.PaymentMethod["type"] = "networkToken"
		request.RecurringProcessingModel = "CardOnFile"
		request.ShopperInteraction = "Ecommerce"
	} else if authRequest.CreditCard.CVV != "" {
		// New customer credit card request
		request.PaymentMethod["cvc"] = authRequest.CreditCard.CVV
		request.PaymentMethod["type"] = "scheme"
		request.ShopperInteraction = "Ecommerce"
		if authRequest.CreditCard.Save {
			// Customer opts in to saving card details
			request.RecurringProcessingModel = "CardOnFile"
			request.StorePaymentMethod = true
		} else {
			// Customer opts out of saving card details
			request.StorePaymentMethod = false
		}
	} else {
		// Existing customer credit card request
		request.PaymentMethod["type"] = "scheme"
		request.RecurringProcessingModel = "CardOnFile"
		request.ShopperInteraction = "ContAuth"
	}
}

func buildLevel3Data(level3Data *sleet.Level3Data) map[string]string {
	additionalData := map[string]string{
		"enhancedSchemeData.customerReference":     sleet.DefaultIfEmpty(level3Data.CustomerReference, level3Default),
		"enhancedSchemeData.destinationPostalCode": level3Data.DestinationPostalCode,
		"enhancedSchemeData.dutyAmount":            sleet.AmountToString(&level3Data.DutyAmount),
		"enhancedSchemeData.freightAmount":         sleet.AmountToString(&level3Data.ShippingAmount),
		"enhancedSchemeData.totalTaxAmount":        sleet.AmountToString(&level3Data.TaxAmount),
	}

	var keyBase string
	for idx, lineItem := range level3Data.LineItems {
		// Maximum of 9 line items allowed in the request
		if idx == 9 {
			break
		}
		keyBase = fmt.Sprintf("enhancedSchemeData.itemDetailLine%d.", idx+1)
		additionalData[keyBase+"commodityCode"] = lineItem.CommodityCode
		additionalData[keyBase+"description"] = sleet.TruncateString(lineItem.Description, maxLineItemDescriptionLength)
		additionalData[keyBase+"discountAmount"] = sleet.AmountToString(&lineItem.ItemDiscountAmount)
		additionalData[keyBase+"productCode"] = lineItem.ProductCode
		additionalData[keyBase+"quantity"] = strconv.Itoa(int(lineItem.Quantity))
		additionalData[keyBase+"totalAmount"] = sleet.AmountToString(&lineItem.TotalAmount)
		additionalData[keyBase+"unitOfMeasure"] = common.ConvertUnitOfMeasurementToCode(lineItem.UnitOfMeasure)
		additionalData[keyBase+"unitPrice"] = sleet.AmountToString(&lineItem.UnitPrice)
	}

	// Omit optional fields if they are empty
	addIfNonEmpty(level3Data.DestinationCountryCode, "enhancedSchemeData.destinationCountryCode", &additionalData)
	addIfNonEmpty(level3Data.DestinationAdminArea, "enhancedSchemeData.destinationStateProvinceCode", &additionalData)

	return additionalData
}

func buildCaptureRequest(captureRequest *sleet.CaptureRequest, merchantAccount string) *payments.ModificationRequest {
	request := &payments.ModificationRequest{
		OriginalReference: captureRequest.TransactionReference,
		ModificationAmount: &payments.Amount{
			Value:    captureRequest.Amount.Amount,
			Currency: captureRequest.Amount.Currency,
		},
		MerchantAccount: merchantAccount,
	}
	return request
}

func buildRefundRequest(refundRequest *sleet.RefundRequest, merchantAccount string) *payments.ModificationRequest {
	request := &payments.ModificationRequest{
		OriginalReference: refundRequest.TransactionReference,
		ModificationAmount: &payments.Amount{
			Value:    refundRequest.Amount.Amount,
			Currency: refundRequest.Amount.Currency,
		}, MerchantAccount: merchantAccount,
	}
	return request
}

func buildVoidRequest(voidRequest *sleet.VoidRequest, merchantAccount string) *payments.ModificationRequest {
	request := &payments.ModificationRequest{
		OriginalReference: voidRequest.TransactionReference,
		MerchantAccount:   merchantAccount,
	}
	return request
}

func addIfNonEmpty(value string, key string, data *map[string]string) {
	if value != "" {
		(*data)[key] = value
	}
}

package adyen

import (
	"fmt"
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/adyen/adyen-go-api-library/v2/src/checkout"
	"github.com/adyen/adyen-go-api-library/v2/src/payments"
	"strconv"
	"time"
)

const maxLineItemDescriptionLength = 26

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

	level3 := authRequest.Level3Data
	if level3 != nil {
		additionalData := map[string]string{
			"enhancedSchemeData.customerReference":            NAifEmpty(level3.CustomerReference),
			"enhancedSchemeData.destinationCountryCode":       NAifEmpty(level3.DestinationCountryCode),
			"enhancedSchemeData.destinationPostalCode":        NAifEmpty(level3.DestinationPostalCode),
			"enhancedSchemeData.destinationStateProvinceCode": NAifEmpty(level3.DestinationAdminArea),
			"enhancedSchemeData.dutyAmount":                   sleet.AmountToString(&level3.DutyAmount),
			"enhancedSchemeData.freightAmount":                sleet.AmountToString(&level3.ShippingAmount),
			// safest to use a default value for any level 2/3 fields, so assuming order is from today
			"enhancedSchemeData.orderDate":      time.Now().Format("020106"),
			"enhancedSchemeData.totalTaxAmount": sleet.AmountToString(&level3.TaxAmount),
		}

		var keyBase string
		for idx, lineItem := range level3.LineItems {
			keyBase = fmt.Sprintf("enhancedSchemeData.itemDetailLine%d.", idx+1)
			additionalData[keyBase+"commodityCode"] = NAifEmpty(lineItem.CommodityCode)
			additionalData[keyBase+"description"] = sleet.TruncateString(lineItem.Description, maxLineItemDescriptionLength)
			additionalData[keyBase+"productCode"] = NAifEmpty(lineItem.ProductCode)
			additionalData[keyBase+"discountAmount"] = sleet.AmountToString(&lineItem.ItemDiscountAmount)
			additionalData[keyBase+"quantity"] = strconv.Itoa(int(lineItem.Quantity))
			additionalData[keyBase+"totalAmount"] = sleet.AmountToString(&lineItem.TotalAmount)
			additionalData[keyBase+"unitOfMeasure"] = sleet.ConvertUnitOfMeasurementToCode(lineItem.UnitOfMeasure)
			additionalData[keyBase+"unitPrice"] = sleet.AmountToString(&lineItem.UnitPrice)
		}

		request.AdditionalData = additionalData
	}

	return request
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

// NAifEmpty returns "NA" if the parameter is an empty string. Adyen best practices
// for Level 2/3 data are to set the value of any fields that are unknown to "NA"
func NAifEmpty(str string) string {
	if str == "" {
		return "NA"
	}
	return str
}

package adyen

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/adyen/adyen-go-api-library/v4/src/checkout"
	"github.com/adyen/adyen-go-api-library/v4/src/payments"
)

const (
	level3Default                = "NA"
	maxLineItemDescriptionLength = 26
	maxProductCodeLength         = 12
)

// Options
const (
	applePayTokenOption = "ApplePayToken"
	shopperIPOption = "ShopperIP"
)

// Shopper Interactions
const (
	shopperInteractionEcommerce = "Ecommerce"
	shopperInteractionContAuth  = "ContAuth"
)

// Recurring Processing Models
const (
	recurringProcessingModelCardOnFile            = "CardOnFile"
	recurringProcessingModelSubscription          = "Subscription"
	recurringProcessingModelUnscheduledCardOnFile = "UnscheduledCardOnFile"
)

var streetNumberRegex = regexp.MustCompile(`^(\d+)\s(.*)`)

// these maps are based on https://docs.adyen.com/online-payments/tokenization/create-and-use-tokens#set-parameters-to-flag-transactions
var initiatorTypeToShopperInteraction = map[sleet.ProcessingInitiatorType]string{
	sleet.ProcessingInitiatorTypeInitialCardOnFile:         shopperInteractionEcommerce,
	sleet.ProcessingInitiatorTypeInitialRecurring:          shopperInteractionEcommerce,
	sleet.ProcessingInitiatorTypeStoredCardholderInitiated: shopperInteractionContAuth,
	sleet.ProcessingInitiatorTypeStoredMerchantInitiated:   shopperInteractionContAuth,
	sleet.ProcessingInitiatorTypeFollowingRecurring:        shopperInteractionContAuth,
}

var initiatorTypeToRecurringProcessingModel = map[sleet.ProcessingInitiatorType]string{
	sleet.ProcessingInitiatorTypeInitialCardOnFile:         recurringProcessingModelCardOnFile,
	sleet.ProcessingInitiatorTypeInitialRecurring:          recurringProcessingModelSubscription,
	sleet.ProcessingInitiatorTypeStoredCardholderInitiated: recurringProcessingModelCardOnFile,
	sleet.ProcessingInitiatorTypeStoredMerchantInitiated:   recurringProcessingModelUnscheduledCardOnFile,
	sleet.ProcessingInitiatorTypeFollowingRecurring:        recurringProcessingModelSubscription,
}

func buildAuthRequest(authRequest *sleet.AuthorizationRequest, merchantAccount string) *checkout.PaymentRequest {
	request := &checkout.PaymentRequest{
		Amount: checkout.Amount{
			Value:    authRequest.Amount.Amount,
			Currency: authRequest.Amount.Currency,
		},
		// Adyen requires a reference in request so this will panic if client doesn't pass it. Assuming this is good for now
		Reference: *authRequest.ClientTransactionReference,
		MerchantAccount:        merchantAccount,
		MerchantOrderReference: authRequest.MerchantOrderReference,

		// https://docs.adyen.com/api-explorer/#/CheckoutService/latest/payments__reqParam_shopperReference
		ShopperReference: authRequest.ShopperReference,
	}

	addPaymentMethod(authRequest, request)
	addPaymentSpecificFields(authRequest, request)
	addShopperData(authRequest, request)
	addAddresses(authRequest, request)

	// overwrites the flag transactions
	if authRequest.ProcessingInitiator != nil {
		if shopperInteraction, ok := initiatorTypeToShopperInteraction[*authRequest.ProcessingInitiator]; ok {
			request.ShopperInteraction = shopperInteraction
		}
		if recurringProcessingModel, ok := initiatorTypeToRecurringProcessingModel[*authRequest.ProcessingInitiator]; ok {
			request.RecurringProcessingModel = recurringProcessingModel
		}
	}

	// overwrites for citiplcc
	if authRequest.CreditCard.Network == sleet.CreditCardNetworkCitiPLCC {
		request.RecurringProcessingModel = recurringProcessingModelSubscription
		request.ShopperInteraction = shopperInteractionEcommerce
	}

	level3 := authRequest.Level3Data
	if level3 != nil {
		request.AdditionalData = buildLevel3Data(level3)
	}

	// Attach results of 3DS verification if performed (and not "R"ejected)
	if authRequest.ThreeDS != nil && authRequest.ThreeDS.PAResStatus != sleet.ThreedsStatusRejected {
		request.MpiData = &checkout.ThreeDSecureData{
			Cavv:              authRequest.ThreeDS.CAVV,
			CavvAlgorithm:     authRequest.ThreeDS.CAVVAlgorithm,
			DirectoryResponse: authRequest.ThreeDS.PAResStatus,
			DsTransID:         authRequest.ThreeDS.DSTransactionID,
			Eci:               authRequest.ECI,
			ThreeDSVersion:    authRequest.ThreeDS.Version,
			Xid:               authRequest.ThreeDS.XID,
		}
		// Only pass these fields for challenge flow
		if !authRequest.ThreeDS.Frictionless {
			request.MpiData.AuthenticationResponse = authRequest.ThreeDS.PAResStatus
		}
	}

	return request
}

// addPaymentMethod add payment method to the adyen payment request.
func addPaymentMethod(authRequest *sleet.AuthorizationRequest, request *checkout.PaymentRequest) {
	if authRequest.Options[applePayTokenOption] != nil {
		request.PaymentMethod = map[string]interface{}{
			"type": "applepay",
			"applePayToken": authRequest.Options[applePayTokenOption].(string),
		}
	} else {
		request.PaymentMethod = map[string]interface{}{
			"expiryMonth": strconv.Itoa(authRequest.CreditCard.ExpirationMonth),
			"expiryYear":  strconv.Itoa(authRequest.CreditCard.ExpirationYear),
			"holderName":  authRequest.CreditCard.FirstName + " " + authRequest.CreditCard.LastName,
			"number":      authRequest.CreditCard.Number,
			"type":        "scheme",
		}
	}
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
		request.PaymentMethod["brand"] = "applepay"
		request.RecurringProcessingModel = recurringProcessingModelSubscription
		request.ShopperInteraction = shopperInteractionEcommerce
	} else if authRequest.CreditCard.CVV != "" {
		// New customer credit card request
		request.PaymentMethod["cvc"] = authRequest.CreditCard.CVV
		request.ShopperInteraction = shopperInteractionEcommerce
		if authRequest.CreditCard.Save {
			// Customer opts in to saving card details
			request.RecurringProcessingModel = recurringProcessingModelCardOnFile
			request.StorePaymentMethod = true
		} else {
			// Customer opts out of saving card details
			request.StorePaymentMethod = false
		}
	} else {
		// Existing customer credit card request
		request.RecurringProcessingModel = recurringProcessingModelCardOnFile
		request.ShopperInteraction = shopperInteractionContAuth
	}
}

// addAddresses adds the billing address and shipping address to the Ayden Payment request if available
func addAddresses(authRequest *sleet.AuthorizationRequest, request *checkout.PaymentRequest) {
	if authRequest.BillingAddress != nil {
		billingStreetNumber, billingStreetName := extractAdyenStreetFormat(common.SafeStr(authRequest.BillingAddress.StreetAddress1))
		request.BillingAddress = &checkout.Address{
			City:              common.SafeStr(authRequest.BillingAddress.Locality),
			Country:           common.SafeStr(authRequest.BillingAddress.CountryCode),
			HouseNumberOrName: billingStreetNumber,
			PostalCode:        common.SafeStr(authRequest.BillingAddress.PostalCode),
			StateOrProvince:   common.SafeStr(authRequest.BillingAddress.RegionCode),
			Street:            billingStreetName,
		}
	}
	if authRequest.ShippingAddress != nil {
		shippingStreetNumber, shippingStreetName := extractAdyenStreetFormat(common.SafeStr(authRequest.ShippingAddress.StreetAddress1))
		request.DeliveryAddress = &checkout.Address{
			City:              common.SafeStr(authRequest.ShippingAddress.Locality),
			Country:           common.SafeStr(authRequest.ShippingAddress.CountryCode),
			HouseNumberOrName: shippingStreetNumber,
			PostalCode:        common.SafeStr(authRequest.ShippingAddress.PostalCode),
			StateOrProvince:   common.SafeStr(authRequest.ShippingAddress.RegionCode),
			Street:            shippingStreetName,
		}
	}
}

// addShopperData adds the shoppers IP and email to the Ayden Payment request if available
func addShopperData(authRequest *sleet.AuthorizationRequest, request *checkout.PaymentRequest) {
	if authRequest.Options[shopperIPOption] != nil {
		request.ShopperIP = authRequest.Options[shopperIPOption].(string)
	}
	if authRequest.BillingAddress.Email != nil {
		request.ShopperEmail = common.SafeStr(authRequest.BillingAddress.Email)
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
		// Due to issues with the credit card networks, dont send any line item if discount amount is 0
		if lineItem.ItemDiscountAmount.Amount > 0 {
			additionalData[keyBase+"discountAmount"] = sleet.AmountToString(&lineItem.ItemDiscountAmount)
		}
		additionalData[keyBase+"commodityCode"] = lineItem.CommodityCode
		additionalData[keyBase+"description"] = sleet.TruncateString(lineItem.Description, maxLineItemDescriptionLength)
		additionalData[keyBase+"productCode"] = sleet.TruncateString(lineItem.ProductCode, maxProductCodeLength)
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

// extractAdyenStreetFormat extracts adyen street format from generic street address
//                          returns (streetNumber, streetName) format
//                          If address does not have leading street number, will return ("", street)
func extractAdyenStreetFormat(streetAddress string) (string, string) {
	streetExtraction := streetNumberRegex.FindStringSubmatch(streetAddress)
	if streetExtraction == nil {
		return "", streetAddress
	}

	return streetExtraction[1], streetExtraction[2]
}

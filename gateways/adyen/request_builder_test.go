// +build unit

package adyen

import (
	"github.com/BoltApp/sleet"
	"github.com/adyen/adyen-go-api-library/v4/src/checkout"
	"github.com/go-test/deep"
	"strconv"
	"testing"

	sleet_testing "github.com/BoltApp/sleet/testing"
)

func TestBuildAuthRequest(t *testing.T) {
	base := sleet_testing.BaseAuthorizationRequest()
	requestWithLevel3Data := sleet_testing.BaseAuthorizationRequest()
	requestWithLevel3Data.Level3Data = sleet_testing.BaseLevel3Data()
	requestWithLevel3ItemDiscount := sleet_testing.BaseAuthorizationRequest()
	requestWithLevel3ItemDiscount.Level3Data = sleet_testing.BaseLevel3Data()
	requestWithLevel3ItemDiscount.Level3Data.LineItems[0].ItemDiscountAmount.Amount = 100

	cases := []struct {
		label string
		in    *sleet.AuthorizationRequest
		want  *checkout.PaymentRequest
	}{
		{
			"Basic Auth Request",
			base,
			&checkout.PaymentRequest{
				Amount: checkout.Amount{
					Currency: "USD",
					Value:    100,
				},
				BillingAddress: &checkout.Address{
					City:            *base.BillingAddress.Locality,
					Country:         *base.BillingAddress.CountryCode,
					PostalCode:      *base.BillingAddress.PostalCode,
					StateOrProvince: *base.BillingAddress.RegionCode,
					Street:          *base.BillingAddress.StreetAddress1,
				},
				MerchantAccount: "merchant-account",
				PaymentMethod: map[string]interface{}{
					"number":      base.CreditCard.Number,
					"expiryMonth": strconv.Itoa(base.CreditCard.ExpirationMonth),
					"expiryYear":  strconv.Itoa(base.CreditCard.ExpirationYear),
					"holderName":  base.CreditCard.FirstName + " " + base.CreditCard.LastName,
					"cvc":         base.CreditCard.CVV,
					"type":        "scheme",
				},
				ShopperInteraction:       "Ecommerce",
				RecurringProcessingModel: "CardOnFile",
				Reference:                *base.ClientTransactionReference,
				StorePaymentMethod:       true,
				ShopperReference: 		  *base.ClientTransactionReference,
			},
		},
		{
			"Auth Request with level3 data no Item Discount",
			requestWithLevel3Data,
			&checkout.PaymentRequest{
				Amount: checkout.Amount{
					Currency: "USD",
					Value:    100,
				},
				BillingAddress: &checkout.Address{
					City:            *requestWithLevel3Data.BillingAddress.Locality,
					Country:         *requestWithLevel3Data.BillingAddress.CountryCode,
					PostalCode:      *requestWithLevel3Data.BillingAddress.PostalCode,
					StateOrProvince: *requestWithLevel3Data.BillingAddress.RegionCode,
					Street:          *requestWithLevel3Data.BillingAddress.StreetAddress1,
				},
				MerchantAccount: "merchant-account",
				PaymentMethod: map[string]interface{}{
					"number":      requestWithLevel3Data.CreditCard.Number,
					"expiryMonth": strconv.Itoa(requestWithLevel3Data.CreditCard.ExpirationMonth),
					"expiryYear":  strconv.Itoa(requestWithLevel3Data.CreditCard.ExpirationYear),
					"holderName":  requestWithLevel3Data.CreditCard.FirstName + " " + requestWithLevel3Data.CreditCard.LastName,
					"cvc":         requestWithLevel3Data.CreditCard.CVV,
					"type":        "scheme",
				},
				ShopperInteraction:       "Ecommerce",
				RecurringProcessingModel: "CardOnFile",
				Reference:                *requestWithLevel3Data.ClientTransactionReference,
				StorePaymentMethod:       true,
				AdditionalData: map[string]string{
					"enhancedSchemeData.totalTaxAmount":                "100",
					"enhancedSchemeData.freightAmount":                 "300",
					"enhancedSchemeData.dutyAmount":                    "400",
					"enhancedSchemeData.itemDetailLine1.description":   "pot",
					"enhancedSchemeData.itemDetailLine1.commodityCode": "cmd",
					"enhancedSchemeData.itemDetailLine1.productCode":   "abc",
					"enhancedSchemeData.itemDetailLine1.unitOfMeasure": "EA",
					"enhancedSchemeData.itemDetailLine1.quantity":      "2",
					"enhancedSchemeData.itemDetailLine1.unitPrice":     "500",
					"enhancedSchemeData.itemDetailLine1.totalAmount":   "1000",
					"enhancedSchemeData.destinationCountryCode":        "US",
					"enhancedSchemeData.customerReference":             "customer",
					"enhancedSchemeData.destinationPostalCode":         "94105",
				},
			},
		},
		{
			"Auth Request with level3 data with Item Discount",
			requestWithLevel3ItemDiscount,
			&checkout.PaymentRequest{
				Amount: checkout.Amount{
					Currency: "USD",
					Value:    100,
				},
				BillingAddress: &checkout.Address{
					City:            *requestWithLevel3Data.BillingAddress.Locality,
					Country:         *requestWithLevel3Data.BillingAddress.CountryCode,
					PostalCode:      *requestWithLevel3Data.BillingAddress.PostalCode,
					StateOrProvince: *requestWithLevel3Data.BillingAddress.RegionCode,
					Street:          *requestWithLevel3Data.BillingAddress.StreetAddress1,
				},
				MerchantAccount: "merchant-account",
				PaymentMethod: map[string]interface{}{
					"number":      requestWithLevel3Data.CreditCard.Number,
					"expiryMonth": strconv.Itoa(requestWithLevel3Data.CreditCard.ExpirationMonth),
					"expiryYear":  strconv.Itoa(requestWithLevel3Data.CreditCard.ExpirationYear),
					"holderName":  requestWithLevel3Data.CreditCard.FirstName + " " + requestWithLevel3Data.CreditCard.LastName,
					"cvc":         requestWithLevel3Data.CreditCard.CVV,
					"type":        "scheme",
				},
				ShopperInteraction:       "Ecommerce",
				RecurringProcessingModel: "CardOnFile",
				Reference:                *requestWithLevel3ItemDiscount.ClientTransactionReference,
				StorePaymentMethod:       true,
				AdditionalData: map[string]string{
					"enhancedSchemeData.totalTaxAmount":                 "100",
					"enhancedSchemeData.freightAmount":                  "300",
					"enhancedSchemeData.dutyAmount":                     "400",
					"enhancedSchemeData.itemDetailLine1.description":    "pot",
					"enhancedSchemeData.itemDetailLine1.commodityCode":  "cmd",
					"enhancedSchemeData.itemDetailLine1.productCode":    "abc",
					"enhancedSchemeData.itemDetailLine1.unitOfMeasure":  "EA",
					"enhancedSchemeData.itemDetailLine1.quantity":       "2",
					"enhancedSchemeData.itemDetailLine1.unitPrice":      "500",
					"enhancedSchemeData.itemDetailLine1.totalAmount":    "1000",
					"enhancedSchemeData.itemDetailLine1.discountAmount": "100",
					"enhancedSchemeData.destinationCountryCode":         "US",
					"enhancedSchemeData.customerReference":              "customer",
					"enhancedSchemeData.destinationPostalCode":          "94105",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			got := buildAuthRequest(c.in, "merchant-account")
			if diff := deep.Equal(got, c.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}

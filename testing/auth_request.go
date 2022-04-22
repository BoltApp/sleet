package testing

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/Pallinder/go-randomdata"
)

// BaseAuthorizationRequest is used as a testing helper method to standardize request calls for integration tests
func BaseAuthorizationRequest() *sleet.AuthorizationRequest {
	amount := sleet.Amount{
		Amount:   100,
		Currency: "USD",
	}
	address := sleet.Address{
		PostalCode:     common.SPtr("94103"),
		CountryCode:    common.SPtr("US"),
		StreetAddress1: common.SPtr("7683 Railroad Street"),
		Locality:       common.SPtr("Zion"),
		RegionCode:     common.SPtr("IL"),
	}
	card := sleet.CreditCard{
		FirstName:       "Bolt",
		LastName:        "Checkout",
		Number:          "4111111111111111",
		ExpirationMonth: 10,
		ExpirationYear:  2023,
		CVV:             "737",
		Save:            true,
	}
	reference := randomdata.Letters(10)
	return &sleet.AuthorizationRequest{Amount: amount, CreditCard: &card, BillingAddress: &address, ClientTransactionReference: &reference, ShopperReference: "test"}
}

func BaseAuthorizationRequestWithEmailPhoneNumber() *sleet.AuthorizationRequest {
	base := BaseAuthorizationRequest()
	base.BillingAddress.Email = common.SPtr("test@bolt.com")
	base.BillingAddress.PhoneNumber = common.SPtr("555-555-5555")
	return base
}

// BaseLevel3Data is used as a testing helper method to standardize request calls for integration tests
func BaseLevel3Data() *sleet.Level3Data {
	return &sleet.Level3Data{
		CustomerReference: "customer",
		TaxAmount: sleet.Amount{
			Amount:   100,
			Currency: "USD",
		},
		DiscountAmount: sleet.Amount{
			Amount:   200,
			Currency: "USD",
		},
		ShippingAmount: sleet.Amount{
			Amount:   300,
			Currency: "USD",
		},
		DutyAmount: sleet.Amount{
			Amount:   400,
			Currency: "USD",
		},
		DestinationPostalCode:  "94105",
		DestinationCountryCode: "US",
		LineItems: []sleet.LineItem{
			{
				Description: "pot",
				ProductCode: "abc",
				UnitPrice: sleet.Amount{
					Amount:   500,
					Currency: "USD",
				},
				Quantity: 2,
				TotalAmount: sleet.Amount{
					Amount:   1000,
					Currency: "USD",
				},
				UnitOfMeasure: "count",
				CommodityCode: "cmd",
			},
		},
	}
}

func BaseLevel3DataMultipleItem() *sleet.Level3Data {
	return &sleet.Level3Data{
		CustomerReference: "customer",
		TaxAmount: sleet.Amount{
			Amount:   100,
			Currency: "USD",
		},
		DiscountAmount: sleet.Amount{
			Amount:   200,
			Currency: "USD",
		},
		ShippingAmount: sleet.Amount{
			Amount:   300,
			Currency: "USD",
		},
		DutyAmount: sleet.Amount{
			Amount:   400,
			Currency: "USD",
		},
		DestinationPostalCode:  "94105",
		DestinationCountryCode: "US",
		LineItems: []sleet.LineItem{
			{
				Description: "pot",
				ProductCode: "abc",
				UnitPrice: sleet.Amount{
					Amount:   500,
					Currency: "USD",
				},
				Quantity: 2,
				TotalAmount: sleet.Amount{
					Amount:   1000,
					Currency: "USD",
				},
				UnitOfMeasure: "count",
				CommodityCode: "cmd",
			},
			{
				Description: "vase",
				ProductCode: "123",
				UnitPrice: sleet.Amount{
					Amount:   1000,
					Currency: "USD",
				},
				Quantity: 5,
				TotalAmount: sleet.Amount{
					Amount:   2000,
					Currency: "USD",
				},
				UnitOfMeasure: "count",
				CommodityCode: "321",
			},
		},
	}
}

func BaseCaptureRequest() *sleet.CaptureRequest {
	clientRef := "222222"

	amount := sleet.Amount{
		Amount:   100,
		Currency: "USD",
	}
	return &sleet.CaptureRequest{
		Amount:                     &amount,
		TransactionReference:       "111111",
		ClientTransactionReference: &clientRef,
	}
}

func BaseVoidRequest() *sleet.VoidRequest {
	clientRef := "222222"

	return &sleet.VoidRequest{
		TransactionReference:       "111111",
		ClientTransactionReference: &clientRef,
	}
}

func BaseRefundRequest() *sleet.RefundRequest {
	clientRef := "222222"

	amount := sleet.Amount{
		Amount:   100,
		Currency: "USD",
	}
	return &sleet.RefundRequest{
		Amount:                     &amount,
		TransactionReference:       "111111",
		ClientTransactionReference: &clientRef,
		Last4:                      "1111",
	}
}

// Base3DS provides a template with 3DS authorization data. Fields are populated by their names
// since no valid values will exist without first having run 3DS.
func Base3DS() *sleet.ThreeDS {
	return &sleet.ThreeDS{
		Frictionless:     false,
		ACSTransactionID: "acs-transaction-id",
		CAVV:             "cavv",
		CAVVAlgorithm:    "cavv-algorithm",
		DSTransactionID:  "ds-transaction-id",
		PAResStatus:      "pares-status",
		UCAFIndicator:    "ucaf-indicator",
		Version:          "version",
		XID:              "xid",
	}
}

func BaseCaptureRequestWithOptions() *sleet.CaptureRequest {
	clientRef := "222222"

	amount := sleet.Amount{
		Amount:   100,
		Currency: "USD",
	}
	return &sleet.CaptureRequest{
		Amount:                     &amount,
		TransactionReference:       "111111",
		ClientTransactionReference: &clientRef,
		Options: map[string]interface{}{
			"captureSequenceNumber": "11",
			"totalCaptureCount":     "99",
		},
	}
}

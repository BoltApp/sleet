package cybersource

import (
	"errors"
	"fmt"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"strconv"
)

const (
	InitiatorTypeMerchant = "merchant"
	InitiatorTypeConsumer = "consumer"
)

const (
	AmexCryptogramMaxLength   = 40
	AmexCryptogramSplitLength = 20
	TransactionTypeInApp = "1"
	PaymentSolutionApplepay = "001"
)

// add mappings
var initiatorTypeToInitiatorType = map[sleet.ProcessingInitiatorType]string{
	sleet.ProcessingInitiatorTypeInitialCardOnFile:         InitiatorTypeConsumer,
	sleet.ProcessingInitiatorTypeInitialRecurring:          InitiatorTypeConsumer,
	sleet.ProcessingInitiatorTypeStoredCardholderInitiated: InitiatorTypeConsumer,
	sleet.ProcessingInitiatorTypeStoredMerchantInitiated:   InitiatorTypeMerchant,
	sleet.ProcessingInitiatorTypeFollowingRecurring:        InitiatorTypeMerchant,
}

var initiatorTypeToCredentialStoredOnFile = map[sleet.ProcessingInitiatorType]bool{
	sleet.ProcessingInitiatorTypeInitialCardOnFile:         true,
	sleet.ProcessingInitiatorTypeInitialRecurring:          true,
	sleet.ProcessingInitiatorTypeStoredCardholderInitiated: false,
	sleet.ProcessingInitiatorTypeStoredMerchantInitiated:   false,
	sleet.ProcessingInitiatorTypeFollowingRecurring:        false,
}

var initiatorTypeToStoredCredentialUsed = map[sleet.ProcessingInitiatorType]bool{
	sleet.ProcessingInitiatorTypeInitialCardOnFile:         false,
	sleet.ProcessingInitiatorTypeInitialRecurring:          false,
	sleet.ProcessingInitiatorTypeStoredCardholderInitiated: true,
	sleet.ProcessingInitiatorTypeStoredMerchantInitiated:   true,
	sleet.ProcessingInitiatorTypeFollowingRecurring:        true,
}

func buildAuthRequest(authRequest *sleet.AuthorizationRequest) (*Request, error) {
	var initiatorType string
	var credentialStoredOnFile bool
	var storedCredentialUsed bool
	if authRequest.ProcessingInitiator != nil {
		initiatorType = initiatorTypeToInitiatorType[*authRequest.ProcessingInitiator]
		credentialStoredOnFile = initiatorTypeToCredentialStoredOnFile[*authRequest.ProcessingInitiator]
		storedCredentialUsed = initiatorTypeToStoredCredentialUsed[*authRequest.ProcessingInitiator]
	}

	amountStr := sleet.AmountToDecimalString(&authRequest.Amount)
	request := &Request{
		ClientReferenceInformation: &ClientReferenceInformation{
			Code: authRequest.MerchantOrderReference,
			Partner: Partner{
				SolutionID: authRequest.Channel,
			},
		},
		ProcessingInformation: &ProcessingInformation{
			Capture:           false, // no autocapture for now
			CommerceIndicator: string(CommerceIndicatorInternet),
			AuthorizationOptions: &AuthorizationOptions{
				Initiator: &Initiator{
					InitiatorType: initiatorType,
					CredentialStoredOnFile: credentialStoredOnFile,
					StoredCredentialUsed: storedCredentialUsed,
				},
			},
		},
		PaymentInformation: &PaymentInformation{
			Card: &CardInformation{
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

	// Apple Pay request
	if authRequest.Cryptogram != "" {
		err := buildApplepayRequest(authRequest, request)
		if err != nil {
			return nil, err
		}
	}

	// If level 3 data is present, and ClientReferenceInformation in that data exists, it will override this.
	if authRequest.ClientTransactionReference != nil {
		request.MerchantDefinedInformation = append(request.MerchantDefinedInformation, MerchantDefinedInformation{
			Key: "1",
			Value: *authRequest.ClientTransactionReference,
		})
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
	amountStr := sleet.AmountToDecimalString(captureRequest.Amount)
	request := &Request{
		OrderInformation: &OrderInformation{
			AmountDetails: AmountDetails{
				Amount:   amountStr,
				Currency: captureRequest.Amount.Currency,
			},
		},
	}
	if captureRequest.MerchantOrderReference != nil {
		request.ClientReferenceInformation = &ClientReferenceInformation{
			Code: *captureRequest.MerchantOrderReference,
		}
	}
	if captureRequest.ClientTransactionReference != nil {
		request.MerchantDefinedInformation = append(request.MerchantDefinedInformation, MerchantDefinedInformation{
			Key: "1",
			Value: *captureRequest.ClientTransactionReference,
		})
	}
	return request, nil
}

func buildVoidRequest(voidRequest *sleet.VoidRequest) (*Request, error) {
	// Maybe add reason / more details, but for now nothing
	request := &Request{}
	if voidRequest.MerchantOrderReference != nil {
		request.ClientReferenceInformation = &ClientReferenceInformation{
			Code: *voidRequest.MerchantOrderReference,
		}
	}
	if voidRequest.ClientTransactionReference != nil {
		request.MerchantDefinedInformation = append(request.MerchantDefinedInformation, MerchantDefinedInformation{
			Key: "1",
			Value: *voidRequest.ClientTransactionReference,
		})
	}
	return request, nil
}

func buildRefundRequest(refundRequest *sleet.RefundRequest) (*Request, error) {
	amountStr := sleet.AmountToDecimalString(refundRequest.Amount)
	request := &Request{
		OrderInformation: &OrderInformation{
			AmountDetails: AmountDetails{
				Amount:   amountStr,
				Currency: refundRequest.Amount.Currency,
			},
		},
	}
	if refundRequest.MerchantOrderReference != nil {
		request.ClientReferenceInformation = &ClientReferenceInformation{
			Code: *refundRequest.MerchantOrderReference,
		}
	}
	if refundRequest.ClientTransactionReference != nil {
		request.MerchantDefinedInformation = append(request.MerchantDefinedInformation, MerchantDefinedInformation{
			Key: "1",
			Value: *refundRequest.ClientTransactionReference,
		})
	}
	return request, nil
}

func buildApplepayRequest(authRequest *sleet.AuthorizationRequest, request *Request) error {
	request.PaymentInformation = &PaymentInformation{
		TokenizedCard: &TokenizedCard{
			Number: authRequest.CreditCard.Number,
			ExpirationYear: strconv.Itoa(authRequest.CreditCard.ExpirationYear),
			ExpirationMonth: fmt.Sprintf("%02d", authRequest.CreditCard.ExpirationMonth),
			TransactionType: TransactionTypeInApp,
			Cryptogram: authRequest.Cryptogram,
		},
	}

	request.ProcessingInformation.PaymentSolution = PaymentSolutionApplepay

	switch authRequest.CreditCard.Network {
	case sleet.CreditCardNetworkVisa:
		request.PaymentInformation.TokenizedCard.Type = string(CardTypeVisa)
		request.ConsumerAuthenticationInformation = &ConsumerAuthenticationInformation{
			Xid: authRequest.Cryptogram,
			Cavv: authRequest.Cryptogram,
		}
	case sleet.CreditCardNetworkMastercard:
		request.PaymentInformation.TokenizedCard.Type = string(CardTypeMastercard)
		request.ProcessingInformation.CommerceIndicator = string(CommerceIndicatorMastercard)
		request.ConsumerAuthenticationInformation = &ConsumerAuthenticationInformation{
			UcafAuthenticationData: authRequest.Cryptogram,
			UcafCollectionIndicator: "2",
		}
	case sleet.CreditCardNetworkAmex:
		request.PaymentInformation.TokenizedCard.Type = string(CardTypeAmex)
		request.ProcessingInformation.CommerceIndicator = string(CommerceIndicatorAmex)
		consumerAuthInfo, err := getAmexConsumerAuthInfo(authRequest.Cryptogram)
		if err != nil {
			return err
		}
		request.ConsumerAuthenticationInformation = consumerAuthInfo
	case sleet.CreditCardNetworkDiscover:
		request.PaymentInformation.TokenizedCard.Type = string(CardTypeDiscover)
		request.ProcessingInformation.CommerceIndicator = string(CommerceIndicatorDiscover)
		request.ConsumerAuthenticationInformation = &ConsumerAuthenticationInformation{
			Cavv: authRequest.Cryptogram,
		}
	default:
		return errors.New("unsupported payment method")
	}

	return nil
}

func getAmexConsumerAuthInfo(cryptogram string) (*ConsumerAuthenticationInformation, error) {
	if len(cryptogram) > AmexCryptogramMaxLength {
		return nil, errors.New("invalid Amex cryptogram length")
	} else if len(cryptogram) == AmexCryptogramMaxLength {
		// For a 40-byte cryptogram, split the cryptogram into two 20-byte binary values (block A and block B).
		// Send the first 20-byte value (block A) in the cavv field. Send the second 20-byte value (block B) in
		// the xid field.
		return &ConsumerAuthenticationInformation{
			Cavv: cryptogram[:AmexCryptogramSplitLength],
			Xid: cryptogram[AmexCryptogramSplitLength:],
		}, nil
	}

	return &ConsumerAuthenticationInformation{
		Cavv: cryptogram,
	}, nil
}

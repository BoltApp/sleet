# Sleet

[![CircleCI status](https://circleci.com/gh/BoltApp/sleet.png?circle-token=d60ceb64eb6ebdfd6a45a4703563c1752598db63 "CircleCI status")](https://circleci.com/gh/BoltApp/sleet)
[![GoDoc](https://godoc.org/github.com/BoltApp/sleet?status.svg)](https://godoc.org/github.com/BoltApp/sleet)

Payment abstraction library - interact with different Payment Service Providers with one unified interface.

## Installation
`go get github.com/BoltApp/sleet`

## Methodology
Wherever possible, we try to use native Golang implementations of the PsP's API. We also assume that the caller can pass along raw credit card information (i.e. are PCI compliant) 

### Supported API Calls
1. Authorize
2. Capture
3. Void
4. Refund

### To run tests
The following environment variables are needed in order to run tests
```shell script
$ export ADYEN_USERNAME="YOUR_ADYEN_WEBSERVICE_USERNAME"
$ export ADYEN_ACCOUNT="YOUR_ADYEN_MERCHANT_ACCOUNT"
$ export ADYEN_PASSWORD="YOUR_ADYEN_WEBSERVICE_PASSWORD"
$ export STRIPE_TEST_KEY="YOUR_STRIPE_API_KEY"
$ export AUTH_NET_LOGIN_ID="YOUR_AUTHNET_LOGIN"
$ export AUTH_NET_TXN_KEY="YOUR_AUTHNET_TXN_KEY"
$ export BRAINTREE_MERCHANT_ID="YOUR_BRAINTREE_MERCHANT_ACCOUNT"
$ export BRAINTREE_PUBLIC_KEY="YOUR_BRAINTREE_PUBLIC_KEY"
$ export BRAINTREE_PRIVATE_ID="YOUR_BRAINTREE_PRIVATE_KEY"
$ export CYBERSOURCE_ACCOUNT="YOUR_CYBS_ACCOUNT"
$ export CYBERSOURCE_API_KEY="YOUR_CYBS_KEY"
$ export CYBERSOURCE_SHARED_SECRET="YOUR_CYBS_SECRET"
``` 
Then run tests with: `go test ./integration-tests/`

#### Code Example for Auth + Capture

```
import (
  "github.com/BoltApp/sleet"
  "github.com/BoltApp/sleet/gateways/authorize_net"
)
// Generate a client using your own credentials
client := authorize_net.NewClient("AUTH_NET_LOGIN_ID", "AUTH_NET_TXN_KEY")

amount := sleet.Amount{
  Amount: 100,
  Currency: "USD",
}
card := sleet.CreditCard{
  FirstName: "Bolt",
  LastName: "Checkout",
  Number: "4111111111111111",
  ExpMonth: 8,
  EpxYear: 2010,
  CVV: "000",
}
streetAddress := "22 Linda St."
locality := "Hoboken"
regionCode := "NJ"
postalCode := "07030"
countryCode := "US"
address := sleet.BillingAddress{
  StreetAddress1: &streetAddress,
  Locality:       &locality,
  RegionCode:     &regionCode,
  PostalCode:     &postalCode,
  CountryCode:    &countryCode,
}
authorizeRequest := sleet.AuthorizationRequest{
  Amount: &amount,
  CreditCard: &card,
  BillingAddress: &address,
}
authorizeResponse, _ := client.Authorize(&authorizeRequest) 

captureRequest := sleet.CaptureRequest{
  Amount:               &amount,
  TransactionReference: authorizeResponse.TransactionReference,
}
client.Capture(&captureRequest)
```

#### Supported Gateways
* [Authorize.Net](https://developer.authorize.net/api/reference/index.html#payment-transactions)
* [CyberSource](https://developer.cybersource.com/api-reference-assets/index.html#payments)
* [Stripe](https://stripe.com/docs/api)
* [Adyen](https://docs.adyen.com/classic-integration/api-integration-ecommerce)
* [Braintree](https://www.braintreepayments.com/)

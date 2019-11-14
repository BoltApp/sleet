### Sleet

[![CircleCI status](https://circleci.com/gh/BoltApp/sleet.png?circle-token=d60ceb64eb6ebdfd6a45a4703563c1752598db63 "CircleCI status")](https://circleci.com/gh/BoltApp/sleet)

Payment abstraction library - one interface for multiple payment processors

#### Supported operations
1. authorize
2. capture
3. void
4. refund

#### Example

```
import (
  "github.com/BoltApp/sleet"
  "github.com/BoltApp/sleet/gateways/authorize_net"
)
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
  TransactionReference: resp.TransactionReference,
}
client.Capture(&captureRequest)
```

#### Supported Gateways
* [Authorize.Net](https://developer.authorize.net/api/reference/index.html#payment-transactions)
* [CyberSource](https://developer.cybersource.com/api-reference-assets/index.html#payments)
* [Stripe](https://stripe.com/docs/api)

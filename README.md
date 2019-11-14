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
client := sleet.NewStripeClient("stripe_api_key", sleet.ModeTest)
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
client.Authorize(amount, card) 
```
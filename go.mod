module github.com/BoltApp/sleet

go 1.14

require (
	github.com/Pallinder/go-randomdata v1.2.0
	github.com/jarcoal/httpmock v1.0.4
	github.com/stretchr/testify v1.4.0
	github.com/stripe/stripe-go v70.11.0+incompatible
	github.com/zhutik/adyen-api-go v0.0.0-20200322135654-771fae870bd2
	golang.org/x/net v0.0.0-20200421231249-e086a090c8fd // indirect
	golang.org/x/text v0.3.2 // indirect
)

replace github.com/zhutik/adyen-api-go => github.com/nirajjayantbolt/adyen-api-go v0.0.0-20200427202302-d60958b1d31e

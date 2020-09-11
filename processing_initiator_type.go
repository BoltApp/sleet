package sleet

type ProcessingInitiatorType string

const (
	ProcessingInitiatorTypeInitialCardOnFile         ProcessingInitiatorType = "initial_card_on_file"
	ProcessingInitiatorTypeInitialRecurring          ProcessingInitiatorType = "initial_recurring"
	ProcessingInitiatorTypeStoredCardholderInitiated ProcessingInitiatorType = "stored_cardholder_initiated"
	ProcessingInitiatorTypeStoredMerchantInitiated   ProcessingInitiatorType = "stored_merchant_initiated"
	ProcessingInitiatorTypeFollowingRecurring        ProcessingInitiatorType = "following_recurring"
)

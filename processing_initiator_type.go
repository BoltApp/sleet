package sleet

// ProcessingInitiatorType type of processing initiator
type ProcessingInitiatorType string

const (
	// ProcessingInitiatorTypeInitialCardOnFile initial non-recurring payment
	ProcessingInitiatorTypeInitialCardOnFile         ProcessingInitiatorType = "initial_card_on_file"
	// ProcessingInitiatorTypeInitialRecurring initiating recurring payment
	ProcessingInitiatorTypeInitialRecurring          ProcessingInitiatorType = "initial_recurring"
	// ProcessingInitiatorTypeStoredCardholderInitiated initiated by cardholder using stored card
	ProcessingInitiatorTypeStoredCardholderInitiated ProcessingInitiatorType = "stored_cardholder_initiated"
	// ProcessingInitiatorTypeStoredMerchantInitiated initiated by merchant using stored card
	ProcessingInitiatorTypeStoredMerchantInitiated   ProcessingInitiatorType = "stored_merchant_initiated"
	// ProcessingInitiatorTypeFollowingRecurring recurring payment
	ProcessingInitiatorTypeFollowingRecurring        ProcessingInitiatorType = "following_recurring"
)

package webhooks

import "github.com/BoltApp/sleet"

// WebhookTranslator Sleet interface which takes an eventBody and translates the body to the
// Sleet TransactionEvent. Normalizes all fields to the structure and enums defined by Sleet.
type WebhookTranslator interface {
	Translate(eventBody *string) (*TransactionEvent, error)
}

type TransactionEventType int

const (
	CaptureEvent TransactionEventType = 1
	VoidEvent	 TransactionEventType = 2
	RefundEvent	 TransactionEventType = 3
)

type TransactionEvent struct {
	// Core event fields, these should be included in every transaction event
	transactionEventType    TransactionEventType  // Normalized event type for the transaction
	transactionReferenceId	*string               // id representing the transaction on the PsP system
	success					bool				  // normalized indicator on event representing success or not

	// Optional event fields, may be available based on the event type or processor implementation
	merchantTransactionReferenceId *string 	      // id representing the transaction on caller's system (as passed to the PsP if supported)
	amount sleet.Amount

}

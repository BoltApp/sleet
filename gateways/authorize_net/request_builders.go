package authorize_net

import (
	"fmt"
	"github.com/BoltApp/sleet"
)

func buildCaptureRequest(merchantName string, transactionKey string, captureRequest *sleet.CaptureRequest) (*Request, error) {
	amount := fmt.Sprintf("%.2f", float64(captureRequest.Amount.Amount) / 100.0))
	request := &Request{
		CreateTransactionRequest:CreateTransactionRequest{
			MerchantAuthentication: MerchantAuthentication{
				Name:           merchantName,
				TransactionKey: transactionKey,
			},
			TransactionRequest: TransactionRequest{
				TransactionType:  "priorAuthCaptureTransaction",
				Amount:           &amount,
				RefTransactionID: captureRequest.TransactionReference,
			},
		},
	}
	return request, nil
}

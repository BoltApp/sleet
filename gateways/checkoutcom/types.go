package checkoutcom

type CVVResponseCode string

// See https://www.checkout.com/docs/resources/codes/cvv-response-codes
const (
	CVVResponseNotPresent    CVVResponseCode = "X"
	CVVResponseNotConfigured CVVResponseCode = "U"
	CVVResponseCVDMissing    CVVResponseCode = "P"
	CVVResponseMatched       CVVResponseCode = "Y"
	CVVResponseNotValid      CVVResponseCode = "D"
	CVVResponseFailed        CVVResponseCode = "N"
)

type AVSResponseCode string

// See https://www.checkout.com/docs/resources/codes/avs-codes
const (
	AVSResponseStreetMatch                                 AVSResponseCode = "A"
	AVSResponseStreetMatchPostalUnverified                 AVSResponseCode = "B"
	AVSResponseStreetAndPostalUnverified                   AVSResponseCode = "C"
	AVSResponseStreetAndPostalMatch                        AVSResponseCode = "D"
	AVSResponseAddressMatchError                           AVSResponseCode = "E"
	AVSResponseStreetAndPostalMatchUK                      AVSResponseCode = "F"
	AVSResponseNotVerifiedOrNotSupported                   AVSResponseCode = "G"
	AVSResponseAddressUnverified                           AVSResponseCode = "I"
	AVSResponseStreetAndPostalMatchMIntl                   AVSResponseCode = "M"
	AVSResponseNoAddressMatch                              AVSResponseCode = "N"
	AVSResponseAVSNotRequested                             AVSResponseCode = "O"
	AVSResponseStreetUnverifiedPostalMatch                 AVSResponseCode = "P"
	AVSResponseAVSUnavailable                              AVSResponseCode = "R"
	AVSResponseAVSUnsupported                              AVSResponseCode = "S"
	AVSResponseMatchNotCapable                             AVSResponseCode = "U"
	AVSResponseNineDigitPostalMatch                        AVSResponseCode = "W"
	AVSResponseStreetAndNineDigitPostalMatch               AVSResponseCode = "X"
	AVSResponseStreetAndFiveDigitPostalMatch               AVSResponseCode = "Y"
	AVSResponseFiveDigitPostalMatch                        AVSResponseCode = "Z"
	AVSResponseCardholderNameIncorrectPostalMatch          AVSResponseCode = "AE1"
	AVSResponseCardholderNameIncorrectStreetAndPostalMatch AVSResponseCode = "AE2"
	AVSResponseCardholderNameIncorrectStreetMatch          AVSResponseCode = "AE3"
	AVSResponseCardholderNameMatch                         AVSResponseCode = "AE4"
	AVSResponseCardholderNameAndPostalMatch                AVSResponseCode = "AE5"
	AVSResponseCardholderNameAndStreetAndPostalMatch       AVSResponseCode = "AE6"
	AVSResponseCardholderNameAndStreetMatch                AVSResponseCode = "AE7"
)

type BalanceTransferRequest struct {
	Source                 string
	Destination            string
	Amount                 int64
	MerchantOrderReference string
	TransferType           *string
	IdempotencyKey         *string
}

// BalanceTransferResponse indicating a successful balance transfers properties
type BalanceTransferResponse struct {
	Success    bool
	ErrorCode  *string
	TransferID *string
}

package promptpay

import (
	"fmt"
	"regexp"

	"github.com/snksoft/crc"
)

// qr consts
const (
	idPayloadFormat                            = "00"
	idPOIMethod                                = "01"
	idMerchantInformationBot                   = "29"
	idTransactionCurrency                      = "53"
	idTransactionAmount                        = "54"
	idCountryCode                              = "58"
	idCRC                                      = "63"
	payloadFormatEMVQRCPSMerchantPresentedMode = "01"
	poiMethodStatic                            = "11"
	poiMethodDynamic                           = "12"
	merchantInformationTemplateIDGUID          = "00"
	botIDMerchantPhoneNumber                   = "01"
	botIDMerchantTaxID                         = "02"
	botIDMerchantEWalletID                     = "03"
	guidPromptpay                              = "A000000677010111"
	transactionCurrencyTHB                     = "764"
	countryCodeTH                              = "TH"
)

// regexp
var (
	reSanitizeTarget = regexp.MustCompile(`[^0-9]`)
	reFormatTarget   = regexp.MustCompile(`^0`)
)

func f(id, value string) string {
	return fmt.Sprintf("%s%02d%s", id, len(value), value)
}

func sanitizeTarget(id string) string {
	return reSanitizeTarget.ReplaceAllString(id, "")
}

func formatTarget(id string) string {
	numbers := sanitizeTarget(id)
	if len(numbers) >= 13 {
		return numbers
	}
	countryCoded := reFormatTarget.ReplaceAllString(id, "66")
	return fmt.Sprintf("%013s", countryCoded)
}

func formatAmount(amount float32) string {
	return fmt.Sprintf("%.2f", amount)
}

func formatCRC(crcValue uint64) string {
	return fmt.Sprintf("%04X", crcValue)
}

// Generate generates promptpay payload
func Generate(target string, amount float32) string {
	target = sanitizeTarget(target)
	var targetType string
	switch {
	case len(target) >= 15:
		targetType = botIDMerchantEWalletID
	case len(target) >= 13:
		targetType = botIDMerchantTaxID
	default:
		targetType = botIDMerchantPhoneNumber
	}

	data := ""
	data += f(idPayloadFormat, payloadFormatEMVQRCPSMerchantPresentedMode)
	if amount != 0 {
		data += f(idPOIMethod, poiMethodDynamic)
	} else {
		data += f(idPOIMethod, poiMethodStatic)
	}
	merchantInfo := f(merchantInformationTemplateIDGUID, guidPromptpay) + f(targetType, formatTarget(target))
	data += f(idMerchantInformationBot, merchantInfo)
	data += f(idCountryCode, countryCodeTH)
	data += f(idTransactionCurrency, transactionCurrencyTHB)
	data += f(idPayloadFormat, payloadFormatEMVQRCPSMerchantPresentedMode)
	if amount != 0 {
		data += f(idTransactionAmount, formatAmount(amount))
	}

	dataToCRC := fmt.Sprintf("%s%s%s", data, idCRC, "04")
	crcValue := crc.CalculateCRC(crc.CCITT, []byte(dataToCRC))
	data += f(idCRC, formatCRC(crcValue))
	return data
}

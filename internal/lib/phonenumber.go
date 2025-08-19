package lib

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

var ErrUnrecognizedPhoneFormat = fmt.Errorf("unrecognized phone format")

func ValidatePhone(phone string) (string, error) {
	phoneNum, err := phonenumbers.Parse(phone, "NG")
	if err != nil {
		return "", ErrUnrecognizedPhoneFormat
	}

	if !phonenumbers.IsValidNumber(phoneNum) {
		return "", ErrUnrecognizedPhoneFormat
	}

	return phonenumbers.Format(phoneNum, phonenumbers.E164), nil
}

package validate

import (
	"errors"
	"net/mail"
	"regexp"
)

func Email(raw string) error {
	_, err := mail.ParseAddress(raw)
	return err
}

func CitizenId(raw, country string) error {
	var p string

	switch country {
	case "MX":
		p = `^[A-Z]{4}[0-9]{6}[HM]{1}[A-Z]{2}[BCDFGHJKLMNPQRSTVWXYZ]{3}[0-9A-Z]{2}$`
	case "US":
		p = `^\d{3}-?\d{2}-?\d{4}$`
	default:
		return errors.New("invalid country code")
	}

	re := regexp.MustCompile(p)
	if !re.MatchString(raw) {
		return errors.New("invalid citizen id")
	}
	return nil
}

func TaxId(raw, country string) error {
	var p string
	switch country {
	case "MX":
		p = `^[A-ZÑ&]{3,4}\d{6}[A-Z0-9]{3}$`
	case "US":
		p = `^\d{2}-?\d{7}$`
	default:
		return errors.New("invalid country code")
	}

	re := regexp.MustCompile(p)
	if !re.MatchString(raw) {
		return errors.New("invalid tax id")
	}

	return nil
}

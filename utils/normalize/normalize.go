package normalize

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/nyaruka/phonenumbers"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

func String(s *string) *string {
	if s == nil {
		return nil
	}
	normalized := norm.NFC.String(strings.TrimSpace(*s))
	c := cases.Title(language.Und)
	result := c.String(normalized)
	return &result
}

func APIAction(query url.Values) (string, error) {
	queryAction := query["action"]
	var action string
	if len(queryAction) > 0 {
		action = queryAction[0]
	}
	if action == "" {
		return "", errors.New("no action")
	}
	return strings.ToUpper(action), nil
}

func Period(value string) (string, error) {
	if value == "" {
		return "", errors.New("no action")
	}
	return strings.ToUpper(value), nil
}

func Time(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func Date(raw string) (string, error) {
	const ISODate = "2006-01-02"
	parsed, err := time.Parse(ISODate, raw)
	if err != nil {
		return "", err
	}
	return parsed.Format(ISODate), nil
}

func Phone(raw string, defaultRegion string) (string, error) {
	num, err := phonenumbers.Parse(raw, defaultRegion)
	if err != nil {
		return "", err
	}

	if !phonenumbers.IsValidNumber(num) {
		return "", errors.New("invalid phone number")
	}

	return phonenumbers.Format(num, phonenumbers.E164), nil
}

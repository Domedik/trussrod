package notifications

import (
	"context"
	"encoding/json"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Twilio struct {
	client     *twilio.RestClient
	fromNumber string
}

func NewTwilioClient(accountSid, authToken, fromNumber string) *Twilio {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	return &Twilio{
		client:     client,
		fromNumber: fromNumber,
	}
}

func (t *Twilio) SendWhatsAppMessage(ctx context.Context, notification *WhatsAppMessage) error {
	variables, err := json.Marshal(notification.Variables)
	if err != nil {
		return err
	}

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(notification.To)
	params.SetFrom(t.fromNumber)
	params.SetContentSid(notification.TemplateID)
	params.SetContentVariables(string(variables))

	_, err = t.client.Api.CreateMessage(params)
	if err != nil {
		return err
	}

	return nil
}

package notifications

import "context"

type WhatsAppMessage struct {
	To         string
	TemplateID string
	Variables  map[string]string
}

type NotificationService interface {
	SendWhatsAppMessage(ctx context.Context, notification *WhatsAppMessage) error
}

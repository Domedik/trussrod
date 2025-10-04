package identity

import (
	"context"
	"time"
)

type Credentials struct {
	AccessKey    string
	SecretKey    string
	Expiration   time.Time
	SessionToken string
}

type Manager interface {
	GetId(ctx context.Context, token string) (string, error)
	GetCredentials(ctx context.Context, token string) (*Credentials, error)
}

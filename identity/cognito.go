package identity

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
)

type Cognito struct {
	client       *cognitoidentity.Client
	userPool     string
	identityPool string
	url          string
}

func (c *Cognito) GetId(ctx context.Context, token string) (string, error) {
	i := &cognitoidentity.GetIdInput{
		IdentityPoolId: aws.String(c.identityPool),
		Logins: map[string]string{
			c.url: token,
		},
	}
	res, err := c.client.GetId(ctx, i)
	if err != nil {
		return "", err
	}
	return *res.IdentityId, nil
}

func (c *Cognito) GetCredentials(ctx context.Context, token string) (*Credentials, error) {
	id, err := c.GetId(ctx, token)
	if err != nil {
		return nil, err
	}

	i := &cognitoidentity.GetCredentialsForIdentityInput{
		IdentityId: aws.String(id),
		Logins: map[string]string{
			c.url: token,
		},
	}

	res, err := c.client.GetCredentialsForIdentity(ctx, i)
	if err != nil {
		return nil, err
	}

	return &Credentials{
		AccessKey:    *res.Credentials.AccessKeyId,
		SecretKey:    *res.Credentials.SecretKey,
		Expiration:   *res.Credentials.Expiration,
		SessionToken: *res.Credentials.SessionToken,
	}, nil
}

func NewCognitoClient(url, identityPool string) (*Cognito, error) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
	)
	if err != nil {
		return nil, err
	}

	client := cognitoidentity.NewFromConfig(cfg)
	manager := &Cognito{
		client:       client,
		identityPool: identityPool,
		url:          url,
	}
	return manager, nil
}

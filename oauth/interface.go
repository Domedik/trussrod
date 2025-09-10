package oauth

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
)

type Client interface {
	Login(ctx context.Context, username, password string) (*types.AuthenticationResultType, error)
	RequestResetPassword(ctx context.Context, email string) (*cognitoidentityprovider.AdminResetUserPasswordOutput, error)
	ConfirmResetPassword(ctx context.Context, email, code, newPassword string) error
	ConfirmUserSignup(ctx context.Context, email, code string) error
	ValidateToken(ctx context.Context, jwt string) (jwt.Claims, error)
	Close() error
	Ping(ctx context.Context) error
}

package keys

import (
	"context"
	"encoding/base64"

	"github.com/Domedik/trussrod/utils/encryption"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

type KMS struct {
	client *kms.Client
	keyId  string
}

func (k *KMS) Decrypt(ctx context.Context, target []byte) ([]byte, error) {
	blob, err := base64.StdEncoding.DecodeString(string(target))
	if err != nil {
		return nil, err
	}
	input := &kms.DecryptInput{
		CiphertextBlob: blob,
	}

	decrypted, err := k.client.Decrypt(ctx, input)
	if err != nil {
		return nil, err
	}
	return decrypted.Plaintext, nil
}

func (k *KMS) CreateDEK(ctx context.Context) ([]byte, []byte, error) {
	input := &kms.GenerateDataKeyInput{
		KeyId:   aws.String(k.keyId),
		KeySpec: "AES_256",
	}
	out, err := k.client.GenerateDataKey(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	return out.Plaintext, out.CiphertextBlob, nil
}

func (k *KMS) Sign(ctx context.Context, params *SignInput) ([]byte, error) {
	hash := encryption.GetSHA256(params.Message)
	input := &kms.SignInput{
		KeyId:            aws.String(params.ARN),
		Message:          hash,
		MessageType:      types.MessageTypeDigest,
		SigningAlgorithm: types.SigningAlgorithmSpecRsassaPssSha256,
		// EncryptionContext: map[string]string{
		//     "doctor_id": "id",
		//     "note_id": "note-id",
		// },
	}

	result, err := k.client.Sign(ctx, input)
	if err != nil {
		return nil, err
	}

	return result.Signature, nil
}

func (k *KMS) Verify(ctx context.Context, params *VerifyInput) (bool, error) {
	hash := encryption.GetSHA256(params.Message)
	input := &kms.VerifyInput{
		KeyId:            aws.String(params.ARN),
		Message:          hash,
		Signature:        params.Signature,
		MessageType:      types.MessageTypeDigest,
		SigningAlgorithm: types.SigningAlgorithmSpecRsassaPssSha256,
	}

	result, err := k.client.Verify(ctx, input)
	if err != nil {
		return false, err
	}

	return result.SignatureValid, nil
}

func NewKMSClient(key string) (*KMS, error) {
	conf, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return &KMS{client: kms.NewFromConfig(conf), keyId: key}, nil
}

package Auth

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jonathanpatta/apartmentservices/Settings"
)

type FirebaseAuthService struct {
	db               *dynamodb.Client
	dynamodbSettings *Settings.DynamoDbSettings
	auth             *Settings.FirebaseAuthSettings
}

func NewFirebaseAuthService(settings *Settings.Settings) (*FirebaseAuthService, error) {
	return &FirebaseAuthService{
		db:               settings.Dynamo.Cli,
		dynamodbSettings: settings.Dynamo,
		auth:             settings.FirebaseAuth,
	}, nil
}

func (s *FirebaseAuthService) SetClaims(id string, claims map[string]interface{}) error {
	err := s.auth.Auth.SetCustomUserClaims(context.Background(), id, claims)
	if err != nil {
		return err
	}
	return nil
}

func (s *FirebaseAuthService) GetProviderId(id string) (string, error) {
	user, err := s.auth.Auth.GetUser(context.Background(), id)
	if err != nil {
		return "", err
	}
	providerId := user.CustomClaims["providerId"]
	if str, ok := providerId.(string); ok {
		if str != "" {
			return str, nil
		} else {
			return "", errors.New("providerId empty")
		}
	} else {
		return "", errors.New("claims not a string")
	}
	return "", nil
}

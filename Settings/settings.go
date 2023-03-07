package Settings

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/jonathanpatta/apartmentservices/Middleware"
	"google.golang.org/api/option"
	"io/ioutil"
	"os"
)

type Settings struct {
	Dynamo            *DynamoDbSettings
	FirebaseAuth      *FirebaseAuthSettings
	Region            string
	AwsCfg            aws.Config
	MiddlewareService *Middleware.MiddlwareService
}

func NewSettings() (*Settings, error) {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("could not read from .env file")
	}
	tableName := os.Getenv("DYNAMO_TABLE_NAME")
	region := os.Getenv("AWS_REGION_CODE")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	firebaseAuthSettings, err := NewFirebaseAuthSettings(s3Client)
	if err != nil {
		return nil, err
	}

	dynoDbSettings, err := NewDynamoDbSettings(cfg, tableName)
	if err != nil {
		return nil, err
	}

	middlewareService, err := Middleware.NewMiddlwareService(firebaseAuthSettings.Auth)
	if err != nil {
		return nil, err
	}

	return &Settings{
		Dynamo:            dynoDbSettings,
		FirebaseAuth:      firebaseAuthSettings,
		MiddlewareService: middlewareService,
		AwsCfg:            cfg,
	}, nil
}

type DynamoDbSettings struct {
	TableName *string
	Cli       *dynamodb.Client
}

func NewDynamoDbSettings(cfg aws.Config, TableName string) (*DynamoDbSettings, error) {
	dynamoDbCli := dynamodb.NewFromConfig(cfg)
	return &DynamoDbSettings{
		TableName: aws.String(TableName),
		Cli:       dynamoDbCli,
	}, nil
}

type FirebaseAuthSettings struct {
	App  *firebase.App
	Auth *auth.Client
}

func NewFirebaseAuthSettings(s3Client *s3.Client) (*FirebaseAuthSettings, error) {
	result, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("SECRETS_BUCKET")),
		Key:    aws.String(os.Getenv("FIREBASE_AUTH_SECRETS_FILENAME")),
	})
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	opt := option.WithCredentialsJSON(body)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}
	auth, err := app.Auth(context.Background())
	if err != nil {
		panic("Firebase load error")
	}
	return &FirebaseAuthSettings{
		App:  app,
		Auth: auth,
	}, nil
}

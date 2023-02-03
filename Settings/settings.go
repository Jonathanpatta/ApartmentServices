package Settings

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Settings struct {
	Dynamo *DynamoDbSettings
	Region string
	AwsCfg aws.Config
}

func NewSettings() (*Settings, error) {
	region := "ap-south-1"
	tableName := "apartment-services"

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	dynoDbSettings, err := NewDynamoDbSettings(cfg, tableName)
	if err != nil {
		return nil, err
	}

	return &Settings{
		Dynamo: dynoDbSettings,
		AwsCfg: cfg,
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

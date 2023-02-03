package Consumers

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/jonathanpatta/apartmentservices/Settings"
	"time"
)

var ConsumerKeyVal = "CONSUMER#"

type Meta struct {
	PK           string `json:"pk,omitempty"`
	SK           string `json:"sk,omitempty"`
	CreatedAt    int64  `json:"created_at,omitempty"`
	LastModified int64  `json:"last_modified,omitempty"`
}
type Consumer struct {
	Meta
	Id     string
	UserId string `json:"user_id,omitempty"`
}

type ConsumerService struct {
	db               *dynamodb.Client
	dynamodbSettings *Settings.DynamoDbSettings
}

func NewConsumerService(settings *Settings.Settings) (*ConsumerService, error) {
	return &ConsumerService{
		db:               settings.Dynamo.Cli,
		dynamodbSettings: settings.Dynamo,
	}, nil
}

func (s *ConsumerService) Create(in *Consumer) (*Consumer, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	in.SK = ConsumerKeyVal + id.String()
	in.Id = id.String()
	in.PK = ConsumerKeyVal
	in.CreatedAt = now
	in.LastModified = now

	item, err := attributevalue.MarshalMap(in)
	if err != nil {
		return nil, err
	}

	_, err = s.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:      item,
		TableName: s.dynamodbSettings.TableName,
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(item, in)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func (s *ConsumerService) Read(consumerId string) (*Consumer, error) {

	keyFilter := expression.Key("PK").Equal(expression.Value(ConsumerKeyVal)).
		And(expression.Key("SK").Equal(expression.Value(consumerId)))

	expr, err := expression.NewBuilder().WithKeyCondition(keyFilter).Build()
	if err != nil {
		return nil, err
	}

	out, err := s.db.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 s.dynamodbSettings.TableName,
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
	})
	if err != nil {
		return nil, err
	}

	var data []Consumer
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return nil, err
	}
	if len(data) != 1 {
		return nil, errors.New(fmt.Sprintf("length of return items not 1 got %v", len(data)))
	}

	return &data[0], nil
}

func (s *ConsumerService) Update(in *Consumer) (*Consumer, error) {

	consumer, err := s.Read(in.SK)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	consumer.LastModified = now
	consumer.UserId = in.UserId

	item, err := attributevalue.MarshalMap(consumer)
	if err != nil {
		return nil, err
	}

	_, err = s.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:      item,
		TableName: s.dynamodbSettings.TableName,
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(item, in)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func (s *ConsumerService) List() ([]*Consumer, error) {

	keyFilter := expression.Key("PK").Equal(expression.Value(ConsumerKeyVal))

	expr, err := expression.NewBuilder().WithKeyCondition(keyFilter).Build()
	if err != nil {
		return nil, err
	}

	out, err := s.db.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 s.dynamodbSettings.TableName,
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
	})
	if err != nil {
		return nil, err
	}

	var data []*Consumer
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ConsumerService) Delete(consumerId string) (*Consumer, error) {
	return nil, nil
}

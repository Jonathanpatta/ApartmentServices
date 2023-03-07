package Consumers

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jonathanpatta/apartmentservices/Settings"
	"github.com/jonathanpatta/apartmentservices/Utils"
)

const ConsumerPrefix = "CONSUMER#"

type Consumer struct {
	Utils.Meta
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
	err := in.New(ConsumerPrefix, "")
	if err != nil {
		return nil, err
	}

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

func (s *ConsumerService) CreateOrGet(in *Consumer) (*Consumer, error) {
	err := in.New(ConsumerPrefix, "")
	if err != nil {
		return nil, err
	}

	userIdConsumer, err := s.ReadFromUserId(in.UserId)
	if err == nil && userIdConsumer != nil && userIdConsumer.UserId == in.UserId {
		return userIdConsumer, nil
	}

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

func (s *ConsumerService) ReadFromUserId(userId string) (*Consumer, error) {
	keyFilter := expression.Key("PK").Equal(expression.Value(ConsumerPrefix))

	filterExpression := expression.Name("UserId").Equal(expression.Value(userId))

	expr, err := expression.NewBuilder().WithKeyCondition(keyFilter).WithFilter(filterExpression).Build()
	if err != nil {
		return nil, err
	}

	out, err := s.db.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 s.dynamodbSettings.TableName,
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
		ConsistentRead:            aws.Bool(false),
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

func (s *ConsumerService) Read(consumerId string) (*Consumer, error) {

	keyFilter := expression.Key("PK").Equal(expression.Value(ConsumerPrefix)).
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

	consumer.SetLastModifiedNow()
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

	keyFilter := expression.Key("PK").Equal(expression.Value(ConsumerPrefix))

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

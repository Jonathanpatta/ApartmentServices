package Producers

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jonathanpatta/apartmentservices/Services"
	"github.com/jonathanpatta/apartmentservices/Settings"
	"github.com/jonathanpatta/apartmentservices/Utils"
)

const ProducerPrefix = "PRODUCER#"

type Producer struct {
	Utils.Meta
	ApartmentNumber string `json:"apartment_number"`
}

type ProducerService struct {
	db               *dynamodb.Client
	dynamodbSettings *Settings.DynamoDbSettings
	servicesCli      *Services.ServiceService
}

func NewProducerService(settings *Settings.Settings) (*ProducerService, error) {
	servicesCli, err := Services.NewServiceService(settings)
	if err != nil {
		return nil, err
	}

	return &ProducerService{
		db:               settings.Dynamo.Cli,
		dynamodbSettings: settings.Dynamo,
		servicesCli:      servicesCli,
	}, nil
}

func (s *ProducerService) Create(in *Producer) (*Producer, error) {
	err := in.New("", ProducerPrefix)
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

func (s *ProducerService) Read(producerId string) (*Producer, error) {
	keyFilter := expression.Key("PK").Equal(expression.Value(ProducerPrefix)).
		And(expression.Key("SK").Equal(expression.Value(producerId)))

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
		ConsistentRead:            aws.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	var data []Producer
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return nil, err
	}
	if len(data) != 1 {
		return nil, errors.New(fmt.Sprintf("length of return items not 1 got %v", len(data)))
	}

	return &data[0], nil
}

func (s *ProducerService) Update(in *Producer) (*Producer, error) {
	producer, err := s.Read(in.SK)
	if err != nil {
		return nil, err
	}

	producer.SetLastModifiedNow()

	producer.ApartmentNumber = in.ApartmentNumber

	item, err := attributevalue.MarshalMap(producer)
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

func (s *ProducerService) List() ([]*Producer, error) {

	keyFilter := expression.Key("PK").Equal(expression.Value(ProducerPrefix))

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
		ConsistentRead:            aws.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	var data []*Producer
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ProducerService) Delete(producerId string) (*Producer, error) {
	return nil, nil
}

type AddServiceInput struct {
	producerId string
	service    *Services.Service
}

type RemoveServiceInput struct {
	producerId string
	serviceId  string
}

func (s *ProducerService) GetServices(producerId string) ([]*Services.Service, error) {
	producer, err := s.Read(producerId)
	if err != nil {
		return nil, err
	}

	keyFilter := expression.Key("PK").Equal(expression.Value(Services.ServicePrefix)).
		And(expression.Key("SK").BeginsWith(producer.SK))

	filter := expression.Name("IsDeleted").NotEqual(expression.Value(true))

	expr, err := expression.NewBuilder().WithKeyCondition(keyFilter).WithFilter(filter).Build()
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

	var data []*Services.Service
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

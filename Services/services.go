package Services

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

const ServicePrefix = "SERVICE#"

type Service struct {
	Utils.Meta
	Name string `json:"name"`
}

const ProducerPrefix = "PRODUCER#"

type Producer struct {
	Utils.Meta
	ApartmentNumber string `json:"apartment_number"`
}

type ServiceService struct {
	db               *dynamodb.Client
	dynamodbSettings *Settings.DynamoDbSettings
}

func NewServiceService(settings *Settings.Settings) (*ServiceService, error) {
	return &ServiceService{
		db:               settings.Dynamo.Cli,
		dynamodbSettings: settings.Dynamo,
	}, nil
}

func (s *ServiceService) Create(producerId string, in *Service) (*Service, error) {
	if producerId == "" {
		return nil, errors.New("producer id required")
	}
	err := s.ProducerCheck(producerId)
	if err != nil {
		return nil, err
	}
	err = in.New(producerId, ServicePrefix)
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

func (s *ServiceService) Read(serviceId string) (*Service, error) {

	keyFilter := expression.Key("PK").Equal(expression.Value(ServicePrefix)).
		And(expression.Key("SK").Equal(expression.Value(serviceId)))

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

	var data []Service
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return nil, err
	}
	if len(data) != 1 {
		return nil, errors.New(fmt.Sprintf("length of return items not 1 got %v", len(data)))
	}

	return &data[0], nil
}

func (s *ServiceService) Update(in *Service) (*Service, error) {

	service, err := s.Read(in.SK)
	if err != nil {
		return nil, err
	}

	service.SetLastModifiedNow()
	service.Name = in.Name

	item, err := attributevalue.MarshalMap(service)
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

func (s *ServiceService) List() ([]*Service, error) {

	keyFilter := expression.Key("PK").Equal(expression.Value(ServicePrefix))

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

	var data []*Service
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ServiceService) Delete(serviceId string) (*Service, error) {
	return nil, nil
}

func (s *ServiceService) ProducerCheck(producerId string) error {
	keyFilter := expression.Key("PK").Equal(expression.Value(ProducerPrefix)).
		And(expression.Key("SK").Equal(expression.Value(producerId)))

	expr, err := expression.NewBuilder().WithKeyCondition(keyFilter).Build()
	if err != nil {
		return err
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
		return err
	}

	var data []Producer
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return err
	}
	if len(data) != 1 {
		return errors.New(fmt.Sprintf("length of return items not 1 got %v", len(data)))
	}
	return nil
}

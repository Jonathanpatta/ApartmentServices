package Orders

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

type Order struct {
	Utils.Meta
	ItemId string `json:"item_id,omitempty"`
}

type OrderService struct {
	db               *dynamodb.Client
	dynamodbSettings *Settings.DynamoDbSettings
}

const OrderPrefix = "ORDER#"

func NewOrderService(settings *Settings.Settings) (*OrderService, error) {
	return &OrderService{
		db:               settings.Dynamo.Cli,
		dynamodbSettings: settings.Dynamo,
	}, nil
}

func (s *OrderService) Create(consumerId string, in *Order) (*Order, error) {

	err := s.ConsumerCheck(consumerId)
	if err != nil {
		return nil, err
	}

	err = in.New(OrderPrefix, consumerId)
	if err != nil {
		return nil, err
	}

	order, err := attributevalue.MarshalMap(in)
	if err != nil {
		return nil, err
	}

	_, err = s.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:      order,
		TableName: s.dynamodbSettings.TableName,
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(order, in)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func (s *OrderService) Read(orderId string) (*Order, error) {

	keyFilter := expression.Key("PK").Equal(expression.Value(OrderPrefix)).
		And(expression.Key("SK").Equal(expression.Value(orderId)))

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

	var data []Order
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return nil, err
	}
	if len(data) != 1 {
		return nil, errors.New(fmt.Sprintf("length of return orders not 1 got %v", len(data)))
	}

	return &data[0], nil
}

func (s *OrderService) Update(in *Order) (*Order, error) {

	prevOrder, err := s.Read(in.SK)
	if err != nil {
		return nil, err
	}

	prevOrder.ItemId = in.ItemId
	prevOrder.SetLastModifiedNow()

	order, err := attributevalue.MarshalMap(prevOrder)
	if err != nil {
		return nil, err
	}

	_, err = s.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:      order,
		TableName: s.dynamodbSettings.TableName,
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(order, in)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func (s *OrderService) List() ([]*Order, error) {

	keyFilter := expression.Key("PK").Equal(expression.Value(OrderPrefix))

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

	var data []*Order
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *OrderService) Delete(orderId string) (*Order, error) {
	return nil, nil
}

func (s *OrderService) ConsumerCheck(consumerId string) error {
	keyFilter := expression.Key("PK").Equal(expression.Value(ConsumerPrefix)).
		And(expression.Key("SK").Equal(expression.Value(consumerId)))

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

	var data []Consumer
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return err
	}
	if len(data) != 1 {
		return errors.New(fmt.Sprintf("length of return orders not 1 got %v", len(data)))
	}
	return nil
}

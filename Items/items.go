package Items

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

const ItemPrefix = "ITEM#"

type Item struct {
	Utils.Meta
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	ImageUrls   []string `json:"image_urls,omitempty"`
	Price       int64    `json:"price,omitempty"`
}

type ItemService struct {
	db               *dynamodb.Client
	dynamodbSettings *Settings.DynamoDbSettings
}

func NewItemService(settings *Settings.Settings) (*ItemService, error) {
	return &ItemService{
		db:               settings.Dynamo.Cli,
		dynamodbSettings: settings.Dynamo,
	}, nil
}

func (s *ItemService) Create(serviceId string, in *Item) (*Item, error) {

	err := s.ServiceCheck(serviceId)
	if err != nil {
		return nil, err
	}

	err = in.New(ItemPrefix, serviceId)
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

func (s *ItemService) Read(itemId string) (*Item, error) {

	keyFilter := expression.Key("PK").Equal(expression.Value(ItemPrefix)).
		And(expression.Key("SK").Equal(expression.Value(itemId)))

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

	var data []Item
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return nil, err
	}
	if len(data) != 1 {
		return nil, errors.New(fmt.Sprintf("length of return items not 1 got %v", len(data)))
	}

	return &data[0], nil
}

func (s *ItemService) Update(in *Item) (*Item, error) {

	prevItem, err := s.Read(in.SK)
	if err != nil {
		return nil, err
	}

	prevItem.Name = in.Name
	prevItem.ImageUrls = in.ImageUrls
	prevItem.Price = in.Price
	prevItem.SetLastModifiedNow()

	item, err := attributevalue.MarshalMap(prevItem)
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

func (s *ItemService) List() ([]*Item, error) {

	keyFilter := expression.Key("PK").Equal(expression.Value(ItemPrefix))

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

	var data []*Item
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ItemService) Delete(itemId string) (*Item, error) {
	return nil, nil
}

func (s *ItemService) ServiceCheck(serviceId string) error {
	keyFilter := expression.Key("PK").Equal(expression.Value(ServicePrefix)).
		And(expression.Key("SK").Equal(expression.Value(serviceId)))

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

	var data []Service
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return err
	}
	if len(data) != 1 {
		return errors.New(fmt.Sprintf("length of return items not 1 got %v", len(data)))
	}
	return nil
}

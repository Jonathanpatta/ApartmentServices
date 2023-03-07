package Subscriptions

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

type Subscription struct {
	Utils.Meta
	ItemId        string `json:"item_id,omitempty"`
	ItemName      string `json:"item_name,omitempty"`
	Note          string `json:"note,omitempty"`
	RecurringType string `json:"recurring_type,omitempty"`
	Cancelled     bool   `json:"cancelled,omitempty"`

	CreatedByUserId      string `json:"created_by_user_id,omitempty"`
	CreatedByName        string `json:"created_by_name,omitempty"`
	CreatedByUserEmail   string `json:"created_by_user_email,omitempty"`
	CreatedByUserPicture string `json:"created_by_user_picture,omitempty"`
}

type SubscriptionService struct {
	db               *dynamodb.Client
	dynamodbSettings *Settings.DynamoDbSettings
}

const SubscriptionPrefix = "SUBSCRIPTION#"

func NewSubscriptionService(settings *Settings.Settings) (*SubscriptionService, error) {
	return &SubscriptionService{
		db:               settings.Dynamo.Cli,
		dynamodbSettings: settings.Dynamo,
	}, nil
}

func (s *SubscriptionService) Create(consumerId string, in *Subscription) (*Subscription, error) {

	err := s.ConsumerCheck(consumerId)
	if err != nil {
		return nil, err
	}

	err = in.New(SubscriptionPrefix, consumerId)
	if err != nil {
		return nil, err
	}

	subscription, err := attributevalue.MarshalMap(in)
	if err != nil {
		return nil, err
	}

	_, err = s.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:      subscription,
		TableName: s.dynamodbSettings.TableName,
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(subscription, in)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func (s *SubscriptionService) Read(subscriptionId string) (*Subscription, error) {

	keyFilter := expression.Key("PK").Equal(expression.Value(SubscriptionPrefix)).
		And(expression.Key("SK").Equal(expression.Value(subscriptionId)))

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

	var data []Subscription
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return nil, err
	}
	if len(data) != 1 {
		return nil, errors.New(fmt.Sprintf("length of return subscriptions not 1 got %v", len(data)))
	}

	return &data[0], nil
}

func (s *SubscriptionService) Update(in *Subscription) (*Subscription, error) {

	prevSubscription, err := s.Read(in.SK)
	if err != nil {
		return nil, err
	}

	prevSubscription.ItemId = in.ItemId
	prevSubscription.SetLastModifiedNow()

	subscription, err := attributevalue.MarshalMap(prevSubscription)
	if err != nil {
		return nil, err
	}

	_, err = s.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:      subscription,
		TableName: s.dynamodbSettings.TableName,
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(subscription, in)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func (s *SubscriptionService) List() ([]*Subscription, error) {

	keyFilter := expression.Key("PK").Equal(expression.Value(SubscriptionPrefix))

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

	var data []*Subscription
	err = attributevalue.UnmarshalListOfMaps(out.Items, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *SubscriptionService) Delete(subscriptionId string) (*Subscription, error) {
	return nil, nil
}

func (s *SubscriptionService) ConsumerCheck(consumerId string) error {
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
		return errors.New(fmt.Sprintf("length of return subscriptions not 1 got %v", len(data)))
	}
	return nil
}

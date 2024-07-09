package outbox

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type DynamoOutboxRepository struct {
	dynamoClient *dynamodb.DynamoDB
	tableName    string
}

func NewDynamoOutboxRepository(dynamoClient *dynamodb.DynamoDB, tableName string) Repository {
	return &DynamoOutboxRepository{dynamoClient: dynamoClient, tableName: tableName}
}

func (r *DynamoOutboxRepository) Store(ctx context.Context, entry *Entry) error {
	item, err := dynamodbattribute.MarshalMap(entry)
	if err != nil {
		return err
	}
	_, err = r.dynamoClient.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(r.tableName),
	})
	return err
}

func (r *DynamoOutboxRepository) MarkAsProcessed(ctx context.Context, entry *Entry) error {
	return r.changeStatus(ctx, entry, "PROCESSED")
}

func (r *DynamoOutboxRepository) MarkAsError(ctx context.Context, entry *Entry) error {
	return r.changeStatus(ctx, entry, "ERROR")
}

func (r *DynamoOutboxRepository) changeStatus(ctx context.Context, entry *Entry, status string) error {
	key, err := dynamodbattribute.MarshalMap(map[string]string{"id": entry.Id})
	if err != nil {
		return err
	}
	update := expression.Set(expression.Name("status"), expression.Value(status))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}
	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(r.tableName),
		Key:                       key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}
	_, err = r.dynamoClient.UpdateItemWithContext(ctx, input)
	return err
}

func (r *DynamoOutboxRepository) Get(ctx context.Context, id string) (*Entry, error) {
	key, err := dynamodbattribute.MarshalMap(map[string]string{"id": id})
	if err != nil {
		return nil, err
	}
	item, err := r.dynamoClient.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	var entry Entry
	err = dynamodbattribute.UnmarshalMap(item.Item, &entry)
	if err != nil {
		return nil, err
	}
	if entry.Id == "" {
		return nil, nil
	}
	return &entry, nil
}

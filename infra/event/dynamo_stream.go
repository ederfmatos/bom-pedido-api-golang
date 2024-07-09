package event

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodbstreams"
	"sync"
	"time"
)

type DynamoStream struct {
	dynamoStreamClient *dynamodbstreams.DynamoDBStreams
	awsSession         *session.Session
	tableName          string
	dynamoDB           *dynamodb.DynamoDB
}

func NewDynamoStream(awsSession *session.Session, tableName string, dynamoDB *dynamodb.DynamoDB) *DynamoStream {
	return &DynamoStream{
		dynamoStreamClient: dynamodbstreams.New(awsSession),
		dynamoDB:           dynamoDB,
		awsSession:         awsSession,
		tableName:          tableName,
	}
}

func (stream *DynamoStream) getStreamArn() (string, error) {
	result, err := stream.dynamoDB.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(stream.tableName)})
	if err != nil {
		return "", err
	}
	if result.Table.StreamSpecification != nil && *result.Table.StreamSpecification.StreamEnabled {
		return *result.Table.LatestStreamArn, nil
	}
	return "", fmt.Errorf("streams not enabled for table %s", stream.tableName)
}

func (stream *DynamoStream) FetchEvents() (chan *dynamodbstreams.Record, error) {
	streamArn, err := stream.getStreamArn()
	if err != nil {
		return nil, err
	}
	events := make(chan *dynamodbstreams.Record)
	describeStreamInput := &dynamodbstreams.DescribeStreamInput{StreamArn: aws.String(streamArn)}
	describeStreamOutput, err := stream.dynamoStreamClient.DescribeStream(describeStreamInput)
	if err != nil {
		return nil, err
	}
	for _, shard := range describeStreamOutput.StreamDescription.Shards {
		go stream.processShard(*shard.ShardId, events, streamArn)
	}
	return events, nil
}

func (stream *DynamoStream) processShard(shardID string, events chan<- *dynamodbstreams.Record, streamArn string) {
	shardIteratorInput := &dynamodbstreams.GetShardIteratorInput{
		StreamArn:         aws.String(streamArn),
		ShardId:           aws.String(shardID),
		ShardIteratorType: aws.String(dynamodbstreams.ShardIteratorTypeTrimHorizon),
	}

	shardIteratorOutput, err := stream.dynamoStreamClient.GetShardIterator(shardIteratorInput)
	if err != nil {
		return
	}

	ShardIterator := shardIteratorOutput.ShardIterator
	for {
		getRecordsInput := &dynamodbstreams.GetRecordsInput{ShardIterator: ShardIterator}
		getRecordsOutput, err := stream.dynamoStreamClient.GetRecords(getRecordsInput)
		if err != nil {
			continue
		}

		waitGroup := sync.WaitGroup{}
		for _, record := range getRecordsOutput.Records {
			if *record.EventName == "INSERT" {
				go func() {
					waitGroup.Add(1)
					events <- record
					waitGroup.Done()
				}()
			}
		}
		waitGroup.Wait()
		ShardIterator = getRecordsOutput.NextShardIterator

		time.Sleep(1 * time.Second)
	}
}

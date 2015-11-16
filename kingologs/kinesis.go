package kingologs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

// KinesisRelay is the main kinesis struct.
type KinesisRelay struct {
	logger    Logger
	awsConfig aws.Config
	config    ConfigValues
	ksvc      *kinesis.Kinesis
	Pipe      chan string
}

// NewKinesisRelay gets the relay started.
func NewKinesisRelay(l Logger, c ConfigValues) *KinesisRelay {
	l.Info.Println("Beginning to start kinesis")

	// New KinesisRelay setup.
	kr := new(KinesisRelay)
	kr.logger = l
	kr.config = c

	// Setup the AWS config.
	creds := credentials.NewEnvCredentials()
	kr.awsConfig = aws.Config{}
	kr.awsConfig.Credentials = creds

	// TODO - handle this via ENV variables as well.
	kr.awsConfig.Region = aws.String(kr.config.Kinesis.Region)

	// TODO configure the channel better in the future.
	kr.Pipe = make(chan string, 100)

	// Prep the stream.
	l.Trace.Printf("We want to use stream: %s", kr.config.Kinesis.StreamName)
	kr.createStreamIfNotExists(kr.config.Kinesis.StreamName)

	l.Info.Println("Done starting kinesis relay")
	return kr
}

// NewMessage puts a new message into the channel.
func (kr KinesisRelay) NewMessage(s string) {
	kr.logger.Trace.Printf("Adding to channel: %s", s)
	kr.Pipe <- s
}

// StartRelay gets new messages from the channel.
func (kr KinesisRelay) StartRelay() {
	kr.logger.Info.Println("Starting to watch channel")

	// Wait for new messages in channel.
	for true {
		select {
		case msg := <-kr.Pipe:
			kr.logger.Trace.Printf("Value in channel: %s", msg)
			kr.putRecord(msg)
		}
	}
}

// Make a stream if it does not yet exist.
func (kr KinesisRelay) createStreamIfNotExists(name string) {
	kr.logger.Info.Printf("Checking to see if we have a '%s' stream already", name)

	svc := kinesis.New(&kr.awsConfig)

	// TODO - do we need to paginate this?
	params := &kinesis.ListStreamsInput{
		ExclusiveStartStreamName: aws.String("StreamName"),
		//Limit: aws.Int64(1),
	}
	resp, err := svc.ListStreams(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		kr.logger.Error.Println(fmt.Sprintf("Error w/ ListStreams: %#v", err.Error()))
		return
	}

	kr.logger.Trace.Println(fmt.Sprintf("Kinesis Streams: %#v", resp))

	kr.logger.Info.Println(fmt.Sprintf("Checking to see if %s stream exists", name))

	// Loop through all streams to see if the one we want exists.
	streamExists := false
	for _, stream := range resp.StreamNames {
		if *stream == name {
			streamExists = true
			break
		}
	}

	if !streamExists {
		kr.logger.Info.Println("Stream did not exist, so we'll attempt to create it")
		kr.createStream(name)
	}
}

// Create a new stream.
func (kr KinesisRelay) createStream(name string) {
	kr.logger.Info.Printf("Attempting to create '%s' stream", name)

	svc := kinesis.New(&kr.awsConfig)

	params := &kinesis.CreateStreamInput{
		StreamName: aws.String(name),
		ShardCount: aws.Int64(1),
	}

	resp, err := svc.CreateStream(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		kr.logger.Error.Printf("Error w/ ListStreams: %#v", err.Error())
		return
	}

	kr.logger.Trace.Printf("CreateStream response:  %#v", resp)

	kr.logger.Info.Println("New stream created!")

}

// Put a single record.
func (kr KinesisRelay) putRecord(msg string) {
	kr.logger.Info.Println("Attempting to put record into stream")

	svc := kinesis.New(&kr.awsConfig)

	// TODO real partition key.
	params := &kinesis.PutRecordInput{
		StreamName:   aws.String(kr.config.Kinesis.StreamName),
		PartitionKey: aws.String("1"),
		Data:         []byte(msg),
	}

	resp, err := svc.PutRecord(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		kr.logger.Error.Printf("Error w/ PutRecord: %#v", err.Error())
		return
	}

	kr.logger.Trace.Printf("PutRecord response:  %#v", resp)

	kr.logger.Info.Println("Done putting record!")
}

// Put multiple records.
func (kr KinesisRelay) putRecords(messages []string, num int) {
	kr.logger.Info.Println("Attempting to put records into stream")

	svc := kinesis.New(&kr.awsConfig)

	records := make([]*kinesis.PutRecordsRequestEntry, num)
	for _, msg := range messages {
		records.Append(records, &kinesis.PutRecordsRequestEntry{
			Data:         []byte(msg),
			PartitionKey: aws.String("1"),
		})
	}

	// TODO real partition key.
	params := &kinesis.PutRecordsInput{
		StreamName: aws.String(kr.config.Kinesis.StreamName),
		Records:    records,
	}

	resp, err := svc.PutRecords(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		kr.logger.Error.Printf("Error w/ PutRecord: %#v", err.Error())
		return
	}

	kr.logger.Trace.Printf("PutRecord response:  %#v", resp)

	kr.logger.Info.Println("Done putting record!")
}

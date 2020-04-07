package squat

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Squat struct {
	client   *sqs.SQS
	queueUrl string
}

// begins reading from stdin
// TODO abstract this to accept the source as arg
func (sq *Squat) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
        err := sq.putrec(s)
        if err != nil {
            fmt.Fprintf(os.Stderr, "ERR: %s", err)
        }
	}
}

// put record to kinesis
func (sq *Squat) putrec(data string) error {
	// completely heenky dedup token creation
	tok := Int64ToStr(rand.Int63())
	dedup := "sqt" + tok + "-" + Int64ToStr(time.Now().UnixNano())
	send_params := &sqs.SendMessageInput{
		MessageBody:            aws.String(data),
		QueueUrl:               aws.String(sq.queueUrl), // Required
        // required for fifo queues
		MessageGroupId:         aws.String("squat"),
		MessageDeduplicationId: aws.String(dedup),
	}
	_,err := sq.client.SendMessage(send_params)
	return err
}

func NewSquat(cfg *RuntimeConfig) (*Squat, error) {
	awscfg := aws.Config{}
	if cfg.Region != "" {
		awscfg.Region = aws.String(cfg.Region)
	}
	awscfg.Credentials = credentials.NewSharedCredentials("", "default")
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            awscfg,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}
	return &Squat{
		client:   sqs.New(sess),
		queueUrl: cfg.QueueUrl,
	}, nil
}

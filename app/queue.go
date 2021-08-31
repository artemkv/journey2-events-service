package app

import (
	"encoding/json"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var AWS_SESSION_FRESHNESS_TIMEOUT = time.Duration(2) * time.Minute

var (
	mu           sync.Mutex       // guards session
	awsSession   *session.Session // do not use directly, use getSession()
	sessionStart time.Time
)

var _actionTopicArn string

func InitSNSConnection(actionTopicArn string) {
	_actionTopicArn = actionTopicArn
	refreshSession()
}

func refreshSession() {
	sessionStart = time.Now()
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	awsSession = sess
}

func getSession() *session.Session {
	sessionDuration := time.Since(sessionStart)
	if sessionDuration > AWS_SESSION_FRESHNESS_TIMEOUT {
		log.Printf("Refreshing AWS session")
		mu.Lock()
		refreshSession()
		mu.Unlock()
	}
	mu.Lock()
	sess := awsSession
	mu.Unlock()
	return sess
}

func EnqueueAction(action *actionOutgoingData) (string, error) {
	svc := sns.New(getSession())

	// serialize action
	bytes, err := json.Marshal(action)
	if err != nil {
		return "", err
	}
	msg := string(bytes)

	// publish
	result, err := svc.Publish(&sns.PublishInput{
		Message:  &msg,
		TopicArn: &_actionTopicArn,
	})

	// handle error
	if err != nil {
		return "", err
	}
	return *result.MessageId, nil
}

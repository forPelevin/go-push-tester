package go_push_tester

import (
	"context"
	"fmt"
	"log"
	"sync"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func NewFCMClient() (*messaging.Client, error) {
	opt := option.WithCredentialsFile("creds.json")

	app, err := firebase.NewApp(context.TODO(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("new firebase app: %w", err)
	}

	fcmCli, err := app.Messaging(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("messaging: %w", err)
	}

	return fcmCli, nil
}

type Sender struct {
	fcmClient *messaging.Client

	mu          sync.Mutex
	pushCounter int
}

func (s *Sender) SendPush(ctx context.Context, msg *messaging.Message) {
	response, err := s.fcmClient.Send(ctx, msg)

	log.Printf("response: %+v, err: %s \n", response, err)

	s.mu.Lock()
	s.pushCounter++
	s.mu.Unlock()
}

func (s *Sender) PushCount() int {
	return s.pushCounter
}

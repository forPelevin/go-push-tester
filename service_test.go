package go_push_tester

import (
	"context"
	"log"
	"testing"

	"firebase.google.com/go/messaging"
)

func TestSender_SendPush(t *testing.T) {
	fcmCli, err := NewFCMClient()
	if err != nil {
		t.Error(err)
	}

	s := &Sender{
		fcmClient: fcmCli,
	}

	s.SendPush(context.TODO(), testMessage())
}

func Benchmark_Service_SendPush(b *testing.B) {
	fcmCli, err := NewFCMClient()
	if err != nil {
		b.Error(err)
	}

	s := &Sender{
		fcmClient: fcmCli,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.SendPush(context.TODO(), testMessage())
		}
	})

	log.Printf("Sent %d push messages by service\n", s.PushCount())
}

func testMessage() *messaging.Message {
	return &messaging.Message{
		Notification: &messaging.Notification{
			Title: "A nice notification title",
			Body:  "A nice notification body",
		},
		Token: "client-push-token", // a token that you received from client
	}
}

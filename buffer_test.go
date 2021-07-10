package go_push_tester

import (
	"log"
	"testing"
	"time"

	"firebase.google.com/go/messaging"
)

func TestBuffer_SendPush(t *testing.T) {
	fcmCli, err := NewFCMClient()
	if err != nil {
		t.Error(err)
	}

	buf := &Buffer{
		fcmClient:        fcmCli,
		dispatchInterval: 1 * time.Millisecond,
		batchCh:          make(chan *messaging.Message),
	}

	buf.Run()
	buf.SendPush(testMessage())
	buf.Stop()
}

func Benchmark_Buffer_SendPush(b *testing.B) {
	fcmCli, err := NewFCMClient()
	if err != nil {
		b.Error(err)
	}

	buf := &Buffer{
		fcmClient:        fcmCli,
		dispatchInterval: 3 * time.Second,
		batchCh:          make(chan *messaging.Message),
	}

	buf.Run()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf.SendPush(testMessage())
		}
	})

	buf.Stop()

	log.Printf("Sent %d push messages by buffer\n", buf.PushCount())
}

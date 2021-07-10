package go_push_tester

import (
	"context"
	"log"
	"sync"
	"time"

	"firebase.google.com/go/messaging"
)

// Buffer batches all incoming push messages and send them periodically.
type Buffer struct {
	fcmClient *messaging.Client

	dispatchInterval time.Duration
	batchCh          chan *messaging.Message
	wg               sync.WaitGroup
	mu               sync.Mutex
	pushCounter      int
}

func (b *Buffer) SendPush(msg *messaging.Message) {
	b.batchCh <- msg
}

func (b *Buffer) sender() {
	defer b.wg.Done()

	// set your interval
	t := time.NewTicker(b.dispatchInterval)

	// we can send up to 500 messages per call to Firebase
	messages := make([]*messaging.Message, 0, 500)

	defer func() {
		t.Stop()

		// send all buffered messages before quit
		b.sendMessages(messages)

		log.Println("batch sender finished")
	}()

	for {
		select {
		case m, ok := <-b.batchCh:
			if !ok {
				return
			}

			messages = append(messages, m)
		case <-t.C:
			b.sendMessages(messages)
			messages = messages[:0]
		}
	}
}

func (b *Buffer) Run() {
	b.wg.Add(1)
	go b.sender()
}

func (b *Buffer) Stop() {
	close(b.batchCh)
	b.wg.Wait()
}

func (b *Buffer) sendMessages(messages []*messaging.Message) {
	if len(messages) == 0 {
		return
	}

	batchResp, err := b.fcmClient.SendAll(context.TODO(), messages)

	log.Printf("batch response: %+v, err: %s \n", batchResp, err)

	b.mu.Lock()
	b.pushCounter += len(messages)
	b.mu.Unlock()
}

func (b *Buffer) PushCount() int {
	return b.pushCounter
}

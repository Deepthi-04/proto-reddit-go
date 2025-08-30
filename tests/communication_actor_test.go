package tests

import (
	proto_actor "reddit-clone/core/proto_actors"
	"reddit-clone/schemas"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

func BenchmarkMessageActor(b *testing.B) {

	system := actor.NewActorSystem()
	messageActor := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return proto_actor.NewMessageManager()
	}))

	b.Run("SendMessage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := system.Root.RequestFuture(messageActor, &proto_actor.SendMessage{
				FromUserID: "userA",
				ToUserID:   "userB",
				Body:       "Benchmark message",
			}, 3*time.Second).Result()
			if err != nil {
				b.Fatalf("SendMessage failed during benchmark: %v", err)
			}
		}
	})

	res, err := system.Root.RequestFuture(messageActor, &proto_actor.SendMessage{
		FromUserID: "userA",
		ToUserID:   "userB",
		Body:       "Benchmark test message",
	}, 3*time.Second).Result()
	if err != nil {
		b.Fatalf("SendMessage setup (benchmark) failed: %v", err)
	}
	message, ok := res.(*schemas.Message)
	if !ok || message == nil || message.ID == "" {
		b.Fatalf("Invalid response for SendMessage setup: %v", res)
	}

	b.Run("FetchMessages", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := system.Root.RequestFuture(messageActor, &proto_actor.FetchMessages{
				UserID: "userA",
			}, 3*time.Second).Result()
			if err != nil {
				b.Fatalf("FetchMessages failed during benchmark: %v", err)
			}
		}
	})

	b.Run("RemoveMessage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {

			res, err := system.Root.RequestFuture(messageActor, &proto_actor.SendMessage{
				FromUserID: "userA",
				ToUserID:   "userB",
				Body:       "Benchmark delete message",
			}, 3*time.Second).Result()
			if err != nil {
				b.Fatalf("SendMessage failed during RemoveMessage benchmark: %v", err)
			}

			msg, ok := res.(*schemas.Message)
			if !ok || msg == nil || msg.ID == "" {
				b.Fatalf("Invalid response for SendMessage during RemoveMessage benchmark: %v", res)
			}

			_, err = system.Root.RequestFuture(messageActor, &proto_actor.RemoveMessage{
				MessageID: msg.ID,
			}, 3*time.Second).Result()
			if err != nil {
				b.Fatalf("RemoveMessage failed during benchmark: %v", err)
			}
		}
	})
}

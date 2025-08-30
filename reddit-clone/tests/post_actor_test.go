package tests

import (
	proto_actor "reddit-clone/core/proto_actors"
	"reddit-clone/schemas"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

func TestPostManager(t *testing.T) {
	system := actor.NewActorSystem()
	postManager := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return proto_actor.NewPostManager()
	}))

	res, err := system.Root.RequestFuture(postManager, &proto_actor.AddPost{
		ForumID:  "test-forum",
		AuthorID: "test-user",
		Text:     "Test content",
	}, 3*time.Second).Result()
	if err != nil {
		t.Fatalf("AddPost failed: %v", err)
	}

	post, ok := res.(*schemas.Post)
	if !ok || post == nil {
		t.Fatalf("Invalid response for AddPost: %v", res)
	}

	res, err = system.Root.RequestFuture(postManager, &proto_actor.RetrievePost{ContentID: post.ID}, 3*time.Second).Result()
	if err != nil {
		t.Fatalf("RetrievePost failed: %v", err)
	}

	retrievedPost, ok := res.(*schemas.Post)
	if !ok || retrievedPost == nil || retrievedPost.ID != post.ID {
		t.Fatalf("Invalid response for RetrievePost: %v", res)
	}

	res, err = system.Root.RequestFuture(postManager, &proto_actor.RetrieveAllPosts{}, 3*time.Second).Result()
	if err != nil {
		t.Fatalf("RetrieveAllPosts failed: %v", err)
	}

	posts, ok := res.([]*schemas.Post)
	if !ok || len(posts) != 1 {
		t.Fatalf("Invalid response for RetrieveAllPosts: %v", res)
	}

	res, err = system.Root.RequestFuture(postManager, &proto_actor.RemovePost{ContentID: post.ID}, 3*time.Second).Result()
	if err != nil {
		t.Fatalf("RemovePost failed: %v", err)
	}

	deleted, ok := res.(bool)
	if !ok || !deleted {
		t.Fatalf("Invalid response for RemovePost: %v", res)
	}
}

func BenchmarkAddPost(b *testing.B) {
	system := actor.NewActorSystem()
	postManager := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return proto_actor.NewPostManager()
	}))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := system.Root.RequestFuture(postManager, &proto_actor.AddPost{
			ForumID:  "benchmark-forum",
			AuthorID: "benchmark-user",
			Text:     "Benchmark content",
		}, 3*time.Second).Result()
		if err != nil {
			b.Fatalf("AddPost failed: %v", err)
		}
	}
}

func BenchmarkRetrievePost(b *testing.B) {
	system := actor.NewActorSystem()
	postManager := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return proto_actor.NewPostManager()
	}))

	res, err := system.Root.RequestFuture(postManager, &proto_actor.AddPost{
		ForumID:  "benchmark-forum",
		AuthorID: "benchmark-user",
		Text:     "Benchmark content",
	}, 3*time.Second).Result()
	if err != nil {
		b.Fatalf("AddPost failed: %v", err)
	}

	post, ok := res.(*schemas.Post)
	if !ok || post == nil {
		b.Fatalf("Invalid response for AddPost: %v", res)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := system.Root.RequestFuture(postManager, &proto_actor.RetrievePost{ContentID: post.ID}, 3*time.Second).Result()
		if err != nil {
			b.Fatalf("RetrievePost failed: %v", err)
		}
	}
}

func BenchmarkRetrieveAllPosts(b *testing.B) {
	system := actor.NewActorSystem()
	postManager := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return proto_actor.NewPostManager()
	}))

	for i := 0; i < 100; i++ {
		_, err := system.Root.RequestFuture(postManager, &proto_actor.AddPost{
			ForumID:  "benchmark-forum",
			AuthorID: "benchmark-user",
			Text:     "Benchmark content",
		}, 3*time.Second).Result()
		if err != nil {
			b.Fatalf("AddPost failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := system.Root.RequestFuture(postManager, &proto_actor.RetrieveAllPosts{}, 3*time.Second).Result()
		if err != nil {
			b.Fatalf("RetrieveAllPosts failed: %v", err)
		}
	}
}

func BenchmarkRemovePost(b *testing.B) {
	system := actor.NewActorSystem()
	postManager := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return proto_actor.NewPostManager()
	}))

	res, err := system.Root.RequestFuture(postManager, &proto_actor.AddPost{
		ForumID:  "benchmark-forum",
		AuthorID: "benchmark-user",
		Text:     "Benchmark content",
	}, 3*time.Second).Result()
	if err != nil {
		b.Fatalf("AddPost failed: %v", err)
	}

	post, ok := res.(*schemas.Post)
	if !ok || post == nil {
		b.Fatalf("Invalid response for AddPost: %v", res)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := system.Root.RequestFuture(postManager, &proto_actor.RemovePost{ContentID: post.ID}, 3*time.Second).Result()
		if err != nil {
			b.Fatalf("RemovePost failed: %v", err)
		}
	}
}

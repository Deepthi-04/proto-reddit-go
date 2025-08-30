package tests

import (
	proto_actor "reddit-clone/core/proto_actors"
	"reddit-clone/schemas"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

func TestForumManager(t *testing.T) {
	
	system := actor.NewActorSystem()
	forumManager := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return proto_actor.NewForumManager()
	}))

	res, err := system.Root.RequestFuture(forumManager, &proto_actor.AddForum{
		Title: "test-forum",
	}, 3*time.Second).Result()
	if err != nil {
		t.Fatalf("AddForum failed: %v", err)
	}

	forum, ok := res.(*schemas.Subreddit)
	if !ok || forum == nil {
		t.Fatalf("Invalid response for AddForum: %v", res)
	}
	if forum.Name != "test-forum" {
		t.Errorf("Unexpected forum name. Got: %v, Expected: %v", forum.Name, "test-forum")
	}

	res, err = system.Root.RequestFuture(forumManager, &proto_actor.RetrieveForum{
		ForumID: forum.ID,
	}, 3*time.Second).Result()
	if err != nil {
		t.Fatalf("RetrieveForum failed: %v", err)
	}

	retrievedForum, ok := res.(*schemas.Subreddit)
	if !ok || retrievedForum == nil || retrievedForum.ID != forum.ID {
		t.Fatalf("Invalid response for RetrieveForum: %v", res)
	}
	if retrievedForum.Name != forum.Name {
		t.Errorf("Mismatched forum name. Got: %v, Expected: %v", retrievedForum.Name, forum.Name)
	}

	res, err = system.Root.RequestFuture(forumManager, &proto_actor.RemoveForum{
		ForumID: forum.ID,
	}, 3*time.Second).Result()
	if err != nil {
		t.Fatalf("RemoveForum failed: %v", err)
	}

	deleted, ok := res.(bool)
	if !ok || !deleted {
		t.Fatalf("Forum deletion failed. Got: %v, Expected: true", deleted)
	}

}


func BenchmarkForumManager(b *testing.B) {

	system := actor.NewActorSystem()
	forumManager := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return proto_actor.NewForumManager()
	}))


	b.Run("AddForum", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := system.Root.RequestFuture(forumManager, &proto_actor.AddForum{
				Title: "benchmark-forum",
			}, 3*time.Second).Result()
			if err != nil {
				b.Fatalf("AddForum failed: %v", err)
			}
		}
	})


	res, err := system.Root.RequestFuture(forumManager, &proto_actor.AddForum{
		Title: "benchmark-forum",
	}, 3*time.Second).Result()
	if err != nil {
		b.Fatalf("AddForum setup failed: %v", err)
	}
	forum, ok := res.(*schemas.Subreddit)
	if !ok || forum == nil {
		b.Fatalf("Invalid response for AddForum setup: %v", res)
	}

	
	b.Run("RetrieveForum", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := system.Root.RequestFuture(forumManager, &proto_actor.RetrieveForum{
				ForumID: forum.ID,
			}, 3*time.Second).Result()
			if err != nil {
				b.Fatalf("RetrieveForum failed: %v", err)
			}
		}
	})


	b.Run("RemoveForum", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			
			res, _ := system.Root.RequestFuture(forumManager, &proto_actor.AddForum{
				Title: "benchmark-delete",
			}, 3*time.Second).Result()
			tempForum, _ := res.(*schemas.Subreddit)

			_, err := system.Root.RequestFuture(forumManager, &proto_actor.RemoveForum{
				ForumID: tempForum.ID,
			}, 3*time.Second).Result()
			if err != nil {
				b.Fatalf("RemoveForum failed: %v", err)
			}
		}
	})
}

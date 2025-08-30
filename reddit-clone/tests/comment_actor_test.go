package tests

import (
	"reddit-clone/core/proto_actors"
	"reddit-clone/schemas"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

func TestCommentService(t *testing.T) {
	system := actor.NewActorSystem()
	commentService := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return proto_actor.NewCommentService()
	}))

	parentRes, err := system.Root.RequestFuture(commentService, &proto_actor.AddComment{
		ParentID: "",
		AuthorID: "test-user-id",
		Content:  "Parent comment",
	}, 3*time.Second).Result()
	if err != nil {
		t.Fatalf("AddComment (parent) failed: %v", err)
	}

	parentComment, ok := parentRes.(*schemas.Comment)
	if !ok || parentComment == nil {
		t.Fatalf("Invalid response for AddComment (parent): %v", parentRes)
	}

	res, err := system.Root.RequestFuture(commentService, &proto_actor.AddComment{
		ParentID: parentComment.ID,
		AuthorID: "test-user-id",
		Content:  "This is a test comment",
	}, 3*time.Second).Result()
	if err != nil {
		t.Fatalf("AddComment failed: %v", err)
	}

	comment, ok := res.(*schemas.Comment)
	if !ok || comment == nil {
		t.Fatalf("Invalid response for AddComment: %v", res)
	}
	if comment.Content != "This is a test comment" {
		t.Errorf("Unexpected comment content. Got: %v, Expected: %v", comment.Content, "This is a test comment")
	}

	res, err = system.Root.RequestFuture(commentService, &proto_actor.FetchComment{
		CommentID: comment.ID,
	}, 3*time.Second).Result()
	if err != nil {
		t.Fatalf("FetchComment failed: %v", err)
	}

	retrievedComment, ok := res.(*schemas.Comment)
	if !ok || retrievedComment == nil || retrievedComment.ID != comment.ID {
		t.Fatalf("Invalid response for FetchComment: %v", res)
	}
	if retrievedComment.Content != comment.Content {
		t.Errorf("Mismatched content. Got: %v, Expected: %v", retrievedComment.Content, comment.Content)
	}

	res, err = system.Root.RequestFuture(commentService, &proto_actor.RemoveComment{
		CommentID: comment.ID,
	}, 3*time.Second).Result()
	if err != nil {
		t.Fatalf("RemoveComment failed: %v", err)
	}

	deleted, ok := res.(bool)
	if !ok || !deleted {
		t.Fatalf("Comment deletion failed. Got: %v, Expected: true", deleted)
	}


}

func BenchmarkCommentService(b *testing.B) {

	system := actor.NewActorSystem()
	commentService := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return proto_actor.NewCommentService()
	}))

	parentRes, err := system.Root.RequestFuture(commentService, &proto_actor.AddComment{
		ParentID: "",
		AuthorID: "benchmark-user-id",
		Content:  "Benchmark parent comment",
	}, 3*time.Second).Result()
	if err != nil {
		b.Fatalf("AddComment (parent) setup failed: %v", err)
	}

	parentComment, ok := parentRes.(*schemas.Comment)
	if !ok || parentComment == nil {
		b.Fatalf("Invalid response for AddComment (parent) setup: %v", parentRes)
	}

	b.Run("AddComment", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := system.Root.RequestFuture(commentService, &proto_actor.AddComment{
				ParentID: parentComment.ID,
				AuthorID: "benchmark-user-id",
				Content:  "Benchmark child comment",
			}, 3*time.Second).Result()
			if err != nil {
				b.Fatalf("AddComment failed: %v", err)
			}
		}
	})

	childRes, err := system.Root.RequestFuture(commentService, &proto_actor.AddComment{
		ParentID: parentComment.ID,
		AuthorID: "benchmark-user-id",
		Content:  "Benchmark child comment",
	}, 3*time.Second).Result()
	if err != nil {
		b.Fatalf("AddComment setup (child) failed: %v", err)
	}
	comment, ok := childRes.(*schemas.Comment)
	if !ok || comment == nil {
		b.Fatalf("Invalid response for AddComment setup (child): %v", childRes)
	}

	b.Run("FetchComment", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := system.Root.RequestFuture(commentService, &proto_actor.FetchComment{
				CommentID: comment.ID,
			}, 3*time.Second).Result()
			if err != nil {
				b.Fatalf("FetchComment failed: %v", err)
			}
		}
	})

	b.Run("RemoveComment", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := system.Root.RequestFuture(commentService, &proto_actor.RemoveComment{
				CommentID: comment.ID,
			}, 3*time.Second).Result()
			if err != nil {
				b.Fatalf("RemoveComment failed: %v", err)
			}
		}
	})
}

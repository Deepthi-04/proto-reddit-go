package tests

import (
	proto_actor "reddit-clone/core/proto_actors"
	"reddit-clone/schemas"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

func TestMemberManager(t *testing.T) {
	
	system := actor.NewActorSystem()
	memberManager := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return proto_actor.NewMemberManager()
	}))

	
	res, err := system.Root.RequestFuture(memberManager, &proto_actor.RegisterUser{
		DisplayName: "test-user",
	}, 3*time.Second).Result()
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	profile, ok := res.(*schemas.Account)
	if !ok || profile == nil {
		t.Fatalf("Invalid response for RegisterUser: %v", res)
	}
	if profile.Username != "test-user" {
		t.Errorf("Unexpected display name. Got: %v, Expected: %v", profile.Username, "test-user")
	}

	
	res, err = system.Root.RequestFuture(memberManager, &proto_actor.FetchUser{
		ProfileID: profile.ID,
	}, 3*time.Second).Result()
	if err != nil {
		t.Fatalf("FetchUser failed: %v", err)
	}

	retrievedProfile, ok := res.(*schemas.Account)
	if !ok || retrievedProfile == nil || retrievedProfile.ID != profile.ID {
		t.Fatalf("Invalid response for FetchUser: %v", res)
	}
	if retrievedProfile.Username != profile.Username {
		t.Errorf("Mismatched display name. Got: %v, Expected: %v", retrievedProfile.Username, profile.Username)
	}

	
	res, err = system.Root.RequestFuture(memberManager, &proto_actor.RemoveUser{
		ProfileID: profile.ID,
	}, 3*time.Second).Result()
	if err != nil {
		t.Fatalf("RemoveUser failed: %v", err)
	}

	removed, ok := res.(bool)
	if !ok || !removed {
		t.Fatalf("User profile removal failed. Got: %v, Expected: true", removed)
	}
}


func BenchmarkMemberManager(b *testing.B) {
	system := actor.NewActorSystem()
	memberManager := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return proto_actor.NewMemberManager()
	}))


	b.Run("RegisterUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := system.Root.RequestFuture(memberManager, &proto_actor.RegisterUser{
				DisplayName: "benchmark-user",
			}, 3*time.Second).Result()
			if err != nil {
				b.Fatalf("RegisterUser failed: %v", err)
			}
		}
	})


	res, err := system.Root.RequestFuture(memberManager, &proto_actor.RegisterUser{
		DisplayName: "benchmark-user",
	}, 3*time.Second).Result()
	if err != nil {
		b.Fatalf("RegisterUser setup failed: %v", err)
	}
	profile, ok := res.(*schemas.Account)
	if !ok || profile == nil {
		b.Fatalf("Invalid response for RegisterUser setup: %v", res)
	}


	b.Run("FetchUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := system.Root.RequestFuture(memberManager, &proto_actor.FetchUser{
				ProfileID: profile.ID,
			}, 3*time.Second).Result()
			if err != nil {
				b.Fatalf("FetchUser failed: %v", err)
			}
		}
	})


	b.Run("RemoveUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := system.Root.RequestFuture(memberManager, &proto_actor.RemoveUser{
				ProfileID: profile.ID,
			}, 3*time.Second).Result()
			if err != nil {
				b.Fatalf("RemoveUser failed: %v", err)
			}
		}
	})
}

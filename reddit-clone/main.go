package main

import (
	"log"
	"reddit-clone/handlers"
	"reddit-clone/core/proto_actors"
	"github.com/asynkron/protoactor-go/actor"
)

var (
	userActor      *actor.PID
	subredditActor *actor.PID
	postActor      *actor.PID
	commentActor   *actor.PID
	messageActor   *actor.PID
)

func main() {
	system := actor.NewActorSystem()

	
	userActor := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor { return 	proto_actor.NewMemberManager() }))
	if userActor == nil {
		log.Fatalf("Failed to initialize UserActor")
	}

	subredditActor := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor { return proto_actor.NewForumManager() }))
	if subredditActor == nil {
		log.Fatalf("Failed to initialize SubredditActor")
	}

	postActor := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor { return 	proto_actor.NewPostManager() }))
	log.Println("PostActor initialized:", postActor != nil)
	if postActor == nil {
		log.Fatalf("Failed to initialize PostActor")
	}

	commentActor := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor { return 	proto_actor.NewCommentService() }))
	if commentActor == nil {
		log.Fatalf("Failed to initialize CommentActor")
	}

	messageActor := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor { return 	proto_actor.NewMessageManager() }))
	if messageActor == nil {
		log.Fatalf("Failed to initialize MessageActor")
	}

	handlers.SubredditActor = subredditActor
	handlers.PostActor = postActor
	handlers.CommentActor = commentActor
	handlers.MessageActor = messageActor


}

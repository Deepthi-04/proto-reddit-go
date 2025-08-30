package engine

import (
	proto_actor "reddit-clone/core/proto_actors"
	"reddit-clone/schemas"
	"sync"
	"github.com/asynkron/protoactor-go/actor"
)


type Engine struct {
	Users        map[string]*schemas.Account   
	Subreddits   map[string]*schemas.Subreddit 
	Messages     map[string][]schemas.Message   
	Mutex        sync.Mutex                    
	System       *actor.ActorSystem
	PostActorPID *actor.PID
}


func NewEngine() *Engine {
	return &Engine{
		Users:      make(map[string]*schemas.Account),
		Subreddits: make(map[string]*schemas.Subreddit),
		Messages:   make(map[string][]schemas.Message),
		System:     actor.NewActorSystem(),
	}
}

func (e *Engine) findPostByID(postID string) *schemas.Post {
	for _, subreddit := range e.Subreddits {
		for _, post := range subreddit.Posts {
			if post.ID == postID {
				return post
			}
		}
	}
	return nil
}

func (e *Engine) InitPostActor() *actor.PID {
	props := actor.PropsFromProducer(func() actor.Actor { return proto_actor.NewPostManager() })
	e.PostActorPID = e.System.Root.Spawn(props)
	return e.PostActorPID
}

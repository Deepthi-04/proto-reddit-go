package proto_actor

import (
	"errors"
	"reddit-clone/schemas"
	"sync"
	"log"
	"github.com/asynkron/protoactor-go/actor"
)



type ForumManager struct {
	forums map[string]*schemas.Subreddit
	lock   sync.Mutex
}

func NewForumManager() *ForumManager {
	return &ForumManager{
		forums: make(map[string]*schemas.Subreddit),
	}
}


type AddForum struct {
	Title string
}

type RetrieveForum struct {
	ForumID string
}

type RemoveForum struct {
	ForumID string
}

func (fm *ForumManager) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *AddForum:
		fm.lock.Lock()
		defer fm.lock.Unlock()

		
		forum := schemas.NewSubreddit(msg.Title)
		fm.forums[forum.ID] = forum
		ctx.Respond(forum)

	case *RetrieveForum:
		fm.lock.Lock()
		defer fm.lock.Unlock()

	
		forum, exists := fm.forums[msg.ForumID]
		if !exists {
			ctx.Respond(errors.New("forum not found"))
		} else {
			ctx.Respond(forum)
		}

	case *RemoveForum:
		fm.lock.Lock()
		defer fm.lock.Unlock()

	
		_, exists := fm.forums[msg.ForumID]
		if exists {
			delete(fm.forums, msg.ForumID)
			ctx.Respond(true)
		} else {
			ctx.Respond(false)
		}
	}
}



type MemberManager struct {
	profiles map[string]*schemas.Account
	lock     sync.Mutex
}

func NewMemberManager() *MemberManager {
	return &MemberManager{
		profiles: make(map[string]*schemas.Account),
	}
}


type RegisterUser struct {
	DisplayName string
}

type FetchUser struct {
	ProfileID string
}

type RemoveUser struct {
	ProfileID string
}


func (mm *MemberManager) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *RegisterUser:
		mm.lock.Lock()
		defer mm.lock.Unlock()

		profile := schemas.NewAccount(msg.DisplayName)
		mm.profiles[profile.ID] = profile
		ctx.Respond(profile)

	case *FetchUser:
		mm.lock.Lock()
		defer mm.lock.Unlock()

		profile, exists := mm.profiles[msg.ProfileID]
		if !exists {
			ctx.Respond(errors.New("user profile not found"))
		} else {
			ctx.Respond(profile)
		}

	case *RemoveUser:
		mm.lock.Lock()
		defer mm.lock.Unlock()

		delete(mm.profiles, msg.ProfileID)
		ctx.Respond(true)
	}
}






type PostManager struct {
	posts map[string]*schemas.Post
	mutex sync.Mutex
}

func NewPostManager() *PostManager {
	return &PostManager{
		posts: make(map[string]*schemas.Post),
	}
}


type AddPost struct {
	ForumID  string
	AuthorID string
	Text     string
}

type RetrievePost struct {
	ContentID string
}

type RetrieveAllPosts struct{} 

type RemovePost struct {
	ContentID string
}

func (pm *PostManager) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *AddPost:
		log.Printf("Received AddPost message: %+v\n", msg)
		pm.mutex.Lock()
		defer pm.mutex.Unlock()

		// Create a new post
		post := schemas.NewPost(msg.AuthorID, msg.ForumID, msg.Text)
		pm.posts[post.ID] = post
		log.Printf("Post added: %+v\n", post)
		ctx.Respond(post)

	case *RetrievePost:
		pm.mutex.Lock()
		defer pm.mutex.Unlock()

		post, exists := pm.posts[msg.ContentID]
		if !exists {
			ctx.Respond(nil)
		} else {
			ctx.Respond(post)
		}

	case *RetrieveAllPosts:
		pm.mutex.Lock()
		defer pm.mutex.Unlock()

		var allPosts []*schemas.Post
		for _, post := range pm.posts {
			allPosts = append(allPosts, post)
		}
		ctx.Respond(allPosts)

	case *RemovePost:
		pm.mutex.Lock()
		defer pm.mutex.Unlock()

		_, exists := pm.posts[msg.ContentID]
		if exists {
			delete(pm.posts, msg.ContentID)
			ctx.Respond(true)
		} else {
			ctx.Respond(false)
		}
	default:
		log.Printf("Unknown message type received: %+v\n", msg)
		ctx.Respond(errors.New("unknown message"))
	}
}




type MessageManager struct {
	messageStore map[string][]schemas.Message
	lock         sync.Mutex
}

func NewMessageManager() *MessageManager {
	return &MessageManager{
		messageStore: make(map[string][]schemas.Message),
	}
}


type SendMessage struct {
	FromUserID string
	ToUserID   string
	Body       string
}

type FetchMessages struct {
	UserID string
}

type RemoveMessage struct {
	MessageID string
}

func (mm *MessageManager) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *SendMessage:
		mm.lock.Lock()
		defer mm.lock.Unlock()

		message := schemas.NewMessage(msg.FromUserID, msg.ToUserID, msg.Body)

		
		mm.messageStore[msg.FromUserID] = append(mm.messageStore[msg.FromUserID], *message)
		mm.messageStore[msg.ToUserID] = append(mm.messageStore[msg.ToUserID], *message)

		ctx.Respond(message)

	case *FetchMessages:
		mm.lock.Lock()
		defer mm.lock.Unlock()

		
		userMessages, exists := mm.messageStore[msg.UserID]
		if !exists {
			
			ctx.Respond([]schemas.Message{})
		} else {
			ctx.Respond(userMessages)
		}

	case *RemoveMessage:
		mm.lock.Lock()
		defer mm.lock.Unlock()

		
		for user, messages := range mm.messageStore {
			for i, message := range messages {
				if message.ID == msg.MessageID {
			
					mm.messageStore[user] = append(messages[:i], messages[i+1:]...)
					break
				}
			}
		}
		ctx.Respond(true)
	}
}



type CommentService struct {
	comments map[string]*schemas.Comment
	mutex    sync.Mutex
}

func NewCommentService() *CommentService {
	return &CommentService{
		comments: make(map[string]*schemas.Comment),
	}
}

type AddComment struct {
	ParentID string
	AuthorID string
	Content  string
}

type FetchComment struct {
	CommentID string
}

type RemoveComment struct {
	CommentID string
}

func (cs *CommentService) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {

	case *AddComment:
		cs.handleAddComment(ctx, msg)

	case *FetchComment:
		cs.handleFetchComment(ctx, msg)

	case *RemoveComment:
		cs.handleRemoveComment(ctx, msg)
	}
}


func (cs *CommentService) handleAddComment(ctx actor.Context, msg *AddComment) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	comment := schemas.NewComment(msg.AuthorID, msg.Content)

	if msg.ParentID != "" {
		parent, exists := cs.comments[msg.ParentID]
		if !exists {
			ctx.Respond(errors.New("parent comment not found"))
			return
		}
		parent.AddReply(comment)
		cs.comments[comment.ID] = comment
	} else {
		cs.comments[comment.ID] = comment
	}

	ctx.Respond(comment)
}


func (cs *CommentService) handleFetchComment(ctx actor.Context, msg *FetchComment) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	comment, exists := cs.comments[msg.CommentID]
	if !exists {
		ctx.Respond(errors.New("comment not found"))
		return
	}
	ctx.Respond(comment)
}


func (cs *CommentService) handleRemoveComment(ctx actor.Context, msg *RemoveComment) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	_, exists := cs.comments[msg.CommentID]
	if !exists {
		ctx.Respond(false)
		return
	}

	for _, c := range cs.comments {
		for i, reply := range c.Replies {
			if reply.ID == msg.CommentID {
				c.Replies = append(c.Replies[:i], c.Replies[i+1:]...)
				break
			}
		}
	}

	delete(cs.comments, msg.CommentID)
	ctx.Respond(true)
}

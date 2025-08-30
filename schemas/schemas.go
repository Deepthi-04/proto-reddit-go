package schemas

import (
	"fmt"
	"math/rand"
	"time"
)


type Post struct {
	ID          string     `json:"id"`
	Content     string     `json:"content"`
	AuthorID    string     `json:"author_id"`
	SubredditID string     `json:"subreddit_id"`
	Upvotes     int        `json:"upvotes"`
	Downvotes   int        `json:"downvotes"`
	Comments    []*Comment `json:"comments"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}


func NewPost(authorID, subredditID, content string) *Post {
	return &Post{
		ID:          GenerateID("post"),
		Content:     content,
		AuthorID:    authorID,
		SubredditID: subredditID, 
		Upvotes:     0,
		Downvotes:   0,
		Comments:    []*Comment{},
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}


func (p *Post) AddComment(comment *Comment) {
	p.Comments = append(p.Comments, comment)
	p.UpdatedAt = time.Now().UTC()
}


func (p *Post) AddUpvote() {
	p.Upvotes++
	p.UpdatedAt = time.Now().UTC()
}


func (p *Post) AddDownvote() {
	p.Downvotes++
	p.UpdatedAt = time.Now().UTC()
}



type Message struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}


func NewMessage(senderID, receiverID, content string) *Message {
	return &Message{
		ID:         GenerateID("message"),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}





type Subreddit struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Members   map[string]bool   `json:"members"` 
	Posts     []*Post           `json:"posts"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func NewSubreddit(name string) *Subreddit {
	return &Subreddit{
		ID:        GenerateID("subreddit"),
		Name:      name,
		Members:   make(map[string]bool),
		Posts:     []*Post{},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}


func (s *Subreddit) AddMember(userID string) {
	s.Members[userID] = true
	s.UpdatedAt = time.Now().UTC()
}


func (s *Subreddit) RemoveMember(userID string) {
	delete(s.Members, userID)
	s.UpdatedAt = time.Now().UTC()
}

func (s *Subreddit) AddPost(post *Post) {
	s.Posts = append(s.Posts, post)
	s.UpdatedAt = time.Now().UTC()
}


type Comment struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"`
	Replies   []*Comment `json:"replies"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


func NewComment(authorID, content string) *Comment {
	return &Comment{
		ID:        GenerateID("comment"),
		Content:   content,
		AuthorID:  authorID,
		Replies:   []*Comment{},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}


func (c *Comment) AddReply(reply *Comment) {
	c.Replies = append(c.Replies, reply)
	c.UpdatedAt = time.Now().UTC()
}



type Account struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Karma     int       `json:"karma"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


func NewAccount(username string) *Account {
	return &Account{
		ID:        GenerateID("user"),
		Username:  username,
		Karma:     0,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}


func (a *Account) IncrementKarma(value int) {
	a.Karma += value
	a.UpdatedAt = time.Now().UTC()
}

func (a *Account) DecrementKarma(value int) {
	a.Karma -= value
	a.UpdatedAt = time.Now().UTC()
}


func GenerateID(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%s_%d", prefix, rand.Int63())
}

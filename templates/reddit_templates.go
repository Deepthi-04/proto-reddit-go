package templates

import "reddit-clone/schemas"


type PostResponse struct {
	ID          string `json:"id"`
	SubredditID string `json:"subreddit_id"`
	AuthorID    string `json:"user_id"`
	Content     string `json:"content"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}


func NewPostResponse(post *schemas.Post) *PostResponse {
	return &PostResponse{
		ID:          post.ID,
		SubredditID: post.SubredditID,
		AuthorID:    post.AuthorID,
		Content:     post.Content,
		CreatedAt:   post.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   post.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}


func NewPostListResponse(posts []*schemas.Post) []*PostResponse {
	responses := make([]*PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = NewPostResponse(post)
	}
	return responses
}



type MessageResponse struct {
	ID         string `json:"id"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

func NewMessageResponse(message *schemas.Message) *MessageResponse {
	return &MessageResponse{
		ID:         message.ID,
		SenderID:   message.SenderID,
		ReceiverID: message.ReceiverID,
		Content:    message.Content,
	}
}




type CommentResponse struct {
	ID       string `json:"id"`
	Content  string `json:"content"`
	AuthorID string `json:"author_id"`
}

func NewCommentResponse(comment *schemas.Comment) *CommentResponse {
	return &CommentResponse{
		ID:       comment.ID,
		Content:  comment.Content,
		AuthorID: comment.AuthorID,
	}
}


type SubredditResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Members int    `json:"members"`
}

func NewSubredditResponse(subreddit *schemas.Subreddit) *SubredditResponse {
	return &SubredditResponse{
		ID:      subreddit.ID,
		Name:    subreddit.Name,
		Members: len(subreddit.Members),
	}
}



type AccountResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Karma    int    `json:"karma"`
}

func NewAccountResponse(account *schemas.Account) *AccountResponse {
	return &AccountResponse{
		ID:       account.ID,
		Username: account.Username,
		Karma:    account.Karma,
	}
}

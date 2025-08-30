package handlers
import "github.com/asynkron/protoactor-go/actor"

import (
    "net/http"
	"reddit-clone/core/proto_actors"
	"reddit-clone/schemas"
	"reddit-clone/templates"
	"time"
    "log"
	"github.com/gin-gonic/gin"
)


var (
    UserActor      *actor.PID
    SubredditActor *actor.PID
    PostActor      *actor.PID
    CommentActor   *actor.PID
    MessageActor   *actor.PID
	RootContext  *actor.RootContext
)


func SubmitPostHandler(c *gin.Context) {
	if RootContext == nil {
		log.Println("Error: RootContext is missing")
		c.JSON(500, gin.H{"error": "Server error: RootContext is missing"})
		return
	}

	if PostActor == nil {
		log.Println("Error: PostActor is unavailable")
		c.JSON(500, gin.H{"error": "Server error: PostActor is unavailable"})
		return
	}

	var req struct {
		ForumID  string `json:"forum_id"`
		AuthorID string `json:"author_id"`
		Text     string `json:"text"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v\n", err)
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	log.Printf("Submitting post: ForumID=%s, AuthorID=%s, Text=%s\n", req.ForumID, req.AuthorID, req.Text)

	result, err := RootContext.RequestFuture(PostActor, &proto_actor.AddPost{
		ForumID:  req.ForumID,
		AuthorID: req.AuthorID,
		Text:     req.Text,
	}, 5*time.Second).Result()

	if err != nil {
		log.Printf("Error from PostActor: %v\n", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	post, ok := result.(*schemas.Post)
	if !ok {
		log.Printf("Unexpected result type from PostActor: %v\n", result)
		c.JSON(500, gin.H{"error": "Failed to process post submission"})
		return
	}

	log.Printf("Post successfully created: %+v\n", post)
	c.JSON(200, templates.NewPostResponse(post))
}

 
func FetchPostHandler(c *gin.Context) {
	contentID := c.Param("id")

	result, err := RootContext.RequestFuture(PostActor, &proto_actor.RetrievePost{
		ContentID: contentID,
	}, 5*time.Second).Result()

	if err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	post, ok := result.(*schemas.Post)
	if !ok {
		c.JSON(500, gin.H{"error": "Failed to process post retrieval"})
		return
	}

	c.JSON(200, templates.NewPostResponse(post))
}


func RemovePostHandler(c *gin.Context) {
	contentID := c.Param("id")

	result, err := RootContext.RequestFuture(PostActor, &proto_actor.RemovePost{
		ContentID: contentID,
	}, 5*time.Second).Result()

	if err != nil || result.(bool) == false {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Post successfully deleted"})
}


func FetchAllPostsHandler(c *gin.Context) {
	if RootContext == nil {
		c.JSON(500, gin.H{"error": "Server error: RootContext is missing"})
		return
	}

	if PostActor == nil {
		c.JSON(500, gin.H{"error": "Server error: PostActor is unavailable"})
		return
	}

	result, err := RootContext.RequestFuture(PostActor, &proto_actor.RetrieveAllPosts{}, 5*time.Second).Result()

	if err != nil {
		c.JSON(500, gin.H{"error": "Error fetching posts"})
		return
	}

	posts, ok := result.([]*schemas.Post)
	if !ok {
		c.JSON(500, gin.H{"error": "Failed to process posts"})
		return
	}

	c.JSON(200, templates.NewPostListResponse(posts))
}




const ActorRequestTimeout = 5 * time.Second


func AddCommentHandler(c *gin.Context) {
	var req struct {
		ParentID string `json:"parent_id"`
		AuthorID string `json:"author_id"`
		Content  string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	result, err := RootContext.RequestFuture(CommentActor, &proto_actor.AddComment{
		ParentID: req.ParentID,
		AuthorID: req.AuthorID,
		Content:  req.Content,
	}, ActorRequestTimeout).Result()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	comment, ok := result.(*schemas.Comment)
	if !ok {
		c.JSON(500, gin.H{"error": "Failed to process comment"})
		return
	}

	c.JSON(200, templates.NewCommentResponse(comment))
}


func FetchCommentHandler(c *gin.Context) {
	commentID := c.Param("id")

	result, err := RootContext.RequestFuture(CommentActor, &proto_actor.FetchComment{
		CommentID: commentID,
	}, ActorRequestTimeout).Result()

	if err != nil {
		c.JSON(404, gin.H{"error": "Comment not found"})
		return
	}

	comment, ok := result.(*schemas.Comment)
	if !ok {
		c.JSON(500, gin.H{"error": "Failed to process comment"})
		return
	}

	c.JSON(200, templates.NewCommentResponse(comment))
}


func RemoveCommentHandler(c *gin.Context) {
	commentID := c.Param("id")

	result, err := RootContext.RequestFuture(CommentActor, &proto_actor.RemoveComment{
		CommentID: commentID,
	}, ActorRequestTimeout).Result()

	if err != nil || !result.(bool) {
		c.JSON(404, gin.H{"error": "Comment not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Comment deleted successfully"})
}





func SendMessageHandler(c *gin.Context) {
	var request struct {
		FromUserID string `json:"from_user_id"`
		ToUserID   string `json:"to_user_id"`
		Body       string `json:"body"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	result, err := RootContext.RequestFuture(MessageActor, &proto_actor.SendMessage{
		FromUserID: request.FromUserID,
		ToUserID:   request.ToUserID,
		Body:       request.Body,
	}, 5*time.Second).Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send the message"})
		return
	}

	c.JSON(http.StatusOK, result)
}


func FetchMessagesHandler(c *gin.Context) {
	userID := c.Query("user_id")


	result, err := RootContext.RequestFuture(MessageActor, &proto_actor.FetchMessages{
		UserID: userID,
	}, 5*time.Second).Result()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No messages available"})
		return
	}

	c.JSON(http.StatusOK, result)
}


func RemoveMessageHandler(c *gin.Context) {
	messageID := c.Param("id")


	result, err := RootContext.RequestFuture(MessageActor, &proto_actor.RemoveMessage{
		MessageID: messageID,
	}, 5*time.Second).Result()

	if err != nil || result.(bool) == false {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found for deletion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message removed successfully"})
}






func AddForumHandler(c *gin.Context) {
	var request struct {
		Title string `json:"title"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}


	result, err := RootContext.RequestFuture(SubredditActor, &proto_actor.AddForum{
		Title: request.Title,
	}, 5*time.Second).Result()

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to add forum"})
		return
	}

	forum, ok := result.(*schemas.Subreddit)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid response from actor"})
		return
	}

	c.JSON(200, templates.NewSubredditResponse(forum))
}


func GetForumHandler(c *gin.Context) {
	forumID := c.Param("id")

	
	result, err := RootContext.RequestFuture(SubredditActor, &proto_actor.RetrieveForum{
		ForumID: forumID,
	}, 5*time.Second).Result()

	if err != nil {
		c.JSON(404, gin.H{"error": "Forum not found"})
		return
	}

	forum, ok := result.(*schemas.Subreddit)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid response from actor"})
		return
	}

	c.JSON(200, templates.NewSubredditResponse(forum))
}


func DeleteForumHandler(c *gin.Context) {
	forumID := c.Param("id")


	result, err := RootContext.RequestFuture(SubredditActor, &proto_actor.RemoveForum{
		ForumID: forumID,
	}, 5*time.Second).Result()

	if err != nil || result.(bool) == false {
		c.JSON(404, gin.H{"error": "Forum not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Forum deleted successfully"})
}





func RegisterUserHandler(c *gin.Context) {
	var request struct {
		DisplayName string `json:"display_name"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input data"})
		return
	}

	
	result, err := RootContext.RequestFuture(MessageActor, &proto_actor.RegisterUser{
		DisplayName: request.DisplayName,
	}, 5*time.Second).Result()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	profile, ok := result.(*schemas.Account)
	if !ok {
		c.JSON(500, gin.H{"error": "Unexpected response from actor"})
		return
	}

	c.JSON(200, templates.NewAccountResponse(profile))
}


func FetchUserHandler(c *gin.Context) {
	profileID := c.Param("id")

	
	result, err := RootContext.RequestFuture(MessageActor, &proto_actor.FetchUser{
		ProfileID: profileID,
	}, 5*time.Second).Result()

	if err != nil {
		c.JSON(404, gin.H{"error": "User profile not found"})
		return
	}

	profile, ok := result.(*schemas.Account)
	if !ok {
		c.JSON(500, gin.H{"error": "Unexpected response from actor"})
		return
	}

	c.JSON(200, templates.NewAccountResponse(profile))
}


func RemoveUserHandler(c *gin.Context) {
	profileID := c.Param("id")

	
	result, err := RootContext.RequestFuture(MessageActor, &proto_actor.RemoveUser{
		ProfileID: profileID,
	}, 5*time.Second).Result()

	if err != nil || result.(bool) == false {
		c.JSON(404, gin.H{"error": "User profile not found"})
		return
	}

	c.JSON(200, gin.H{"message": "User profile removed successfully"})
}

package engine
import "errors"




func (e *Engine) AddUpvote(postID, userID string) error {
    post := e.findPostByID(postID)
    if post == nil {
        return errors.New("post not found")
    }

    post.AddUpvote()
    author, exists := e.Users[post.AuthorID]
    if exists {
        author.IncrementKarma(1)
    }
    return nil
}


func (e *Engine) AddDownvote(postID, userID string) error {
    post := e.findPostByID(postID)
    if post == nil {
        return errors.New("post not found")
    }

    post.AddDownvote()
    author, exists := e.Users[post.AuthorID]
    if exists {
        author.DecrementKarma(1)
    }
    return nil
}

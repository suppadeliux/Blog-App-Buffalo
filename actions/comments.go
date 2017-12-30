package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/mikaelm1/blog_app/models"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// CommentsCreate default implementation.
func CommentsCreatePost(c buffalo.Context) error {
	comment := &models.Comment{}
	user := c.Value("current_user").(*models.User)
	if err := c.Bind(comment); err != nil {
		return errors.WithStack(err)
	}
	tx := c.Value("tx").(*pop.Connection)
	comment.AuthorID = user.ID
	postID, err := uuid.FromString(c.Param("pid"))
	if err != nil {
		return errors.WithStack(err)
	}
	comment.PostID = postID
	verrs, err := tx.ValidateAndCreate(comment)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Flash().Add("danger", "There was an error adding your comment.")
		return c.Redirect(302, "/posts/detail/%s", c.Param("pid"))
	}
	c.Flash().Add("success", "Comment added successfully.")
	return c.Redirect(302, "/posts/detail/%s", c.Param("pid"))
}

// CommentsEdit default implementation.
func CommentsEditGet(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := c.Value("current_user").(*models.User)
	comment := &models.Comment{}
	if err := tx.Find(comment, c.Param("cid")); err != nil {
		return c.Error(404, err)
	}
	// make sure the comment was made by the logged in user
	if user.ID != comment.AuthorID {
		c.Flash().Add("danger", "You are not authorized to view that page.")
		return c.Redirect(302, "/posts/detail/%s", comment.PostID)
	}
	c.Set("comment", comment)
	return c.Render(200, r.HTML("comments/edit.html"))
}

// CommentsEdit default implementation.
func CommentsEditPost(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	comment := &models.Comment{}
	if err := tx.Find(comment, c.Param("cid")); err != nil {
		return c.Error(404, err)
	}
	if err := c.Bind(comment); err != nil {
		return errors.WithStack(err)
	}
	user := c.Value("current_user").(*models.User)
	// make sure the comment was made by the logged in user
	if user.ID != comment.AuthorID {
		c.Flash().Add("danger", "You are not authorized to view that page.")
		return c.Redirect(302, "/posts/detail/%s", comment.PostID)
	}
	verrs, err := tx.ValidateAndUpdate(comment)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("comment", comment)
		c.Set("errors", verrs.Errors)
		return c.Render(422, r.HTML("comments/edit.html"))
	}
	c.Flash().Add("success", "Comment updated successfully.")
	return c.Redirect(302, "/posts/detail/%s", comment.PostID)
}

// CommentsDelete default implementation.
func CommentsDelete(c buffalo.Context) error {
	return c.Render(200, r.HTML("comments/delete.html"))
}

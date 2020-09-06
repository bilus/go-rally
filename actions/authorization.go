package actions

import (
	"context"
	"fmt"
	"rally/models"

	log "github.com/sirupsen/logrus"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
)

func authorizePostManagement(post *models.Post, c buffalo.Context) error {
	if !canManagePost(post, c) {
		return fmt.Errorf("no permission to manage post %q", post.ID)
	}
	return nil
}

func canManagePost(post interface{}, help context.Context) bool {
	p := toPostPtr(post)
	u, err := CurrentUser(help)
	if err != nil {
		return false
	}
	return p.AuthorID == u.ID
}

func authorizeCommentManagement(comment *models.Comment, c buffalo.Context) error {
	if !canManageComment(comment, c) {
		return fmt.Errorf("no permission to manage comment %q", comment.ID)
	}
	return nil
}

func canManageComment(comment interface{}, help context.Context) bool {
	c := toCommentPtr(comment)
	u, err := CurrentUser(help)
	if err != nil {
		return false
	}
	return c.AuthorID == u.ID
}

func authorizeBoardManagement(board *models.Board, c buffalo.Context) error {
	if !canManageBoard(board, c) {
		return fmt.Errorf("no permission to manage board %q", board.ID)
	}
	return nil
}

func canManageBoard(board interface{}, help context.Context) bool {
	b := toBoardPtr(board)

	tx, ok := help.Value("tx").(*pop.Connection)
	if !ok {
		log.Error("No transaction found")
		return false
	}
	u, err := CurrentUser(help)
	if err != nil {
		log.Errorf("Error retrieving current user: %v", err)
		return false
	}

	isOwner, err := tx.Where("board_id = ? AND user_id = ? AND is_owner", b.ID, u.ID).Exists(&models.BoardMember{})
	if err != nil {
		log.Errorf("Error checking board ownership: %v", err)
		return false
	}

	return isOwner
}

func toCommentPtr(comment interface{}) *models.Comment {
	ptr, ok := comment.(*models.Comment)
	if ok {
		return ptr
	}
	val, ok := comment.(models.Comment)
	if ok {
		return &val
	}
	log.Fatalf("Expecting models.Comment or *models.Comment")
	return nil
}

func toPostPtr(post interface{}) *models.Post {
	ptr, ok := post.(*models.Post)
	if ok {
		return ptr
	}
	val, ok := post.(models.Post)
	if ok {
		return &val
	}
	log.Fatalf("Expecting models.Post or *models.Post")
	return nil
}

func toBoardPtr(board interface{}) *models.Board {
	ptr, ok := board.(*models.Board)
	if ok {
		return ptr
	}
	val, ok := board.(models.Board)
	if ok {
		return &val
	}
	log.Fatalf("Expecting models.Board or *models.Board")
	return nil
}

package services

import (
	"rally/models"
)

type ReactionStore interface {
	AddReactionToPost(user *models.User, post *models.Post, emoji string) error
	RemoveReactionToPost(user *models.User, post *models.Post, emoji string) error
	ListReactionToPost(post *models.Post) ([]models.Reaction, error)
}

type ReactionsService struct {
	store ReactionStore
}

func NewReactionsService(store ReactionStore) ReactionsService {
	return ReactionsService{store}
}

func (s ReactionsService) AddReactionToPost(user *models.User, post *models.Post, emoji string) error {
	return s.store.AddReactionToPost(user, post, emoji)
}

func (s ReactionsService) RemoveReactionToPost(user *models.User, post *models.Post, emoji string) error {
	return s.store.RemoveReactionToPost(user, post, emoji)
}

func (s ReactionsService) ListDetailedReactionsToPost(user *models.User, post *models.Post) ([]models.Reaction, error) {
	return s.store.ListReactionToPost(post)
}

func (s ReactionsService) ListAggregateReactionsToPost(user *models.User, post *models.Post) ([]models.Reaction, error) {
	return s.store.ListReactionToPost(post)
}

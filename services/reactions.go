package services

import (
	"rally/models"
)

type ReactionStore interface {
	AddReactionToPost(user *models.User, post *models.Post, key string) error
	RemoveReactionToPost(user *models.User, post *models.Post, key string) error
	ListReactionToPost(post *models.Post, aggregate bool) ([]models.Reaction, error)
}

type ReactionsService struct {
	store ReactionStore
}

func NewReactionsService(store ReactionStore) ReactionsService {
	return ReactionsService{store}
}

func (s ReactionsService) AddReactionToPost(ReactionStore, user *models.User, post *models.Post, key string) error {
	return s.store.AddReactionToPost(user, post, key)
}

func (s ReactionsService) RemoveReactionToPost(ReactionStore, user *models.User, post *models.Post, key string) error {
	return s.store.RemoveReactionToPost(user, post, key)
}

func (s ReactionsService) ListDetailedReactionsToPost(ReactionStore, user *models.User, post *models.Post) ([]models.Reaction, error) {
	return s.store.ListReactionToPost(post, false)
}

func (s ReactionsService) ListAggregateReactionsToPost(ReactionStore, user *models.User, post *models.Post) ([]models.Reaction, error) {
	return s.store.ListReactionToPost(post, true)
}

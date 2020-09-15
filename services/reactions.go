package services

import (
	"rally/models"
)

type ReactionStore interface {
	AddReactionToPost(user *models.User, post *models.Post, key string) error
	RemoveReactionToPost(user *models.User, post *models.Post, key string) error
	ListReactionToPost(post *models.Post, aggregate bool) ([]models.Reaction, error)
}

func AddReactionToPost(store ReactionStore, user *models.User, post *models.Post, key string) error {
	return store.AddReactionToPost(user, post, key)
}

func RemoveReactionToPost(store ReactionStore, user *models.User, post *models.Post, key string) error {
	return store.RemoveReactionToPost(user, post, key)
}

func ListDetailedReactionsToPost(store ReactionStore, user *models.User, post *models.Post) ([]models.Reaction, error) {
	return store.ListReactionToPost(post, false)
}

func ListAggregateReactionsToPost(store ReactionStore, user *models.User, post *models.Post) ([]models.Reaction, error) {
	return store.ListReactionToPost(post, true)
}

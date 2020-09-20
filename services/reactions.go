package services

import (
	"rally/models"
	"sort"

	"github.com/gofrs/uuid"
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
	reactions, err := s.store.ListReactionToPost(post)
	if err != nil {
		return nil, err
	}
	for i, r := range reactions {
		reactions[i].IsMadeByCurrentUser = isMadeByUser(r, user.ID)
	}

	sort.Slice(reactions, func(i, j int) bool {
		return reactions[i].Emoji < reactions[j].Emoji
	})
	return reactions, nil
}

func (s ReactionsService) ListAggregateReactionsToPost(user *models.User, post *models.Post) ([]models.Reaction, error) {
	return s.ListDetailedReactionsToPost(user, post)
}

func isMadeByUser(reaction models.Reaction, id uuid.UUID) bool {
	for _, u := range reaction.Users {
		if u.String() == id.String() {
			return true
		}
	}
	return false
}

package services

import (
	"rally/models"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
)

type ReactionStore interface {
	AddReactionToPost(user *models.User, post *models.Post, emoji string) error
	RemoveReactionToPost(user *models.User, post *models.Post, emoji string) error
	ListReactionsToPost(post *models.Post) ([]models.Reaction, error)
}

type ReactionsService struct {
	store ReactionStore
	Tx    *pop.Connection
}

func NewReactionsService(store ReactionStore, tx *pop.Connection) ReactionsService {
	return ReactionsService{
		store: store,
		Tx:    tx,
	}
}

func (s ReactionsService) AddReactionToPost(user *models.User, post *models.Post, emoji string) error {
	return s.store.AddReactionToPost(user, post, emoji)
}

func (s ReactionsService) RemoveReactionToPost(user *models.User, post *models.Post, emoji string) error {
	return s.store.RemoveReactionToPost(user, post, emoji)
}

func (s ReactionsService) ListReactionsToPost(user *models.User, post *models.Post) ([]models.Reaction, error) {
	reactions, err := s.store.ListReactionsToPost(post)
	if err != nil {
		return nil, err
	}
	for i, r := range reactions {
		reactions[i].IsMadeByCurrentUser = isMadeByUser(r, user.ID)
	}

	// TODO: Cache user data in Redis or pull all users with a single query.
	for i, r := range reactions {
		if err := s.Tx.Where("id IN (?)", interfaces(r.UserIDs)...).All(&reactions[i].Users); err != nil {
			return nil, err
		}
	}
	return reactions, nil
}

func interfaces(uuids []uuid.UUID) []interface{} {
	xs := make([]interface{}, len(uuids))
	for i, u := range uuids {
		xs[i] = u.String()
	}
	return xs
}

func isMadeByUser(reaction models.Reaction, id uuid.UUID) bool {
	for _, u := range reaction.UserIDs {
		if u.String() == id.String() {
			return true
		}
	}
	return false
}

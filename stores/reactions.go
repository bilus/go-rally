package stores

import (
	"rally/models"

	"github.com/gofrs/uuid"
)

type ReactionStore struct {
	storage Storage
}

func NewReactionStore(storage Storage) ReactionStore {
	return ReactionStore{storage}
}

type ReactionInfo struct {
	Reactions map[string][]string `json:"rs"`
}

func (s ReactionStore) AddReactionToPost(user *models.User, post *models.Post, emoji string) error {
	info := &ReactionInfo{}
	return s.storage.UpdateJSON(
		scopedKey("post-reactions", post.ID),
		info,
		func() error {
			return addReaction(emoji, user, info)
		})
}

func addReaction(emoji string, user *models.User, info *ReactionInfo) error {
	if info.Reactions == nil {
		info.Reactions = make(map[string][]string)
	}

	userIDs := info.Reactions[emoji]
	_, found := findUserReaction(user.ID.String(), userIDs)
	if found {
		return nil
	}
	info.Reactions[emoji] = append(userIDs, user.ID.String())
	return nil
}

func (s ReactionStore) RemoveReactionToPost(user *models.User, post *models.Post, emoji string) error {
	info := &ReactionInfo{}
	return s.storage.UpdateJSON(
		scopedKey("post-reactions", post.ID),
		info,
		func() error {
			return removeReaction(emoji, user, info)
		})
}

func removeReaction(emoji string, user *models.User, info *ReactionInfo) error {
	if info.Reactions == nil {
		info.Reactions = make(map[string][]string)
	}

	userIDs := info.Reactions[emoji]
	i, found := findUserReaction(user.ID.String(), userIDs)
	if !found {
		return nil
	}
	userIDs = remove(i, userIDs)
	if len(userIDs) == 0 {
		delete(info.Reactions, emoji)
	} else {
		info.Reactions[emoji] = userIDs
	}
	return nil
}

func remove(i int, ids []string) []string {
	// Remove the element at index i from a.
	ids[i] = ids[len(ids)-1] // Copy last element to index i.
	ids[len(ids)-1] = ""     // Erase last element (write zero value).
	ids = ids[:len(ids)-1]   // Truncate slice.
	return ids
}

func findUserReaction(id string, ids []string) (int, bool) {
	for i, x := range ids {
		if id == x {
			return i, true
		}
	}

	return 0, false
}

func (s ReactionStore) ListReactionToPost(post *models.Post) ([]models.Reaction, error) {
	info := &ReactionInfo{}
	if _, err := s.storage.GetJSON(scopedKey("post-reactions", post.ID), info); err != nil {
		return nil, err
	}

	reactions := make([]models.Reaction, 0, len(info.Reactions))
	for emoji, userIDs := range info.Reactions {
		uids, err := uids(userIDs)
		if err != nil {
			return nil, err
		}
		reactions = append(reactions, models.Reaction{
			Emoji: emoji,
			Users: uids,
		})
	}
	return reactions, nil
}

func uids(xs []string) ([]uuid.UUID, error) {
	var err error
	uids := make([]uuid.UUID, len(xs))
	for i, x := range xs {
		uids[i], err = uuid.FromString(x)
		if err != nil {
			return nil, err
		}
	}
	return uids, nil
}

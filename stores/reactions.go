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

func (s ReactionStore) AddReactionToPost(user *models.User, post *models.Post, emoji string) error {
	info := &ReactionInfo{}
	return s.storage.UpdateJSON(
		scopedKey("post-reactions", post.ID),
		info,
		func() error {
			return info.addReaction(emoji, user)
		})
}

func (s ReactionStore) RemoveReactionToPost(user *models.User, post *models.Post, emoji string) error {
	info := &ReactionInfo{}
	return s.storage.UpdateJSON(
		scopedKey("post-reactions", post.ID),
		info,
		func() error {
			return info.removeReaction(emoji, user)
		})
}

func (s ReactionStore) ListReactionsToPost(post *models.Post) ([]models.Reaction, error) {
	info := &ReactionInfo{}
	if _, err := s.storage.GetJSON(scopedKey("post-reactions", post.ID), info); err != nil {
		return nil, err
	}
	return info.Reactions, nil
}

type ReactionInfo struct {
	Reactions []models.Reaction `json:"rs"`
}

func (i *ReactionInfo) addReaction(emoji string, user *models.User) error {
	reaction, _, found := i.findUserReaction(emoji, user)
	if found {
		return nil
	}
	if reaction != nil {
		reaction.UserIDs = append(reaction.UserIDs, user.ID)
		return nil
	}
	i.Reactions = append(i.Reactions, models.Reaction{
		Emoji:   emoji,
		UserIDs: []uuid.UUID{user.ID},
	})
	return nil
}

func (info *ReactionInfo) findUserReaction(emoji string, user *models.User) (*models.Reaction, int, bool) {
	for i, r := range info.Reactions {
		if r.Emoji == emoji {
			for j, uuid := range r.UserIDs {
				if uuid == user.ID {
					return &info.Reactions[i], j, true
				}
			}
			return &info.Reactions[i], -1, false
		}
	}
	return nil, -1, false
}

func (info *ReactionInfo) removeReaction(emoji string, user *models.User) error {
	reaction, i, found := info.findUserReaction(emoji, user)
	if !found {
		return nil
	}
	reaction.UserIDs = removeUUID(i, reaction.UserIDs)
	if len(reaction.UserIDs) == 0 {
		info.Reactions = removeReactionByEmoji(emoji, info.Reactions)
	}
	return nil
}

func removeUUID(i int, ids []uuid.UUID) []uuid.UUID {
	// Remove the element at index i from a.
	ids[i] = ids[len(ids)-1]      // Copy last element to index i.
	ids[len(ids)-1] = uuid.UUID{} // Erase last element (write zero value).
	ids = ids[:len(ids)-1]        // Truncate slice.
	return ids
}

func removeReactionByEmoji(emoji string, reactions []models.Reaction) []models.Reaction {
	i, ok := findReactionByEmoji(emoji, reactions)
	if !ok {
		return reactions
	}
	return removeReaction(i, reactions)
}

func findReactionByEmoji(emoji string, reactions []models.Reaction) (int, bool) {
	for i, r := range reactions {
		if r.Emoji == emoji {
			return i, true
		}
	}
	return -1, false
}

func removeReaction(i int, reactions []models.Reaction) []models.Reaction {
	// Prevent order.
	copy(reactions[i:], reactions[i+1:])            // Shift reactions[i+1:] left one index.
	reactions[len(reactions)-1] = models.Reaction{} // Erase last element (write zero value)..
	reactions = reactions[:len(reactions)-1]        // Truncate slice.
	return reactions
}

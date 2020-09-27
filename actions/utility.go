package actions

import (
	"fmt"
	"net/url"
	"os"
	"rally/models"

	log "github.com/sirupsen/logrus"
)

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

func avatarURL(user models.User, size string, round bool) string {
	//
	// TODO: Resizing doesn't really work for Google default profile pictures
	// so we return an unresized url and <img> has width and height.
	if user.AvatarURL.Valid {
		cloudinaryID, ok := os.LookupEnv("CLOUDINARY_ACCOUNT")
		if !ok {
			log.Warn("CLOUDINARY_ACCOUNT environment variable not set, avatars will be ugly, resized by <img>")
			return user.AvatarURL.String
		}
		w := sizeToPixels(size)
		h := w
		cloudinaryURL := "https://res.cloudinary.com/%v/image/fetch/c_thumb,f_auto,fl_alpha,h_%v,r_30,w_%v/%v"
		return fmt.Sprintf(cloudinaryURL, cloudinaryID, h, w, user.AvatarURL.String)
	}

	return generateAvatarURL(user.Email, size, round)
}

func anonymousAvatarURL(size string, round bool) string {
	return generateAvatarURL(anonymousAvatarSeed, size, round)
}
func sizeToPixels(size string) int {
	if size == "large" {
		return 64
	} else if size == "small" {
		return 32
	} else {
		return 24
	}
}

func generateAvatarURL(seed, size string, round bool) string {
	px := sizeToPixels(size)
	var opts string
	if round {
		opts = "&b=%%23ffffff&m=8&r=50"
	}
	return fmt.Sprintf("https://avatars.dicebear.com/api/bottts/%v.com.svg?colorful=1&w=%v&h=%v&deterministic=1%v", url.QueryEscape(seed), px, px, opts)
}

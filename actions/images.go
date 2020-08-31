package actions

import (
	"fmt"
	"net/http"
	"rally/buffalox"
	"rally/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
)

type UploadLocation struct {
	FilePath string `json:"filePath"`
}

type UploadSuccess struct {
	Data UploadLocation `json:"data"`
}

type UploadError struct {
	Error string `json:"error"`
}

func ImagesCreate(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	postId, err := uuid.FromString(c.Param("post_id"))
	if err != nil {
		return c.Render(400, r.JSON(&UploadError{err.Error()}))
	}

	f, err := c.File("image")
	if err != nil {
		return c.Render(400, r.JSON(&UploadError{err.Error()}))
	}

	image := &models.Attachment{
		PostID:   postId,
		Filename: f.FileHeader.Filename,
	}

	err = image.Save(f)
	if err != nil {
		return c.Render(500, r.JSON(&UploadError{err.Error()}))
	}

	verrs, err := tx.ValidateAndCreate(image)
	if err != nil {
		return c.Render(500, r.JSON(&UploadError{err.Error()}))
	}
	if verrs.HasAny() {
		return c.Render(500, r.JSON(&UploadError{verrs.String()}))
	}

	success := &UploadSuccess{
		Data: UploadLocation{
			FilePath: fmt.Sprintf("/posts/%s/images/%s", postId, image.ID),
		},
	}

	return c.Render(http.StatusOK, r.JSON(success))
}

func ImagesShow(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	postId := c.Param("post_id")
	imageId := c.Param("image_id")

	// Find the attachment.
	image := &models.Attachment{}
	if err := tx.Find(image, imageId); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Check if it belongs to the post.
	if image.PostID.String() != postId {
		return c.Error(http.StatusNotFound, fmt.Errorf("image not found"))
	}

	// TODO: Check if extension is right.

	rc, err := image.Open()
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	defer rc.Close()

	return c.Render(http.StatusOK, buffalox.Inline(c, image.Filename, rc))
}

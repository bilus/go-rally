package actions

import (
	"fmt"
	"net/http"
	"rally/buffalox"
	"rally/models"
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

func (c PostsController) ImagesCreate() error {
	if err := c.RequirePostWithWriteAccess(); err != nil {
		return c.Render(400, r.JSON(&UploadError{err.Error()}))
	}

	f, err := c.File("image")
	if err != nil {
		return c.Render(400, r.JSON(&UploadError{err.Error()}))
	}

	image := &models.Attachment{
		PostID:   c.Post.ID,
		Filename: f.FileHeader.Filename,
	}

	if err := image.Save(f); err != nil {
		return c.Render(500, r.JSON(&UploadError{err.Error()}))
	}

	verrs, err := c.Tx.ValidateAndCreate(image)
	if err != nil {
		return c.Render(500, r.JSON(&UploadError{err.Error()}))
	}
	if verrs.HasAny() {
		return c.Render(500, r.JSON(&UploadError{verrs.String()}))
	}

	success := &UploadSuccess{
		Data: UploadLocation{
			FilePath: fmt.Sprintf("/posts/%s/images/%s", c.Post.ID, image.ID),
		},
	}

	return c.Render(http.StatusOK, r.JSON(success))
}

func (c PostsController) ImagesShow() error {
	if err := c.RequirePost(); err != nil {
		return err
	}

	// Find the attachment.
	imageId := c.Param("image_id")
	image := &models.Attachment{}
	if err := c.Tx.Find(image, imageId); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Check if it belongs to the post.
	if image.PostID != c.Post.ID {
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

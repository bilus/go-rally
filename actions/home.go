package actions

import (
	"fmt"
	"net/http"
	"rally/models"

	"github.com/gofrs/uuid"
)

func (c AuthenticatedController) Home() error {
	boardID, found, err := c.GetLastBoardID()
	if err != nil || !found {
		return c.Redirect(http.StatusSeeOther, "/dashboard/")
	}

	// TODO: Use route helper.
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/boards/%s", boardID.String()))
}

func (c AuthenticatedController) Changelog() error {
	return c.Render(http.StatusOK, r.HTML("/changelog.md"))
}

func (c AuthenticatedController) GetLastBoardID() (uuid.UUID, bool, error) {
	v := c.Session().Get("last_board_id")
	if v == nil {
		return uuid.UUID{}, false, nil
	}
	s, ok := v.(string)
	if !ok {
		return uuid.UUID{}, false, nil
	}
	boardID, err := uuid.FromString(s)
	if err != nil {
		return uuid.UUID{}, false, err
	}

	q := models.DB.Where("id = ?", boardID.String())
	exists, err := q.Exists(&models.Boards{})
	if err != nil {
		return uuid.UUID{}, false, err
	}
	if !exists {
		return uuid.UUID{}, false, fmt.Errorf("board not found")
	}
	return boardID, true, nil
}

func (c AuthenticatedController) SetLastBoardID(boardID uuid.UUID) {
	c.Session().Set("last_board_id", boardID.String())
}

package services

import (
	"rally/models"

	"github.com/gobuffalo/pop/v5"
)

type Dashboard struct {
	StarredBoards []models.Board `json:"starred_boards"`
	ShowWelcome   bool           `json:"show_welcome"`
}

type DashboardService struct {
	db *pop.Connection

	starService StarService
}

func NewDashboardService(db *pop.Connection) DashboardService {
	return DashboardService{
		db:          db,
		starService: NewStarService(db),
	}
}

func (s DashboardService) GetUserDashboard(user models.User) (Dashboard, error) {
	limit := 5
	starredBoards, err := s.starService.ListStarredBoards(&user, &limit)
	if err != nil {
		return Dashboard{}, err
	}
	numBoards, err := s.db.Count(&models.Board{})
	if err != nil {
		return Dashboard{}, err
	}

	return Dashboard{
		StarredBoards: starredBoards,
		ShowWelcome:   numBoards == 0,
	}, nil
}

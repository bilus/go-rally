package services

import (
	"rally/models"

	"github.com/gobuffalo/pop/v5"
	"github.com/johncgriffin/yogofn"
)

type Dashboard struct {
	OwnBoards     []models.Board `json:"own_boards"`
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
	ownBoards, err := s.ownBoards(user)
	if err != nil {
		return Dashboard{}, err
	}

	return Dashboard{
		OwnBoards:     ownBoards,
		StarredBoards: starredBoards,
		ShowWelcome:   numBoards == 0,
	}, nil
}

func (s DashboardService) ownBoards(user models.User) ([]models.Board, error) {
	ownerships := []models.BoardMember{}
	if err := s.db.Where("user_id =  ?", user.ID).All(&ownerships); err != nil {
		return nil, err
	}
	uids := yogofn.Map(func(o models.BoardMember) interface{} {
		return o.BoardID
	}, ownerships).([]interface{})

	boards := []models.Board{}
	if err := s.db.Where("id IN (?)", uids...).All(&boards); err != nil {
		return nil, err
	}
	return boards, nil
}

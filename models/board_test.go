package models_test

import "rally/models"

func (t *ModelSuite) Test_Board_VotingStrategySerialization() {
	t.MustCreateBoardWithVotingStrategy(models.VotingStrategy{BoardMax: 1})

	b := &models.Board{}
	t.NoError(t.DB.First(b))

	t.Equal(1, b.VotingStrategy.BoardMax)
}

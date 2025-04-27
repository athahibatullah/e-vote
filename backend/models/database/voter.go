package models

type Voter struct {
	ID         string `gorm:"primaryKey"`
	VoterName  string
	Email      string
	Password   string
	VoteStatus bool
}

func (Voter) TableName() string {
	return "t_voter_detail"
}

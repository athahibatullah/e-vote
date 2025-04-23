package models

type Candidate struct {
	ID            string `gorm:"primaryKey"`
	CandidateName string
	VisionMission string
}

func (Candidate) TableName() string {
	return "t_candidate_detail"
}

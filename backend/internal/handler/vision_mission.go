package handler

import (
	"e-vote-system/config"
	database "e-vote-system/models/database"
	"encoding/json"
	"net/http"
)

func VisionMissionHandler(w http.ResponseWriter, r *http.Request) {
	var candidates []database.Candidate

	err := config.DB.Find(&candidates).Error
	if err != nil {
		http.Error(w, "Failed to fetch candidates", http.StatusInternalServerError)
		return
	}

	type CandidateResponse struct {
		ID            string `json:"id"`
		CandidateName string `json:"candidate_name"`
		VisionMission string `json:"vision_mission"`
	}

	var response []CandidateResponse
	for _, c := range candidates {
		response = append(response, CandidateResponse{
			ID:            c.ID,
			CandidateName: c.CandidateName,
			VisionMission: c.VisionMission,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"candidates": response,
	})
}

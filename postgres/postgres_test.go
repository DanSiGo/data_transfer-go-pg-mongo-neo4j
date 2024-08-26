package postgres

import (
	"encoding/json"
	"testing"
)

func TestGetPostgresData(t *testing.T) {
	data, err := GetPostgres()
	if err != nil {
		t.Errorf("Error getting data from Postgres: %v", err)
		return
	}

	if len(data) == 0 {
		t.Errorf("Expected non-empty data, got empty slice")
		return
	}

	var testData []struct {
		LearningObjectID string `json:"learning_object_id"`
		ID               string `json:"id"`
		TranscriptID     string `json:"transcript_id"`
		Mimetype         string `json:"mimetype"`
	}

	err = json.Unmarshal(data, &testData)
	if err != nil {
		t.Errorf("Error unmarshaling data: %v", err)
		return
	}

	// Verificar se os dados estão corretos
	for _, row := range testData {
		if row.LearningObjectID == "" {
			t.Errorf("Expected non-zero learning_object_id, got 0")
			return
		}
		if row.ID == "" {
			t.Errorf("Expected non-zero id, got 0")
			return
		}
		if row.TranscriptID == "" {
			t.Errorf("Expected non-zero transcript_id, got 0")
			return
		}
		if row.Mimetype == "" {
			t.Errorf("Expected non-empty mimetype, got empty string")
			return
		}
	}
}

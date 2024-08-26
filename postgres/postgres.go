package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func GetPostgres() ([]byte, error) {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5430/homero?sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to POSTGRES:", err)
		return nil, err
	}
	defer db.Close()

	fmt.Println("Connected to POSTGRES successfully!")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT learning_object_id, id, transcript_id, mimetype FROM media WHERE transcript_id IS NOT NULL LIMIT 10")
	if err != nil {
		fmt.Println("Error selecting from POSTGRES:", err)
		return nil, err
	}
	defer rows.Close()

	var data []struct {
		LearningObjectID string    `json:"learning_object_id"`
		ID               string    `json:"id"`
		TranscriptID     string    `json:"transcript_id"`
		Mimetype         string `json:"mimetype"`
	}

	for rows.Next() {
		var row struct {
			LearningObjectID string    `json:"learning_object_id"`
			ID               string    `json:"id"`
			TranscriptID     string    `json:"transcript_id"`
			Mimetype         string `json:"mimetype"`
		}
		err := rows.Scan(&row.LearningObjectID, &row.ID, &row.TranscriptID, &row.Mimetype)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		data = append(data, row)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		return nil, err
	}

	return jsonData, nil
}

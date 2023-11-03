package gcp

import (
	"encoding/json"
	"os"
)

type ServiceAccount struct {
	ProjectID string `json:"project_id"`
}

func FetchProjectID(filePath string) (ServiceAccount, error) {
	jsonFile, err := os.ReadFile(filePath)
	var serviceAccount ServiceAccount
	if err != nil {
		return serviceAccount, err
	}

	err = json.Unmarshal(jsonFile, &serviceAccount)
	return serviceAccount, err

}

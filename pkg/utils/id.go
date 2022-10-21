package utils

import "github.com/google/uuid"

func GenerateId() (string, error) {
	id, err := uuid.NewUUID()

	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func Parse(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

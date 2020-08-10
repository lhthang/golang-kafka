package uuid

import (
	uuid2 "github.com/google/uuid"
)

func GenerateUUID() string {
	uuid, err := uuid2.NewRandom()
	if err != nil {
		return ""
	}
	return uuid.String()
}

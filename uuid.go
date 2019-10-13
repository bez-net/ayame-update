package main

import (
	"github.com/google/uuid"
)

func getStringUUID() (str string) {
	uuid := uuid.New()
	return uuid.String()
}

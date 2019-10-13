package main

import (
	"github.com/google/uuid"
)

func strUUID() (str string) {
	uuid := uuid.New()
	return uuid.String()
}

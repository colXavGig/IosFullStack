package main

import (
	"ios_full_stack/data"
	"ios_full_stack/dto"
)

func init() {
	data.AddModelsToGormSetUp(
		&dto.User{},
		&dto.Conversation{},
		&dto.Message{},
	)
	data.MightInitDB()
}

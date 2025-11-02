package web

import (
	"encoding/json"
	"fmt"
	"ios_full_stack/dto"
	"net/http"
)

func GetUserFromBody(r *http.Request) (*dto.User, error) {
	var user dto.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("error decoding body: %w", err)
	}
	return &user, nil
}

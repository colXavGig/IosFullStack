package web

import (
	"encoding/json"
	"fmt"
	"ios_full_stack/dto"
	"ios_full_stack/models"
	"net/http"
	"strconv"
)

func HandleGetUsersConversations(w http.ResponseWriter, r *http.Request) {
	Chain(
		RecoveryFromPanic,
		StartGormTransaction,
		RequireAuth,
	).Complete(http.HandlerFunc(getUsersConversation)).ServeHTTP(w, r)
}

func getUsersConversation(w http.ResponseWriter, r *http.Request) {
	var (
		ctx           = r.Context()
		user          *models.AppUser
		conversations []dto.Conversation
		page          uint

		httpErr *HttpError = nil
	)
	defer func() {
		if httpErr != nil {
			_ = JsonResponseWriter(w, httpErr)
		}
	}()

	pagestr := r.URL.Query().Get("page")
	if pagestr == "" {
		page = 1
	} else {
		pageUint, err := strconv.ParseUint(pagestr, 10, 32)
		if err != nil {
			httpErr = &HttpError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Details: fmt.Sprintf("page must be a number"),
			}
			return
		}
		page = uint(pageUint)
	}

	user = GetUserFromContext(ctx)

	conversations, err := models.GetAllConversationOfUser(ctx, user.User, page)
	if err != nil {
		httpErr = &HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Details: fmt.Sprintf("error getting conversation of user"),
		}
		return
	}

	if err := json.NewEncoder(w).Encode(conversations); err != nil {
		httpErr = &HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Details: fmt.Sprintf("error encoding conversation of user"),
		}
		return
	}

}

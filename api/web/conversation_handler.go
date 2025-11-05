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

func HandleGetConversationMessages(w http.ResponseWriter, r *http.Request) {
	Chain(
		RecoveryFromPanic,
		StartGormTransaction,
		RequireAuth,
	).Complete(http.HandlerFunc(getConversationMessages)).ServeHTTP(w, r)
}

func HandlePostMessageToConversation(w http.ResponseWriter, r *http.Request) {
	Chain(
		RecoveryFromPanic,
		StartGormTransaction,
		RequireAuth,
	).Complete(http.HandlerFunc(postMessageToConversation)).ServeHTTP(w, r)
}

func HandleCreateConversation(w http.ResponseWriter, r *http.Request) {
	Chain(
		RecoveryFromPanic,
		StartGormTransaction,
		RequireAuth,
	).Complete(http.HandlerFunc(createConversation)).ServeHTTP(w, r)
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

func getConversationMessages(w http.ResponseWriter, r *http.Request) {
	var (
		ctx            = r.Context()
		conversationID uint64
		conversation   *models.AppConversation
		messages       []dto.Message
		page           uint
		httpErr        *HttpError = nil
	)

	defer func() {
		if httpErr != nil {
			_ = JsonResponseWriter(w, httpErr)
		}
	}()

	conversationIDStr := r.PathValue("id")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 32)
	if err != nil {
		httpErr = &HttpError{
			Code:    http.StatusBadRequest,
			Message: "Invalid conversation ID",
			Details: err.Error(),
		}
		return
	}

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

	conversation, err = models.TryFindConversationById(ctx, uint(conversationID))
	if err != nil {
		httpErr = &HttpError{
			Code:    http.StatusNotFound,
			Message: "Conversation not found",
			Details: err.Error(),
		}
		return
	}

	messages, err = conversation.GetMessagePage(ctx, page)
	if err != nil {
		httpErr = &HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Error getting messages",
			Details: err.Error(),
		}
		return
	}

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		httpErr = &HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Details: fmt.Sprintf("error encoding messages"),
		}
		return
	}
}

func postMessageToConversation(w http.ResponseWriter, r *http.Request) {
	var (
		ctx            = r.Context()
		conversationID uint64
		conversation   *models.AppConversation
		message        dto.Message
		httpErr        *HttpError = nil
	)

	defer func() {
		if httpErr != nil {
			_ = JsonResponseWriter(w, httpErr)
		}
	}()

	conversationIDStr := r.PathValue("id")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 32)
	if err != nil {
		httpErr = &HttpError{
			Code:    http.StatusBadRequest,
			Message: "Invalid conversation ID",
			Details: err.Error(),
		}
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		httpErr = &HttpError{
			Code:    http.StatusBadRequest,
			Message: "Invalid message format",
			Details: err.Error(),
		}
		return
	}

	conversation, err = models.TryFindConversationById(ctx, uint(conversationID))
	if err != nil {
		httpErr = &HttpError{
			Code:    http.StatusNotFound,
			Message: "Conversation not found",
			Details: err.Error(),
		}
		return
	}

	user := GetUserFromContext(ctx)
	message.SenderID = user.ID

	if err := conversation.AppendMessage(ctx, message); err != nil {
		httpErr = &HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Error sending message",
			Details: err.Error(),
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func createConversation(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		body struct {
			UserID uint `json:"user_id"`
		}
		conversation dto.Conversation
		httpErr      *HttpError = nil
	)

	defer func() {
		if httpErr != nil {
			_ = JsonResponseWriter(w, httpErr)
		}
	}()

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpErr = &HttpError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Details: err.Error(),
		}
		return
	}

	user := GetUserFromContext(ctx)
	otherUser, err := models.TryFindUserById(ctx, body.UserID)
	if err != nil {
		httpErr = &HttpError{
			Code:    http.StatusNotFound,
			Message: "User not found",
			Details: "User to start conversation with not found",
		}
		return
	}

	conversation.Members = []dto.User{user.User, otherUser.User}

	if err := models.CreateConversation(ctx, conversation); err != nil {
		httpErr = &HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Error creating conversation",
			Details: err.Error(),
		}
		return
	}

	if err := json.NewEncoder(w).Encode(conversation); err != nil {
		httpErr = &HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Details: fmt.Sprintf("error encoding conversation"),
		}
		return
	}
}

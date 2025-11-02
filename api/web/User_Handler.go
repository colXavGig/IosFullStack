package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"ios_full_stack/dto"
	"ios_full_stack/models"
	"log"
	"net/http"
)

func HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	Chain(
		RecoveryFromPanic,
		StartGormTransaction,
	).Complete(http.HandlerFunc(createUser)).ServeHTTP(w, r)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var (
		userInfo *dto.User
		cookie   *http.Cookie

		err        error
		httpStatus int = http.StatusAccepted
	)
	defer func() {
		if err != nil {
			_ = JsonResponseWriter(w, &HttpError{
				Code:    httpStatus,
				Message: err.Error(),
				Details: err.Error(),
			})
		}
	}()

	userInfo, err = GetUserFromBody(r)
	if err != nil {
		httpStatus = http.StatusBadRequest
		return
	}

	err = models.RegisterNewAppUser(r.Context(), userInfo)
	if errors.Is(err, models.ErrDuplicatedUser) {
		err = errors.New("user already exists")
		httpStatus = http.StatusBadRequest
		return
	}
	if err != nil {
		return
	}

	cookie, err = GenerateToken(*userInfo)
	if err != nil {
		httpStatus = http.StatusInternalServerError
		return
	}

	http.SetCookie(w, cookie)
}

func HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	Chain(
		RecoveryFromPanic,
		StartGormTransaction,
	).Complete(http.HandlerFunc(authentifyUser)).ServeHTTP(w, r)
}

func authentifyUser(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = r.Context()
		loginInfo dto.User
		cookie    *http.Cookie

		httpErr *HttpError = nil
	)
	defer func() {
		if httpErr != nil {
			w.WriteHeader(httpErr.Code)
			_ = json.NewEncoder(w).Encode(httpErr)
		}
	}()

	if err := json.NewDecoder(r.Body).Decode(&loginInfo); err != nil {
		log.Fatalln(err)
	}

	appUser, found := models.TryFindUserByUsername(ctx, loginInfo.Username)
	if !found {
		httpErr = &HttpError{
			Code:    http.StatusUnauthorized,
			Message: "user not found",
			Details: fmt.Sprintf("user %s not found", loginInfo.Username),
		}
		return
	}

	if !appUser.IsCorrectPassword(loginInfo.Password) {
		httpErr = &HttpError{
			Code:    http.StatusUnauthorized,
			Message: "password incorrect",
			Details: fmt.Sprintf("user %s password is incorrect", loginInfo.Username),
		}
		return
	}

	cookie, err := GenerateToken(loginInfo)
	if err != nil {
		log.Panicln(err)
		return
	}

	http.SetCookie(w, cookie)

}

func HandleGetUserInfo(w http.ResponseWriter, r *http.Request) {
	Chain(
		RecoveryFromPanic,
		StartGormTransaction,
	).Complete(http.HandlerFunc(getUserByUsername)).ServeHTTP(w, r)
}

func getUserByUsername(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      context.Context = r.Context()
		username string
		user     *models.AppUser

		err        error
		httpStatus int = http.StatusAccepted
	)
	defer func() {
		if err != nil {
			_ = JsonResponseWriter(w, &HttpError{
				Code:    httpStatus,
				Message: err.Error(),
				Details: fmt.Sprintf("user %s not found", username),
			})
		}
	}()

	username = r.PathValue("username")
	if username == "" {
		httpStatus = http.StatusBadRequest
		err = errors.New("username is empty")
		return
	}

	user, ok := models.TryFindUserByUsername(ctx, username)
	if !ok {
		httpStatus = http.StatusNotFound
		err = models.ErrUserNotFound
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		httpStatus = http.StatusInternalServerError
		return
	}

}

func HandleSetNewPassword(w http.ResponseWriter, r *http.Request) {
	Chain(
		RecoveryFromPanic,
		StartGormTransaction,
	).Complete(http.HandlerFunc(updatePasswd)).ServeHTTP(w, r)
}

func updatePasswd(w http.ResponseWriter, r *http.Request) {
	var (
		userFromCtx *models.AppUser
		newPasswd   string

		err        error
		httpStatus int = http.StatusAccepted
	)
	defer func() {
		if err != nil {
			_ = JsonResponseWriter(w, &HttpError{
				Code:    httpStatus,
				Message: err.Error(),
				Details: fmt.Sprintf("user %s not found", userFromCtx.Username),
			})
		}
	}()

	ctx := r.Context()

	userFromCtx = GetUserFromContext(ctx)

	var jsonBody map[string]string
	if err = json.NewDecoder(r.Body).Decode(&jsonBody); err != nil {
		httpStatus = http.StatusBadRequest
		return
	}

	newPasswd = jsonBody["new_passwd"]

	if err = userFromCtx.UpdatePassword(ctx, newPasswd); err != nil {
		httpStatus = http.StatusInternalServerError
		return
	}
}

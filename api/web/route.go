package web

import "net/http"

func SetRoutesToMuxiplier(mux *http.ServeMux) {
	mux.HandleFunc("POST /login", HandleLoginUser)
	mux.HandleFunc("POST /register", HandleRegisterUser)
	mux.HandleFunc("GET/user/{username}", HandleGetUserInfo)
	mux.HandleFunc("PATCH /user/password", HandleSetNewPassword)

	mux.HandleFunc("GET /conversation", HandleGetUsersConversations)

}

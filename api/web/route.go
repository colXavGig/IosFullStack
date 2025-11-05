package web

import "net/http"

func SetRoutesToMuxiplier(mux *http.ServeMux) {
	mux.HandleFunc("POST /login", HandleLoginUser)
	mux.HandleFunc("POST /register", HandleRegisterUser)
	mux.HandleFunc("GET/user/{username}", HandleGetUserInfo)
	mux.HandleFunc("PATCH /user/password", HandleSetNewPassword)
	mux.HandleFunc("GET /user/search/{query}", HandleSearchUser)

	mux.HandleFunc("POST /conversation", HandleCreateConversation)
	mux.HandleFunc("GET /conversation", HandleGetUsersConversations)
	mux.HandleFunc("GET /conversation/{id}/message", HandleGetConversationMessages)
	mux.HandleFunc("POST /conversation/{id}/message", HandlePostMessageToConversation)

}

package web

import (
	"encoding/json"
	"ios_full_stack/data"
	"net/http"

	"gorm.io/gorm"
)

func StartGormTransaction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			db *gorm.DB

			err        error
			httpStatus int = http.StatusAccepted
		)
		defer func() {
			if err != nil {
				jsonErr, _ := json.Marshal(map[string]any{
					"error": err,
				})
				http.Error(w, string(jsonErr), httpStatus)
			}
		}()

		db, err = data.GetGormDB()
		if err != nil {
			httpStatus = http.StatusInternalServerError
			return
		}
		err = db.Transaction(func(tx *gorm.DB) error {
			ctx := data.GetContextWithTransaction(r.Context(), tx)
			next.ServeHTTP(w, r.WithContext(ctx))
			return nil
		})
		if err != nil {
			httpStatus = http.StatusInternalServerError
			return
		}
	})
}

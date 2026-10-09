package handlers

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/coderianx/lantern/internal/config"
	"github.com/coderianx/lantern/internal/server/middlewares"
	"github.com/coderianx/lantern/internal/server/models"
	"github.com/coderianx/lantern/internal/store"
	sqlite "modernc.org/sqlite"
	lib "modernc.org/sqlite/lib"
)

func CreateNewBucket(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	claims, ok := middlewares.ClaimsFromContext(ctx)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "unauthorized",
		})

		return
	}

	if claims.Role != "admin" {
		w.WriteHeader(http.StatusForbidden)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "You are not admin",
		})

		return
	}

	var req models.CreateBucketReq

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "Invalid request",
		})

		return
	}

	path := filepath.Join(config.GetLanternPath(), req.BucketName)

	err = os.Mkdir(path, 0755)

	if err != nil {

		if errors.Is(err, fs.ErrExist) {
			w.WriteHeader(http.StatusConflict)

			json.NewEncoder(w).Encode(map[string]any{
				"error": "Bucket folder already exists",
			})

			return
		}

		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "Create Bucket folder error",
		})

		return
	}

	_, err = store.DB.Exec(
		`
		INSERT INTO buckets (
			name
		)
		VALUES (?)
		`,
		req.BucketName,
	)

	if err != nil {
		var sqliteErr *sqlite.Error

		if errors.As(err, &sqliteErr) {

			if sqliteErr.Code() == lib.SQLITE_CONSTRAINT_UNIQUE {
				w.WriteHeader(http.StatusConflict)

				json.NewEncoder(w).Encode(map[string]any{
					"error": "Bucket already exists",
				})

				return
			}
		}

		if removeErr := os.Remove(path); removeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{
				"error": "Database error and delete bucket folder error",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"error": "Database error",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]any{
		"message": "Bucket created successfully",
	})
}

func DeleteBucket(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	claims, ok := middlewares.ClaimsFromContext(ctx)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "unauthorized",
		})

		return
	}

	if claims.Role != "admin" {
		w.WriteHeader(http.StatusForbidden)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "you are not admin",
		})

		return
	}

}

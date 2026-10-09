package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/coderianx/lantern/internal/server/models"
	"github.com/coderianx/lantern/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

func Login(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")

	var req models.LoginReq

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "Invalid Request",
		})

		return
	}

	if len(req.Username) == 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "username is required",
		})

		return
	}

	var (
		userID     int64
		passwordDB string
		roleDB     string
	)

	err = store.DB.QueryRow(
		`
		SELECT id, password, role
		FROM users
		WHERE username = ?
		`,
		req.Username,
	).Scan(
		&userID,
		&passwordDB,
		&roleDB,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)

			json.NewEncoder(w).Encode(map[string]any{
				"error": "User not found",
			})

			return
		}

		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "Database error",
		})

		return
	}

	claims := models.Claims{
		UserID:   userID,
		Username: req.Username,

		Role: roleDB,

		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),

			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(30 * time.Minute),
			),
		},
	}

	access := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	secret := os.Getenv("JWT_SECRET")

	if len(secret) == 0 {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "JWT_SECRET error",
		})

		fmt.Println("JWT_SECRET is not set in .env")

		os.Exit(1)
	}

	accessToken, err := access.SignedString(secret)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "create jwt error",
		})

		return
	}

	w.WriteHeader(200)

	json.NewEncoder(w).Encode(map[string]any{
		"message":     "Success",
		"accessToken": accessToken,
	})
}

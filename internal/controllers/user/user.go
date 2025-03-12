package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	storage "github.com/BinayRajbanshi/go-rest-API/database"
	"github.com/BinayRajbanshi/go-rest-API/internal/models"
	"github.com/BinayRajbanshi/go-rest-API/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request received: '/api/v1/users'")
		var User models.User
		// decode the json in to struct
		err := json.NewDecoder(r.Body).Decode(&User)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.BaseError(fmt.Errorf("body is required")))
			return
		} else if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.BaseError(err))
			return
		}
		fmt.Println("Decoded Body is: ", User)

		// after decoding the body, validate it
		validate := validator.New()

		err = validate.Struct(User)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationErrors(err.(validator.ValidationErrors)))
			return
		}

		userId, err := db.CreateUser(User.Username, User.Email, User.Password)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, err)
		}

		successRes := map[string]int64{
			"success": 1,
			"id":      userId,
		}
		// send back json response
		response.WriteJson(w, http.StatusCreated, successRes)
	}
}

func GetAll(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := db.GetUsers()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.BaseError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, users)
	}
}

func Delete(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.BaseError(fmt.Errorf("valid (only number) id is required in path.")))
			return
		}
		deletedId, err := db.DeleteUser(intId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.BaseError(err))
			return
		}

		successMsg := map[string]int64{
			"success":   1,
			"deletedId": deletedId,
		}
		response.WriteJson(w, http.StatusOK, successMsg)
	}
}

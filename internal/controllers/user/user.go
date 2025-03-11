package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/BinayRajbanshi/go-rest-API/internal/models"
	"github.com/BinayRajbanshi/go-rest-API/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
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

		successRes := map[string]int{
			"success": 1,
		}
		// send back json response
		response.WriteJson(w, http.StatusCreated, successRes)
	}
}

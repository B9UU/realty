package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/b9uu/realty/internal/data"
	"github.com/b9uu/realty/internal/validator"
)

func (app *application) AuthToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidPassword(v, input.Password)
	if !v.Valid() {
		app.failedValidationRespone(w, r, v.Errors)
		return
	}

	user, err := app.models.User.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	match, err := user.Password.Match(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}
	if !user.Activated {
		v.AddError("error", "should activate account please check your email")
		app.failedValidationRespone(w, r, v.Errors)

		token, err := app.models.Token.New(user.ID, time.Hour*3*24, data.ScopeActivation)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		app.Background(func() {
			data := map[string]any{
				"activationToken": token.Plaintext,
				"userID":          user.ID,
				"appName":         "Realty",
				"method":          http.MethodPut,
				"URI":             "/users",
			}
			err = app.mailer.Send(user.Email, "welcome.html", data)
			if err != nil {
				app.logger.PrintError(err, nil)
			}
		})
		return
	}
	token, err := app.models.Token.New(user.ID, time.Hour*24, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, data.Envelope{"token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

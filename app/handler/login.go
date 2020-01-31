package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/lacazethomas/goTodo/app/model"
)

func CreateAccount(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	account := &model.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp,err := account.Create(db) //Create account
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

func Authenticate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	account := &model.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp,err := model.Login(account.Email, account.Password,db)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	type login struct {
		Token string
	}

	respondJSON(w, http.StatusOK, login{resp.Token})
}

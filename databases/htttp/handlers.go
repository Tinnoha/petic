package htttp

import (
	databases "databases/htttp/database"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type HttpHandler struct {
	database databases.Database
}

func NewHttpHandler(db databases.Database) *HttpHandler {
	return &HttpHandler{
		database: db,
	}
}

func HTTPError(w http.ResponseWriter, err error, status int) {
	errdto := databases.ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}

	b, err := json.MarshalIndent(errdto, "", "    ")

	if err != nil {
		panic(err)
	}

	http.Error(w, string(b), status)
}

func (h *HttpHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	answer, err := h.database.GetUsers()

	if err != nil {
		HTTPError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(answer); err != nil {
		fmt.Println("Error to write answer from database")
	}
}

func (h *HttpHandler) Insert(w http.ResponseWriter, r *http.Request) {
	var userdto databases.UserDTO
	err := json.NewDecoder(r.Body).Decode(&userdto)

	if err != nil {
		HTTPError(w, err, http.StatusBadGateway)
	}

	err = h.database.Insert(userdto.FIO, userdto.Username, userdto.Email, userdto.Age)

	if err != nil {
		HTTPError(w, err, http.StatusBadGateway)
	}
}

func (h *HttpHandler) AddCash(w http.ResponseWriter, r *http.Request) {
	username, ok := mux.Vars(r)["Username"]

	if !ok {
		HTTPError(w, errors.New("Ошибка нет такого юзера ты пиздабол"), http.StatusBadRequest)
	}

	fmt.Println("\n Usermae:", username)

	var cashdto databases.CashReciverDTO
	err := json.NewDecoder(r.Body).Decode(&cashdto)

	if err != nil {
		HTTPError(w, err, http.StatusBadGateway)
	}

	err = h.database.AddCash(username, cashdto.Count)

	if err != nil {
		HTTPError(w, err, http.StatusBadGateway)
	}

	user := h.database.GetOneUser(username)

	b, err := json.MarshalIndent(user, "", "    ")

	if err != nil {
		HTTPError(w, err, http.StatusBadGateway)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error to write answer from database")
	}
}

func (h *HttpHandler) Buy(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["Username"]

	var buydto databases.BuyingOperationDTO
	err := json.NewDecoder(r.Body).Decode(&buydto)

	if err != nil {
		HTTPError(w, err, http.StatusBadGateway)
	}

	err = h.database.DelCash(username, buydto.Count)

	if err != nil {
		HTTPError(w, err, http.StatusBadGateway)
	}

	user := h.database.GetOneUser(username)

	b, err := json.MarshalIndent(user, "", "    ")

	if err != nil {
		HTTPError(w, err, http.StatusBadGateway)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error to write answer from database")
	}
}

func (h *HttpHandler) TransferCash(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["Username"]

	var transferdto databases.TransferDTO
	err := json.NewDecoder(r.Body).Decode(&transferdto)

	if err != nil {
		HTTPError(w, err, http.StatusBadGateway)
	}

	err = h.database.TransferCash(username, transferdto.UserTo, transferdto.HowMuch)

	if err != nil {
		HTTPError(w, err, http.StatusBadGateway)
	}

	user := h.database.GetOneUser(username)

	b, err := json.MarshalIndent(user, "", "    ")

	if err != nil {
		HTTPError(w, err, http.StatusBadGateway)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error to write answer from database")
	}
}

func (h *HttpHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["Username"]

	err := h.database.DeleteUser(username)

	if err != nil {
		HTTPError(w, err, http.StatusBadGateway)
	}

	w.WriteHeader(http.StatusNoContent)
}

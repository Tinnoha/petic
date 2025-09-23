package htttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"htttp/htps/repositoriy"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	users repositoriy.Polzovately
}

func NewHTTPHandler(pol repositoriy.Polzovately) HTTPHandlers {
	return HTTPHandlers{
		users: pol,
	}
}

func HTTPError(w http.ResponseWriter, err error, status int) {
	errdto := repositoriy.ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}

	b, err := json.MarshalIndent(errdto, "", "    ")

	if err != nil {
		panic(err)
	}

	http.Error(w, string(b), status)
}

/*
-patern: "/users"
-metgod: GET
-info: -

OK:
-Status: 200 OK
-Answer: JSON with all users

Fail:
-Status: 500 InternalServerError
-Answer: JSON with message error and time
*/

func (h *HTTPHandlers) HandlerGetAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandlerGetAllUsers\n")
	b, err := h.users.GetUsers()

	if err != nil {
		HTTPError(w, err, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		log.Fatal(err)
	}
}

/*
-patern: "/users"
-metgod: Post
-info: JSON with user info

OK:
-Status: 201 Created
-Answer: JSON with this user

Fail:
-Status: 400 Forbidden, 500 InternalServerError
-Answer: JSON with message error and time
*/
func (h *HTTPHandlers) HandlerNewUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandlerNewUser\n")

	userdto := repositoriy.UserDTO{}

	err := json.NewDecoder(r.Body).Decode(&userdto)

	if err != nil {
		HTTPError(w, err, http.StatusBadRequest)
		return
	}

	kolya := repositoriy.NewUser(userdto.FIO, userdto.Username, userdto.Email, userdto.Age, 0)

	b, err := h.users.NewUser(kolya)
	fmt.Println("b in handler", string(b))

	if err != nil {
		if errors.Is(err, repositoriy.ThisNameIsExist) {
			HTTPError(w, err, http.StatusForbidden)
			return
		} else {
			HTTPError(w, err, http.StatusInternalServerError)
			return
		}
	}

	fmt.Println("Пользователь создан:\n")

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		log.Fatal(err)
	}
}

/*
-patern: "/users/cash/{username}"
-metgod: Patch
-info: JSON with username and count of cash

OK:
-Status: 200 OK
-Answer: JSON with this user

Fail:
-Status: 400 Bad Request, 500 InternalServerError
-Answer: JSON with message error and time
*/
func (h *HTTPHandlers) HandlerCashReciver(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandlerCashReciver\n")
	username, ok := mux.Vars(r)["username"]

	if !ok {
		fmt.Println("Ты пиздаюол")
	}

	cashDTO := repositoriy.CashReciverDTO{}

	err := json.NewDecoder(r.Body).Decode(&cashDTO)

	if err != nil {
		HTTPError(w, err, http.StatusBadRequest)
		return
	}
	fmt.Println("pered editbalance")
	serega, err := h.users.EditBalance(cashDTO.Count, username, "Cash", "")
	fmt.Println(string(serega))

	if err != nil {
		HTTPError(w, err, http.StatusInternalServerError)
	}

	fmt.Println("Пользователь пополнил баланс:")

	w.WriteHeader(http.StatusAccepted)
	if _, err := w.Write(serega); err != nil {
		fmt.Println(err)
	}
}

/*
-patern: "/users/transfer/{username}"
-metgod: Patch
-info: JSON with username, count of operation and for what

OK:
-Status: 200 OK
-Answer: JSON with this user

Fail:
-Status: 400 Bad Request,403 Forbidden 500 InternalServerError
-Answer: JSON with message error and time
*/
func (h *HTTPHandlers) HandlerTransferOperation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandlerTransferOperation\n")
	username := mux.Vars(r)["username"]

	transferOperationDTO := repositoriy.TransferDTO{}

	err := json.NewDecoder(r.Body).Decode(&transferOperationDTO)

	if err != nil {
		HTTPError(w, err, http.StatusBadRequest)
		return
	}

	serega, err := h.users.EditBalance(transferOperationDTO.HowMuch, username, "Transfer", transferOperationDTO.UserTo)

	if err != nil {
		if errors.Is(err, repositoriy.NotEnouhgMoney) {
			HTTPError(w, err, http.StatusInternalServerError)
		} else if errors.Is(err, repositoriy.ThisNameIsNotExist) {
			HTTPError(w, err, http.StatusBadRequest)
		} else {
			HTTPError(w, err, http.StatusInternalServerError)
		}
	}

	fmt.Println("Пользователь купил что-то ! Юзер:")
	fmt.Println(string(serega))

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(serega); err != nil {
		log.Fatal(err)
	}
}

/*
-patern: "/users/buy/{username}"
-metgod: Patch
-info: JSON with username and count of cash

OK:
-Status: 200 OK
-Answer: JSON with this user

Fail:
-Status: 400 Bad Request,403 Forbidden 500 InternalServerError
-Answer: JSON with message error and time
*/
func (h *HTTPHandlers) HandlerBuynigOperation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandlerBuynigOperation\n")
	username := mux.Vars(r)["username"]

	buyingOperationDTO := repositoriy.BuyingOperationDTO{}

	err := json.NewDecoder(r.Body).Decode(&buyingOperationDTO)

	if err != nil {
		HTTPError(w, err, http.StatusBadRequest)
		return
	}

	serega, err := h.users.EditBalance(buyingOperationDTO.Count, username, "Buy", buyingOperationDTO.ForWhat)
	fmt.Println("serega", string(serega))

	if err != nil {
		if errors.Is(err, repositoriy.NotEnouhgMoney) {
			HTTPError(w, err, http.StatusInternalServerError)
		} else if errors.Is(err, repositoriy.ThisNameIsNotExist) {
			HTTPError(w, err, http.StatusBadRequest)
		} else {
			HTTPError(w, err, http.StatusInternalServerError)
		}
	}

	fmt.Println("Пользователь купил шо то Юзер:")
	fmt.Println(string(serega))

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(serega); err != nil {
		log.Fatal(err)
	}
}

/*
-patern: "/users/{username}"
-metgod: DELETE
-info: -

OK:
-Status: 204 No content
-Answer: -

Fail:
-Status: 400 Bad Request, 500 InternalServerError
-Answer: JSON with message error and time
*/
func (h *HTTPHandlers) HandlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandlerDeleteUser\n")
	username := mux.Vars(r)["username"]

	err := h.users.DeleteUser(username)

	if err != nil {
		if errors.Is(err, repositoriy.ThisNameIsNotExist) {
			HTTPError(w, err, http.StatusBadRequest)
		} else {
			HTTPError(w, err, http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
-patern: "/users/stop"
-metgod: GET
-info: -

OK:
-Status: 204 No content
-Answer: -

Fail:
-Status: 500 InternalServerError
-Answer: JSON with message error and time
*/
func (h *HTTPHandlers) HandlerStop(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandlerStop\n")

	h.users.Stop()
}

package htttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"htttp/repositoriy"
	"log"
	"net/http"
	"time"
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
	users, b := h.users.GetUsers()

	fmt.Println("Пользователи:")
	fmt.Println(users, "\n")

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
		panic(err)
	}

	kolya := repositoriy.NewUser(userdto.FIO, userdto.Username, userdto.Email, userdto.Age)

	err = h.users.NewUser(kolya)

	if err != nil {
		if errors.Is(err, repositoriy.ThisNameIsExist) {
			HTTPError(w, err, http.StatusForbidden)
			return
		} else {
			HTTPError(w, err, http.StatusInternalServerError)
			return
		}
	}

	b, err := json.MarshalIndent(kolya, "", "    ")

	if err != nil {
		panic(err)
	}

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

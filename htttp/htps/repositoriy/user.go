package repositoriy

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

type user struct {
	ID int

	FIO      string
	Username string
	Email    string
	Age      int

	Balance int

	mtx sync.Mutex
}

func NewUser(fio string, username string, email string, age int) user {
	return user{
		FIO:      fio,
		Username: username,
		Email:    email,
		Age:      age,
		Balance:  0,
	}
}

func (u *user) AddBalanceCash(count int) ([]byte, error) {
	client := &http.Client{}
	urlstr := "http://db-service:8081/users/cash/" + u.Username

	cout := CashReciverDTO{Count: count}

	b, err := json.MarshalIndent(cout, "", "    ")

	if err != nil {
		return nil, err
	}

	data := bytes.NewReader(b)

	req, err := http.NewRequest(http.MethodPatch, urlstr, data)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (u *user) DelBalance(count int, ForWhat string) ([]byte, error) {
	urlstr := "http://db-service:8081/users/buy/" + u.Username

	cout := BuyingOperationDTO{Count: count, ForWhat: ForWhat}

	b, err := json.MarshalIndent(cout, "", "    ")

	if err != nil {
		return nil, err
	}

	data := bytes.NewReader(b)

	req, err := http.NewRequest(http.MethodPatch, urlstr, data)

	if err != nil {
		return nil, err
	}

	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (u *user) PerevodBalance(count int, usernameTo string) ([]byte, error) {
	urlstr := "http://db-service:8081/users/transfer/" + u.Username

	cout := PerevodDTO{UserTo: usernameTo, UserFrom: u.Username, HowMuch: count}

	b, err := json.MarshalIndent(cout, "", "    ")

	if err != nil {
		return nil, err
	}

	data := bytes.NewReader(b)

	req, err := http.NewRequest(http.MethodPatch, urlstr, data)

	if err != nil {
		return nil, err
	}

	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

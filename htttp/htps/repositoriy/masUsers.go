package repositoriy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Polzovately struct {
	users map[string]user
}

func NewPolzovately() Polzovately {
	return Polzovately{
		users: make(map[string]user),
	}
}

func (p *Polzovately) NewUser(vasya user) ([]byte, error) {
	// проверка на существование ника
	if _, ok := p.users[vasya.Username]; ok {
		return nil, ThisNameIsExist
	}

	p.users[vasya.Username] = vasya

	// настройка http запроса
	client := &http.Client{}

	// создаем юрл куда будем отправлять запрос на другой микросервис
	urlstr := "http://db-service:8081/users"

	// Переводим данные из стуктуры в байты и обрабатываем ошибку
	b, err := json.MarshalIndent(vasya, "", "    ")

	if err != nil {
		panic(err)
	}

	// переводим в байты в интерфейс ио.Реадер чтоб можно было отправить запроом информацию
	data := bytes.NewReader(b)

	// формирую запрос и обрабатываю ошибку
	req, err := http.NewRequest(http.MethodPost, urlstr, data)

	if err != nil {
		return nil, err
	}

	// ставлю тип отправляемого файла
	req.Header.Set("Content-Type", "application/json")

	// Выполняю запрос и обрабатываю ошибку
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// всегда закрывать тело запроса
	defer req.Body.Close()

	// закрываю тело ответа
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		return body, err
	}
	// считываю тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// отдаю ответ
	return body, nil
}

func (p *Polzovately) GetUsers() ([]byte, error) {

	client := &http.Client{}

	urlstr := "http://db-service:8081/users"

	req, err := http.NewRequest(http.MethodGet, urlstr, nil)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		return body, err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (p *Polzovately) EditBalance(count int, username string, typeOfOperation string, DopInformation string) ([]byte, error) {
	if _, ok := p.users[username]; !ok {
		return nil, errors.New((ThisNameIsNotExist).Error() + username)
	}

	petya := p.users[username]

	switch typeOfOperation {
	case "Cash":
		bUser, err := petya.AddBalanceCash(count)

		if err != nil {
			return bUser, err
		}
	case "Buy":
		if petya.Balance < count {
			return nil, NotEnouhgMoney
		}
		bUser, err := petya.DelBalance(count, DopInformation)
		if err != nil {
			return bUser, err
		}
	case "Transfer":
		if petya.Balance < count {
			return nil, NotEnouhgMoney
		}
		if _, ok := p.users[DopInformation]; !ok {
			return nil, errors.New((ThisNameIsNotExist).Error() + DopInformation)
		}

		bUser, err := petya.PerevodBalance(count, DopInformation)

		if err != nil {
			return bUser, err
		}
	}

	b, err := json.MarshalIndent(petya, "", "    ")

	if err != nil {
		panic(err)
	}

	return b, nil
}

func (p *Polzovately) DeleteUser(username string) error {
	if _, ok := p.users[username]; !ok {
		return errors.New((ThisNameIsNotExist).Error() + username)
	}

	client := &http.Client{}

	urlstr := "db-service:8081/users/" + username

	req, err := http.NewRequest(http.MethodDelete, urlstr, nil)

	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		_, err := io.ReadAll(resp.Body)
		return err
	}

	_, err = io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	return nil
}

func (p *Polzovately) Stop() {
	client := &http.Client{}

	urlstr := "db-service:8081/stop"

	req, _ := http.NewRequest(http.MethodDelete, urlstr, nil)

	resp, _ := client.Do(req)

	if resp.StatusCode != http.StatusNoContent {
		fmt.Println("Неудача закончить")
	}

	req.Body.Close()
}

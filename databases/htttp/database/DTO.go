package databases

import "time"

type UserDTO struct {
	FIO      string `json:"fio"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Balance  int    `json:"balance"`
}

type PerevodDTO struct {
	UserFrom string `json:"From"`
	UserTo   string `json:"To"`
	HowMuch  int    `json:"Cost"`
}

type TransferDTO struct {
	UserTo  string `json:"To"`
	HowMuch int    `json:"Cost"`
}

type BuyingOperationDTO struct {
	ForWhat string `json:"forwhat"`
	Count   int    `json:"count"`
}

type CashReciverDTO struct {
	Count int `json:"count"`
}

type ErrorDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

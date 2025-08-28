package repositoriy

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

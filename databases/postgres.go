package databases

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

type User struct {
	ID       int     `json:"id"`
	FIO      string  `json:"fio"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Age      int     `json:"age"`
	Balance  float64 `json:"balance"`
}

type database struct {
	database *sql.DB
}

func NewDatabase() database {
	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return database{
		database: db,
	}
}

func (db *database) Start() error {
	rezult, err := db.database.Exec(`
	Create table useris(
		ID SERIAL PRIMARY KEY,
		FIO varchar(100) NOT NULL,
		Username varchar(30) NOT NULL UNIQUE,
		Email varchar(50) NOT NULL UNIQUE,
		Age smallint NOT NULL,
		Balance money DEFAULT 0
	)
`)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rezult)

	return nil
}

func (db *database) GetOneUser(Username string) User {
	Row := db.database.QueryRow("SELECT * FROM useri WHERE Username = $1", Username)
	var user User

	err := Row.Scan(&user.FIO, &user.Username, &user.Email, &user.Age)

	if err != nil {
		log.Fatal(err)
	}

	return user
}

func (db *database) Insert(fio string, Username string, Email string, Age int) error {
	rezult, err := db.database.Exec("INSERT INTO useris (fio, username, email, age) VALUES ($1, $2, $3, $4)",
		fio,
		Username,
		Email,
		Age,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rezult)

	return nil
}

func (db *database) AddCash(username string, count int) error {
	rezult, err := db.database.Exec(`
			UPDATE useris
			SET BALANCE = Balance + $1
			WHERE USERNAME = $2
			`,
		count,
		username,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rezult)

	return nil
}

func (db *database) DelCash(username string, count int) error {
	rezult, err := db.database.Exec(`
			UPDATE useris
			SET BALANCE = Balance - $1
			WHERE USERNAME = $2
			`,
		count,
		username,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rezult)

	return nil
}

func (db *database) TransferCash(From string, To string, count int) error {
	rezult, err := db.database.Exec(`
			BEGIN;
			UPDATE useris
			SET BALANCE = Balance - $1
			WHERE USERNAME = $2;
			UPDATE useris
			SET BALANCE = Balance + $1
			WHERE USERNAME = $3;
			COMMIT;
			`,
		From,
		To,
		count,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rezult)

	return nil
}

func (db *database) GetUsers() error {
	rows, err := db.database.Query("SELECT ID, FIO, Username, Email, Age, balance from useri")
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}

	useri := []User{}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.FIO, &user.Username, &user.Email, &user.Age, &user.Balance)

		if err != nil {
			return fmt.Errorf("error: %s", err)
		}

		useri = append(useri, user)

	}

	b, err := json.MarshalIndent(useri, "", "    ")

	if err != nil {
		return fmt.Errorf("error: %s", err)
	}

	PostKafkaRequest("Users", "", b)

	return nil

}

func (db *database) Stop() error {
	rezult, err := db.database.Exec("DROP TABLE useris")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rezult)

	return nil
}

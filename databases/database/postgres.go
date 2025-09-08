package databases

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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
	databasePostgr *sql.DB
	databaseRedis  *redis.Client
}

func NewDatabase() database {
	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	cfg := NewConfigRedis("redis:6379", "1234", "timoha", 0, 5, 10*time.Second, 5*time.Second)

	redis, err := NewClientRedis(context.Background(), *cfg)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return database{
		databasePostgr: db,
		databaseRedis:  redis,
	}
}

func (db *database) Start() error {
	rezult, err := db.databasePostgr.Exec(`
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
	val, err := db.databaseRedis.HGetAll(context.Background(), Username).Result()
	id, _ := strconv.Atoi(val["ID"])
	age, _ := strconv.Atoi(val["Age"])
	balance, _ := strconv.ParseFloat(val["Balance"], 64)
	if err == nil {
		fmt.Println("Взято из redis-хранилища")
		return User{
			ID:       id,
			FIO:      val["FIO"],
			Username: val["Username"],
			Email:    val["Email"],
			Age:      age,
			Balance:  float64(balance),
		}
	}

	Row := db.databasePostgr.QueryRow("SELECT * FROM useri WHERE Username = $1", Username)
	var user User

	err = Row.Scan(&user.FIO, &user.Username, &user.Email, &user.Age)

	if err != nil {
		log.Fatal(err)
	}

	db.databaseRedis.HSet(context.Background(), user.Username, "ID", strconv.Itoa(user.ID), "FIO", user.FIO, "Username", user.Username, "Email", user.Email, "Age", strconv.Itoa(user.Age), "Balance", strconv.FormatFloat(user.Balance, 'f', -1, 64))

	db.databaseRedis.Expire(context.Background(), user.Username, 30*time.Second)

	return user
}

func (db *database) Insert(fio string, Username string, Email string, Age int) error {
	rezult, err := db.databasePostgr.Exec("INSERT INTO useris (fio, username, email, age) VALUES ($1, $2, $3, $4)",
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
	rezult, err := db.databasePostgr.Exec(`
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
	rezult, err := db.databasePostgr.Exec(`
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
	rezult, err := db.databasePostgr.Exec(`
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
	rows, err := db.databasePostgr.Query("SELECT ID, FIO, Username, Email, Age, balance from useri")
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

func (db *database) DeleteUser(username string) error {
	rezult, err := db.databasePostgr.Exec("DELETE FROM useri  WHERE USERNAME = $1", username)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rezult)

	return nil
}

func (db *database) Stop() error {
	rezult, err := db.databasePostgr.Exec("DROP TABLE useris")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rezult)

	db.databasePostgr.Close()

	return nil
}

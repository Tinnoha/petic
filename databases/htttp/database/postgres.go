package databases

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

type Database struct {
	databasePostgr *sql.DB
	databaseRedis  *redis.Client
}

func NewDatabase() Database {
	// Получаем параметры из environment variables
	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		postgresHost = "postgres"
	}

	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	if postgresPassword == "" {
		postgresPassword = "postgres" // fallback
	}

	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = "admin123"
	}

	fmt.Println("Environment variables:")
	fmt.Printf("POSTGRES_HOST: %s\n", os.Getenv("POSTGRES_HOST"))
	fmt.Printf("POSTGRES_PASSWORD: %s\n", os.Getenv("POSTGRES_PASSWORD"))
	fmt.Printf("REDIS_URL: %s\n", os.Getenv("REDIS_URL"))
	fmt.Printf("REDIS_PASSWORD: %s\n", os.Getenv("REDIS_PASSWORD"))

	// Подключение к PostgreSQL с правильным паролем
	connStr := fmt.Sprintf("host=%s port=5432 user=postgres password=%s dbname=postgres sslmode=disable",
		postgresHost, postgresPassword)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open PostgreSQL connection:", err)
	}

	// Подключение к Redis
	cfg := NewConfigRedis(
		redisAddr,
		redisPassword,
		0,
		5,
		10*time.Second,
		5*time.Second,
	)

	redisClient, err := NewClientRedis(context.Background(), *cfg)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// Проверка подключения к PostgreSQL с retry
	var pgErr error
	for i := 0; i < 5; i++ {
		if pgErr = db.Ping(); pgErr == nil {
			break
		}
		log.Printf("PostgreSQL ping attempt %d failed: %v", i+1, pgErr)
		time.Sleep(2 * time.Second)
	}

	if pgErr != nil {
		log.Fatal("Failed to ping PostgreSQL after retries:", pgErr)
	}

	log.Println("Successfully connected to both PostgreSQL and Redis")

	return Database{
		databasePostgr: db,
		databaseRedis:  redisClient,
	}
}

func (db *Database) Start() error {
	db.databasePostgr.Exec("DROP TABLE useris")
	rezult, err := db.databasePostgr.Exec(`
	Create table useris(
		ID SERIAL PRIMARY KEY,
		FIO varchar(100) NOT NULL,
		Username varchar(30) NOT NULL UNIQUE,
		Email varchar(50) NOT NULL UNIQUE,
		Age smallint NOT NULL,
		Balance decimal DEFAULT 0
	)
`)

	if err != nil {
		log.Fatal("WONDERFUUUUL", err)
	}

	fmt.Println(rezult)

	return nil
}

func (db *Database) GetOneUser(Username string) User {
	val, err := db.databaseRedis.HGetAll(context.Background(), Username).Result()
	if err == nil && len(val) > 0 {

		fmt.Println("Взято из redis-хранилища")

		id, _ := strconv.Atoi(val["ID"])
		age, _ := strconv.Atoi(val["Age"])
		balance, _ := strconv.ParseFloat(val["Balance"], 64)

		return User{
			ID:       id,
			FIO:      val["FIO"],
			Username: val["Username"],
			Email:    val["Email"],
			Age:      age,
			Balance:  float64(balance),
		}
	} else {
		Row := db.databasePostgr.QueryRow("SELECT * FROM useris WHERE Username = $1", Username)
		var user User

		err = Row.Scan(&user.ID, &user.FIO, &user.Username, &user.Email, &user.Age, &user.Balance)

		if err != nil {
			fmt.Println("Та самая ошибка", err)
		}

		db.databaseRedis.HSet(context.Background(), user.Username, "ID", strconv.Itoa(user.ID), "FIO", user.FIO, "Username", user.Username, "Email", user.Email, "Age", strconv.Itoa(user.Age), "Balance", strconv.FormatFloat(user.Balance, 'f', -1, 64))

		db.databaseRedis.Expire(context.Background(), user.Username, 30*time.Second)

		return user
	}
}

func (db *Database) Insert(fio string, Username string, Email string, Age int) error {
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

func (db *Database) AddCash(username string, count int) error {
	rezult, err := db.databasePostgr.Exec(`
			UPDATE useris
			SET BALANCE = Balance + $1
			WHERE USERNAME = $2
			`,
		count,
		username,
	)
	fmt.Println("Мы в addCash Postgrees1")
	if err != nil {
		fmt.Println("Та самая ошибка", err)
		return err
	}

	fmt.Println("Мы в addCash Postgrees2")

	fmt.Println("rexult", rezult)

	return nil
}

func (db *Database) DelCash(username string, count int) error {
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

func (db *Database) TransferCash(From string, To string, count int) error {
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

func (db *Database) GetUsers() ([]byte, error) {
	rows, err := db.databasePostgr.Query("SELECT ID, FIO, Username, Email, Age, Balance from useris")
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	useri := []User{}

	for rows.Next() {
		var user User
		fmt.Printf("%#v", user.Balance)
		err := rows.Scan(&user.ID, &user.FIO, &user.Username, &user.Email, &user.Age, &user.Balance)

		if err != nil {
			return nil, fmt.Errorf("errorScAn: %s", err)
		}

		useri = append(useri, user)

	}

	b, err := json.MarshalIndent(useri, "", "    ")

	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	return b, nil

}

func (db *Database) DeleteUser(username string) error {
	rezult, err := db.databasePostgr.Exec("DELETE FROM useris  WHERE USERNAME = $1", username)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rezult)

	return nil
}

func (db *Database) Stop() error {
	rezult, err := db.databasePostgr.Exec("DROP TABLE useris")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rezult)

	db.databasePostgr.Close()

	return nil
}

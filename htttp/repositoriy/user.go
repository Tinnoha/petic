package repositoriy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
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
	ctx, ctxcancel := context.WithCancel(context.Background())

	fmt.Printf("Пользователь %s пытается положить на счет %d денег!\n", u.Username, count)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		Topic:        "AddBalance",
		Balancer:     &kafka.Hash{},
		RequiredAcks: 1,
	})

	defer writer.Close()

	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(string(u.Username)),
		Topic: "AddBalance",
		Value: []byte(string(count)),
	})

	u.mtx.Lock()
	u.Balance += count
	u.mtx.Unlock()

	if err != nil {
		log.Fatal("Error to write Kafka message: ", err)
		return nil, err
	}

	fmt.Printf("Пользователь %s удачно положил деньги", u.Username)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "AddBalanceAnswer",
	})

	defer reader.Close()

	msg, err := reader.ReadMessage(ctx)

	if err != nil {
		log.Fatal("Error to write Kafka message: ", err)
		return nil, err
	}

	ctxcancel()

	return msg.Value, nil
}

func (u *user) DelBalance(count int, ForWhat string) ([]byte, error) {
	ctx, ctxcancel := context.WithCancel(context.Background())

	fmt.Printf("Пользователь %s пытается купить %s за %d денег!\n", u.Username, ForWhat, count)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		Topic:        "DelBalance",
		Balancer:     &kafka.Hash{},
		RequiredAcks: 1,
	})

	defer writer.Close()

	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(string(u.Username)),
		Topic: "DelBalance",
		Value: []byte(string(count)),
	})

	if err != nil {
		log.Fatal("Error to write Kafka message: ", err)
		return nil, err
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "DelBalanceAnswer",
	})

	defer reader.Close()

	msg, err := reader.ReadMessage(ctx)

	if err != nil {
		log.Fatal("Error to write Kafka message: ", err)
		return nil, err
	}

	ctxcancel()

	fmt.Printf("Пользователь %s удачно купил %s деньги", u.Username, ForWhat)

	return msg.Value, nil
}

func (u *user) PerevodBalance(count int, usernameTo string) ([]byte, error) {
	ctx, ctxcancel := context.WithCancel(context.Background())

	fmt.Printf("Пользователь %s пытается перевести пользователю %s %d денег!\n", u.Username, usernameTo, count)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		Topic:        "PerevodBalance",
		Balancer:     &kafka.Hash{},
		RequiredAcks: 1,
	})

	defer writer.Close()

	perevoddto := PerevodDTO{
		UserFrom: u.Username,
		UserTo:   usernameTo,
		HowMuch:  count,
	}

	mes, err := json.MarshalIndent(perevoddto, "", "    ")

	if err != nil {
		panic(err)
	}

	err = writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(string(u.Username)),
		Topic: "PerevodBalance",
		Value: []byte(mes),
	})

	if err != nil {
		log.Fatal("Error to write Kafka message: ", err)
		return nil, err
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "PerevodBalanceAnswer",
	})

	defer reader.Close()

	msg, err := reader.ReadMessage(ctx)

	if err != nil {
		log.Fatal("Error to write Kafka message: ", err)
		return nil, err
	}

	ctxcancel()

	fmt.Printf("Пользователь %s удачно перевел пользователю %s %d деньги", u.Username, usernameTo, count)

	return msg.Value, nil
}

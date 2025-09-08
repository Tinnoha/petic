package repositoriy

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"
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
	if _, ok := p.users[vasya.Username]; ok {
		return nil, ThisNameIsExist
	}

	ctx, ctxcancel := context.WithCancel(context.Background())

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		Topic:        "NewUser",
		Balancer:     &kafka.Hash{},
		RequiredAcks: 1,
	})

	msg, err := json.MarshalIndent(vasya, "", "    ")

	if err != nil {
		panic(err)
	}

	err = writer.WriteMessages(ctx, kafka.Message{
		Topic: "NewUser",
		Key:   []byte(vasya.Username),
		Value: msg,
	})

	if err != nil {
		log.Fatal(err)
	}

	writer.Close()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "NewUserAnswer",
	})

	msg1, err := reader.ReadMessage(ctx)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	ctxcancel()

	return msg1.Value, nil
}

func (p *Polzovately) GetUsers() (map[string]user, []byte) {
	useri := make(map[string]user)

	ctx, ctxcancel := context.WithCancel(context.Background())

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		Topic:        "UsersGet",
		Balancer:     &kafka.Hash{},
		RequiredAcks: 1,
	})

	err := writer.WriteMessages(ctx, kafka.Message{
		Topic: "UsersGet",
		Value: []byte(string("GetUsers")),
	})

	if err != nil {
		writer.Close()
		log.Fatal(err)
	}

	writer.Close()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "Users",
	})

	defer reader.Close()

	msg, err := reader.ReadMessage(ctx)

	if err != nil {
		log.Fatal(err)
	}

	ctxcancel()

	for k, v := range p.users {
		useri[k] = v
	}

	return useri, msg.Value
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

	ctx, ctxcancel := context.WithCancel(context.Background())

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		RequiredAcks: 1,
		Balancer:     &kafka.Hash{},
		Topic:        "DeleteUser",
	})

	defer writer.Close()

	err := writer.WriteMessages(ctx, kafka.Message{
		Topic: "DeleteUser",
		Key:   []byte(username),
		Value: []byte(username),
	})

	if err != nil {
		log.Fatal(err)
	}

	ctxcancel()

	delete(p.users, username)

	return nil
}

func (p *Polzovately) Stop() {
	ctx, ctxcancel := context.WithCancel(context.Background())

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		RequiredAcks: 1,
		Balancer:     &kafka.Hash{},
		Topic:        "Stop",
	})

	defer writer.Close()

	err := writer.WriteMessages(ctx, kafka.Message{
		Topic: "Stop",
	})

	if err != nil {
		log.Fatal(err)
	}

	ctxcancel()
}

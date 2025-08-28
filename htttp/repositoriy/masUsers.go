package repositoriy

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"
)

type polzovately struct {
	users map[string]user
}

func NewPolzovately() polzovately {
	return polzovately{
		users: make(map[string]user),
	}
}

func (p *polzovately) NewUser(vasya user) error {
	if _, ok := p.users[vasya.Username]; ok {
		return ThisNameIsExist
	}

	ctx, ctxcancel := context.WithCancel(context.Background())

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		Topic:        "NewUser",
		Balancer:     &kafka.Hash{},
		RequiredAcks: 1,
	})

	defer writer.Close()

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

	ctxcancel()

	p.users[vasya.Username] = vasya

	return nil
}

func (p *polzovately) GetUsers() (map[string]user, []byte) {
	useri := make(map[string]user)

	ctx, ctxcancel := context.WithCancel(context.Background())

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

func (p *polzovately) EditBalance(count int, username string, typeOfOperation string, DopInformation string) (user, error) {
	if _, ok := p.users[username]; !ok {
		return user{}, errors.New((ThisNameIsNotExist).Error() + username)
	}

	petya := p.users[username]

	switch typeOfOperation {
	case "Cash":
		petya.AddBalanceCash(count)
	case "Buy":
		petya.DelBalance(count, DopInformation)
	case "Transfer":
		if _, ok := p.users[DopInformation]; !ok {
			return user{}, errors.New((ThisNameIsNotExist).Error() + DopInformation)
		}
		petya.PerevodBalance(count, DopInformation)
	}

	return petya, nil
}

func (p *polzovately) DeleteUser(username string) error {

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

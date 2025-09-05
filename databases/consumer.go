package databases

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

func GetKafkaRequest(TopicName string) kafka.Message {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   TopicName,
	})

	ctx, ctxcancel := context.WithCancel(context.Background())

	defer reader.Close()

	msg, err := reader.ReadMessage(ctx)

	if err != nil {
		log.Fatal(err)
	}

	ctxcancel()

	return msg
}

func PostKafkaRequest(TopicName string, Key string, Value []byte) {
	ctx, ctxcancel := context.WithCancel(context.Background())

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhsot:9092"},
		Topic:        TopicName,
		Balancer:     &kafka.Hash{},
		RequiredAcks: 1,
	})

	defer writer.Close()

	err := writer.WriteMessages(ctx, kafka.Message{
		Topic: TopicName,
		Key:   []byte(Key),
		Value: Value,
	})

	if err != nil {
		log.Fatal(err)
	}

	ctxcancel()
}

func GetUsersGo(db database) {
	GetKafkaRequest("UsersGet")

	err := db.GetUsers()

	if err != nil {
		log.Fatal(err)
	}
}

func PostUsersGo(db database) {
	msg := GetKafkaRequest("NewUser")
	var user User

	err := json.Unmarshal(msg.Value, &user)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Insert(user.FIO, user.Username, user.Email, user.Age)

	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(user, "", "    ")

	if err != nil {
		panic(err)
	}

	PostKafkaRequest("NewUserAnswer", "", b)
}

func TransferCashGo(db database) {
	msg := GetKafkaRequest("PerevodBalance")
	transferdto := PerevodDTO{}

	err := json.Unmarshal(msg.Value, &transferdto)

	if err != nil {
		log.Fatal(err)
	}

	err = db.TransferCash(transferdto.UserFrom, transferdto.UserTo, transferdto.HowMuch)

	if err != nil {
		log.Fatal(err)
	}

}

func AddCashGo(db database) {
	msg := GetKafkaRequest("AddBalance")
	cash := CashReciverDTO{}

	err := json.Unmarshal(msg.Value, &cash)

	if err != nil {
		log.Fatal(err)
	}

	err = db.AddCash(string(msg.Key), cash.Count)

	if err != nil {
		log.Fatal(err)
	}

	kopatich := db.GetOneUser(string(msg.Key))

	b, err := json.MarshalIndent(kopatich, "", "    ")

	if err != nil {
		panic(err)
	}

	PostKafkaRequest("AddBalanceAnswer", string(msg.Key), b)
}

func DelBalance(db database) {
	msg := GetKafkaRequest("DelBalance")
	baldto := BuyingOperationDTO{}

	err := json.Unmarshal(msg.Value, &baldto)

	if err != nil {
		log.Fatal(err)
	}

	err = db.DelCash(string(msg.Key), baldto.Count)

	if err != nil {
		log.Fatal(err)
	}

	kopatich := db.GetOneUser(string(msg.Key))

	b, err := json.MarshalIndent(kopatich, "", "    ")

	if err != nil {
		panic(err)
	}
	PostKafkaRequest("DelBalanceAnswer", string(msg.Key), b)
}

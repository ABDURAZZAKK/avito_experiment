package main

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"

	"github.com/ABDURAZZAKK/avito_experiment/config"
	"github.com/ABDURAZZAKK/avito_experiment/internal/repo/pgdb"
	"github.com/ABDURAZZAKK/avito_experiment/pkg/broker"
	"github.com/ABDURAZZAKK/avito_experiment/pkg/postgres"
	log "github.com/sirupsen/logrus"
)

func createCSVFromUsersSegments(pg *postgres.Postgres, msg map[string]interface{}) {
	usersSegmentsRepo := pgdb.NewUsersSegmentsRepo(pg)
	records, err := usersSegmentsRepo.GetStatsPerPeriod(context.Background(),
		int(msg["year"].(float64)),
		int(msg["month"].(float64)),
	)
	if err != nil {
		log.Fatalf("Consumer usersSegmentsRepo.GetStatsPerPeriod: %v", err)
	}

	columns := []string{"Пользователь", "Сегмент", "Операция", "Дата и время"}
	f, err := os.Create(msg["filename"].(string))
	defer f.Close()

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()
	w.Write(columns)
	for _, record := range records {
		row := []string{strconv.Itoa(record.User), record.Segment, string(record.Operation), record.Created_at.Format("2006-01-02 15:04:05")}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	log.Printf("Succses Create File: %s", msg["filename"].(string))
}

func main() {

	cfg, err := config.NewConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	rabbit, err := broker.NewRabbitMQ(cfg.BROKER.URL)
	// rabbit, err := broker.NewRabbitMQ("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err.Error())
	}
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	// pg, err := postgres.New("postgres://postgre:postgre@localhost:5433/postgre", postgres.MaxPoolSize(20))
	if err != nil {
		panic(err.Error())
	}
	messages, err := rabbit.Channel.Consume(
		rabbit.Queue.Name, // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		panic(err.Error())
	}

	var forever chan struct{}

	go func() {
		for message := range messages {
			msg, err := broker.MsgDeserialize(message.Body)
			if err != nil {
				panic(err)
			}
			switch msg["task"] {
			case "createCSVFromUsersSegments":
				createCSVFromUsersSegments(pg, msg)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

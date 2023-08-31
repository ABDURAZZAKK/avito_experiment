package main

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"
	"time"
	_ "time/tzdata"

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
		log.Fatalf("consumer createCSVFromUsersSegments - usersSegmentsRepo.GetStatsPerPeriod: %v", err)
	}

	columns := []string{"Пользователь", "Сегмент", "Операция", "Дата и время"}
	f, err := os.Create(msg["filename"].(string))
	if err != nil {
		log.Fatalf("consumer createCSVFromUsersSegments - os.Create: %v", err)
	}
	defer f.Close()

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

func callAt(callTime string, f func()) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal(err)
	}
	ctime, err := time.ParseInLocation("2006-01-02 15:04:05", callTime, loc)
	if err != nil {
		log.Fatal(err)
	}
	if time.Now().In(loc).Before(ctime) {
		duration := ctime.Sub(time.Now().In(loc))
		log.Printf("Segment delete from User after %v", duration)
		time.Sleep(duration)
		f()
	}
}

func DeleteSegmentFromUserOnTime(pg *postgres.Postgres, msg map[string]interface{}) {
	var users []int
	if msg["users"] != nil {
		for _, i := range msg["users"].([]interface{}) {
			users = append(users, int(i.(float64)))
		}
	} else {
		users = []int{int(msg["user"].(float64))}
	}

	usersSegmentsRepo := pgdb.NewUsersSegmentsRepo(pg)
	m := msg["segments"].([]interface{})
	segments := make([]string, 0, len(m))
	for _, v := range m {
		segments = append(segments, v.(string))
	}
	err := usersSegmentsRepo.DeleteSegmentFromUser(context.Background(),
		users,
		segments,
	)
	if err != nil {
		log.Fatalf("consumer DeleteSegmentFromUserOnTime - usersSegmentsRepo.DeleteSegmentFromUser: %v", err)
	} else {
		log.Printf("Succses delete Segment from User at time: %s", msg["time"].(string))
	}

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
				go createCSVFromUsersSegments(pg, msg)
			case "DeleteSegmentFromUserOnTime":
				go callAt(msg["time"].(string), func() { DeleteSegmentFromUserOnTime(pg, msg) })
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

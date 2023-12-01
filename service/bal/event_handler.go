package bal

import (
	"aci-adr-go-base/model/common"
	"aci-adr-go-base/model/entity"
	"aci-adr-go-base/service/dal"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.mongodb.org/mongo-driver/bson"
	"go.opentelemetry.io/otel/metric"
)

func Connect(meter metric.Meter, js jetstream.JetStream) {
	//create DB services as needed.
	var db dal.Database[entity.ForexData] = &dal.MongoDbService[entity.ForexData]{Collection: "forex_data"}
	histogram, _ := meter.Float64Histogram(
		os.Getenv("STAGE_NAME")+"_duration",
		metric.WithDescription("The duration of task execution."),
		metric.WithUnit("s"),
	)

	apiCounter, _ := meter.Int64Counter(
		os.Getenv("STAGE_NAME")+"_total_processed",
		metric.WithDescription("Number of API calls."),
		metric.WithUnit("{call}"),
	)

	listenStream, listenStreamErr := js.Stream(context.Background(), os.Getenv("STREAM"))

	if listenStreamErr != nil {
		log.Fatal("Unable to connect to listen to specified stream", os.Getenv("STREAM"))
	}

	cons, _ := listenStream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
		Durable:       os.Getenv("CONSUMER"),
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: os.Getenv("LISTEN_SUBJECT"),
		MaxWaiting:    0,
	})

	cc, err := cons.Consume(func(msg jetstream.Msg) {
		go func() {
			start := time.Now()
			handle(msg, js, db)
			duration := time.Since(start)
			histogram.Record(context.Background(), duration.Seconds())
			apiCounter.Add(context.Background(), 1)
		}()
	})
	if err != nil {
		log.Fatal("Error in Consuming Message.")
	}
	defer cc.Stop()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh
	log.Println("Shutting down gracefully.")
}

func handle(msg jetstream.Msg, js jetstream.JetStream, db dal.Database[entity.ForexData]) {
	startedOn := time.Now().UnixMilli()
	var message common.Message
	unmarshalErr := json.Unmarshal(msg.Data(), &message)
	if unmarshalErr != nil {
		return
	}
	result, _ := db.GetOne(bson.D{
		{Key: "tenantId", Value: message.TenantId},
		{Key: "bankId", Value: message.BankId},
		{Key: "baseCurrency", Value: message.BaseCurrency},
		{Key: "targetCurrency", Value: message.TargetCurrency},
		{Key: "tier", Value: message.Tier},
	})
	if result.BuyRate < 0 {
		log.Println("No BuyRate found")
	}
	postedTime := message.StartedOn
	if len(message.Stages) > 0 {
		postedTime = message.Stages[len(message.Stages)-1].CompletedOn
	}
	message.Stages = append(message.Stages, common.ProcessingInfo{
		Stage:          os.Getenv("STAGE_NAME"),
		NetworkTime:    startedOn - postedTime,
		ProcessingTime: time.Now().UnixMilli() - startedOn,
		StartedOn:      startedOn,
		CompletedOn:    time.Now().UnixMilli(),
		Status:         "COMPLETED",
	})

	processedData, _ := json.Marshal(message)
	_, err := js.PublishAsync(os.Getenv("PUBLISH_SUBJECT"), processedData)
	if err != nil {
		return
	}

	ackErr := msg.Ack()
	if ackErr != nil {
		return
	}
}

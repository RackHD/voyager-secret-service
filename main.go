package main

import (
	"flag"
	"log"

	"github.com/RackHD/voyager-secret-service/models"
	"github.com/RackHD/voyager-utilities/random"
)

var (
	uri          = flag.String("uri", "amqp://guest:guest@rabbitmq:5672/", "AMQP URI")
	exchange     = flag.String("exchange", "voyager-secret-service", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType = flag.String("exchange-type", "topic", "Exchange type - direct|fanout|topic|x-custom")
	queue        = flag.String("queue", "voyager-secret-service-queue", "Ephemeral AMQP queue name")
	bindingKey   = flag.String("key", "requests", "AMQP binding key")
	consumerTag  = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
)

func init() {
	flag.Parse()
}

func main() {

	handler := models.NewHandler(*uri)
	if handler.MQ == nil {
		log.Fatalf("Could not connect to RabbitMQ server\n")
	}
	defer handler.MQ.Close()

	secretServiceQueueName := random.RandQueue()
	_, secretServiceMessages, err := handler.MQ.Listen(*exchange, *exchangeType, secretServiceQueueName, *bindingKey, *consumerTag)
	if err != nil {
		log.Fatalf("Error Listening to RabbitMQ: %s\n", err)
	}

	go func() {
		// Listen for messages in the background in infinite loop
		for m := range secretServiceMessages {
			log.Printf(
				"got %dB delivery on exchange %s: [%v] %s",
				len(m.Body),
				m.Exchange,
				m.DeliveryTag,
				m.Body,
			)
			m.Ack(true)
			go handler.ProcessMessage(&m)
		}
	}()

	select {}

}

package models_test

import (
	"os"

	"github.com/RackHD/voyager-secret-service/models"
	"github.com/RackHD/voyager-utilities/random"
	samqp "github.com/streadway/amqp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {

	Describe("AMQP Message Handling", func() {
		var rabbitMQURL string
		var handler *models.Handler
		var testExchange string
		var testExchangeType string
		var testQueueName string
		var testBindingKey string
		var testConsumerTag string
		var testMessage string
		var testReplyTo string
		var testCorID string

		BeforeEach(func() {
			rabbitMQURL = os.Getenv("RABBITMQ_URL")
			handler = models.NewHandler(rabbitMQURL)
			testExchange = "voyager-secret-service"
			testExchangeType = "topic"
			testQueueName = random.RandQueue()
			testBindingKey = "requests"
			testConsumerTag = "testTag"
			testCorID = "TEST-CORRELATION-ID"
			testReplyTo = testBindingKey

			Expect(handler.MQ).ToNot(Equal(nil))

		})
		AfterEach(func() {
			handler.MQ.Close()
		})

		Context("When a message comes in to voyager-secret-service", func() {

			var deliveries <-chan samqp.Delivery
			var err error

			It("INTEGRATION should handle a password request on the 'voyager-secret-service' exchange", func() {
				_, deliveries, err = handler.MQ.Listen(testExchange, testExchangeType, testQueueName, testBindingKey, testConsumerTag)
				Expect(err).ToNot(HaveOccurred())

				testMessage = `generatePassword`
				err = handler.MQ.Send(testExchange, testExchangeType, testBindingKey, testMessage, testCorID, testReplyTo)
				Expect(err).ToNot(HaveOccurred())

				d := <-deliveries
				d.Ack(false)
				Expect(d.CorrelationId).To(Equal(testCorID))

				err = handler.ProcessMessage(&d)
				Expect(err).ToNot(HaveOccurred())

				d = <-deliveries
				d.Ack(false)
				Expect(string(d.Body)).To(Equal(`{"username":"admin","password":"V0yag3r!"}`))

			})

		})

	})

})

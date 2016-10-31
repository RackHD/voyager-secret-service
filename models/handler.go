package models

import (
  "encoding/json"
  "fmt"
  "log"

  "github.com/RackHD/voyager-utilities/amqp"
  "github.com/RackHD/voyager-utilities/models"

  samqp "github.com/streadway/amqp"
)

type Handler struct {
  MQ *amqp.Client
}

// NewHandler creates new amqp handler
func NewHandler(amqpHandler string) *Handler {
  handler := Handler{}
  handler.MQ = amqp.NewClient(amqpHandler)
  if handler.MQ == nil {
    log.Fatalf("Could not connect to RabbitMQ handler: %s\n", amqpHandler)
  }

  return &handler
}

// ProcessMessage processes a message
func (h *Handler) ProcessMessage(m *samqp.Delivery) error {
  switch m.Exchange {

  case "voyager-secret-service":
    return h.processSecretService(m)

  default:
    err := fmt.Errorf("Unknown exchange name: %s\n", m.Exchange)
    log.Printf("Error: %s", err)
    return err

  }
}

// processSecretService processes a message from the voyager-secret-service exchange
func (h *Handler) processSecretService(d *samqp.Delivery) error {

  log.Printf("Received message on Secret-Service Exchange: %s\n", d.Body)
  switch string(d.Body) {
  case "generatePassword":
    credentials := models.Credentials{
      Username: "admin",
      Password: "V0yag3r!",
    }

    credetialMessage, err := json.Marshal(credentials)
    if err != nil {
      log.Fatalf("Error marshalling credentials %s\n", err)
    }
    err = h.MQ.Send(d.Exchange, "topic", d.ReplyTo, string(credetialMessage), d.CorrelationId, d.RoutingKey)
    if err != nil {
      log.Fatalf("Error sending to RabbitMQ: %s\n", err)
    }
  }
  return nil
}

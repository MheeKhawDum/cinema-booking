package queue

import (
    "encoding/json"
    amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
    ch *amqp.Channel
}

func NewPublisher(url string) (*Publisher, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }
    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }
    ch.QueueDeclare("booking.events", true, false, false, false, nil)
    return &Publisher{ch: ch}, nil
}

func (p *Publisher) Publish(routingKey string, data interface{}) error {
    body, _ := json.Marshal(data)
    return p.ch.Publish(
        "", routingKey, false, false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
}
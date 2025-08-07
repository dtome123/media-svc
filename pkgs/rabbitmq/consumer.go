package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func New(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.Qos(1, 0, false) // Prefetch 1 message per worker
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, channel: ch}, nil
}

func (c *Client) Close() {
	c.channel.Close()
	c.conn.Close()
}

func (c *Client) Consume(ctx context.Context, queue string, handler func([]byte) error) error {
	_, err := c.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := c.channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed")
			}
			err := handler(msg.Body)
			if err != nil {
				log.Printf("Job failed: %v", err)
				msg.Nack(false, true)
				continue
			}
			if err := msg.Ack(false); err != nil {
				log.Printf("Failed to ack message: %v", err)
			}
		}
	}
}

func (c *Client) Publish(queue string, data []byte) error {

	_, err := c.channel.QueueDeclare(
		queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	err = c.channel.Publish(
		"",    // exchange ("" means default)
		queue, // routing key (queue name)
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent, // ensure message survives restart
		},
	)

	return err
}

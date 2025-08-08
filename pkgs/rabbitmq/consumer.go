package rabbitmq

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

func NewPublisher(url string) (*Publisher, error) {
	p := &Publisher{url: url}
	if err := p.connect(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Publisher) connect() error {
	conn, err := amqp.Dial(p.url)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}
	p.conn = conn
	p.channel = ch
	return nil
}

func (p *Publisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}

// Publish gửi message lên queue với cơ chế tái kết nối nếu connection/chanel bị đóng
func (p *Publisher) Publish(queue string, data []byte) error {
	// Nếu channel hoặc conn đóng thì reconnect
	if p.conn.IsClosed() || p.channel.IsClosed() {
		if err := p.connect(); err != nil {
			return fmt.Errorf("reconnect failed: %w", err)
		}
	}

	_, err := p.channel.QueueDeclare(
		queue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"",    // default exchange
		queue, // routing key = queue name
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent, // tin nhắn bền vững
		},
	)

	return err
}

// --------------------------

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

func NewConsumer(url string) (*Consumer, error) {
	c := &Consumer{url: url}
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Consumer) connect() error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	err = ch.Qos(1, 0, false) // prefetch 1 message
	if err != nil {
		ch.Close()
		conn.Close()
		return err
	}

	c.conn = conn
	c.channel = ch
	return nil
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

// Consume chạy loop lấy message và xử lý, có thể dùng context để dừng
func (c *Consumer) Consume(ctx context.Context, queue string, handler func([]byte) error) error {
	_, err := c.channel.QueueDeclare(
		queue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare queue failed: %w", err)
	}

	msgs, err := c.channel.Consume(
		queue,
		"",    // consumer tag auto generated
		false, // autoAck false vì muốn ack thủ công
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		return fmt.Errorf("register consumer failed: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				// Channel closed, cần reconnect hoặc thoát
				return fmt.Errorf("channel closed")
			}

			err := handler(msg.Body)
			if err != nil {
				log.Printf("handler error: %v", err)
				// Nack và requeue lại message
				msg.Nack(false, true)
				continue
			}

			// Ack message thành công
			if err := msg.Ack(false); err != nil {
				log.Printf("ack error: %v", err)
			}
		}
	}
}

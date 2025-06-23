package producer

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

// Producer описывает интерфейс продюсера, способного отправлять
// сообщения в брокер.
type Producer interface {
	// SendStringMessage отправляет строку в определенный топик брокера.
	SendStringMessage(event string) error
}

// KafkaProducer имплементирует интерфейс Producer.
type KafkaProducer struct {
	producer sarama.AsyncProducer
	brokers  []string
	topic    string
}

// NewKafkaProducer создает новый экземпляр KafkaProducer.
//
// Параметры:
//   - brokers - слайс с адресами брокеров
//   - topic - название топика
func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		brokers: brokers,
		topic:   topic,
	}
}

// StartProducer запускает работу продюсера.
func (p *KafkaProducer) StartProducer() error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err := sarama.NewAsyncProducer(p.brokers, config)
	if err != nil {
		return fmt.Errorf("failed to start producer: %v", err)
	}

	p.producer = producer

	// горутина для обработки успехов
	go func() {
		for msg := range producer.Successes() {
			log.Printf("message successfully sended: partition=%d, offset=%d\n", msg.Partition, msg.Offset)
		}
	}()

	// горутина для обработки ошибок
	go func() {
		for err := range producer.Errors() {
			log.Printf("failed to send message: %v\n", err.Err)
		}
	}()

	return nil
}

// StopProducer останавливает отправку сообщений продюсером.
func (p *KafkaProducer) StopProducer() {
	if err := p.producer.Close(); err != nil {
		log.Printf("failed to gracefully stop producer: %v", err)
	}
}

// SendStringMessage отправляет строку в определенный топик брокера.
func (p *KafkaProducer) SendStringMessage(event string) error {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(event),
	}

	p.producer.Input() <- msg

	return nil
}

package consumer

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/child6yo/logger-bot/notification-service/internal/handler"
)

func configSarama() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	config.Consumer.Group.Session.Timeout = 6 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 2 * time.Second

	return config
}

// Connection - структура, определяющая соединение с Kafka-брокером.
type Connection struct {
	brokers, topics []string        // список адресов брокеров, список обрабатываемых топиков
	groupID         string          // айди группы консьюмеров
	handler         handler.Handler // обработчик событий
	numPartitions   int             //количество партиций в топике

	ctx    *context.Context
	cancel *context.CancelFunc

	wg sync.WaitGroup
}

// NewConnection создает новый экземлпяр Connection.
//
// Параметры:
//   - brokers - список адресов брокеров
//   - topics - список обрабатываемых топиков
//   - groupID - айди группы консьюмеров
//   - handler - обработчик событий
//   - numPart - количество партиций в топике
func NewConnection(brokers, topics []string, groupID string, numPart int, handler handler.Handler) *Connection {
	return &Connection{
		brokers:       brokers,
		topics:        topics,
		groupID:       groupID,
		handler:       handler,
		numPartitions: numPart,
	}
}

// RunConsumers запускает консьюмеры в количестве, соответсвующем количеству партиций.
func (c *Connection) RunConsumers() error {
	config := configSarama()

	consumerGroup, err := sarama.NewConsumerGroup(c.brokers, c.groupID, config)
	if err != nil {
		return fmt.Errorf("run consumers: failed to start consumer group: %w", err)
	}

	go func() {
		for err := range consumerGroup.Errors() {
			log.Printf("consumer error: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	c.ctx = &ctx
	c.cancel = &cancel

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, c.topics, GroupHandler{c.handler}); err != nil {
				log.Printf("consumers: error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	log.Println("consumer group is running...")
	return nil
}

// StopConsumers останавливает группу консьюмеров отменой контекста.
// Дожидается завершения всех горутин.
func (c *Connection) StopConsumers() {
	cnl := *c.cancel
	cnl()

	c.wg.Wait()
}

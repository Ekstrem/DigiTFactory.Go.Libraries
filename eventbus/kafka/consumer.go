package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/events"
	kafkago "github.com/segmentio/kafka-go"
)

// Consumer — Kafka consumer для доменных событий.
// Управляет подписками и background consumer goroutines.
type Consumer struct {
	mu       sync.RWMutex
	handlers map[string][]events.DomainEventHandler
	options  Options
	logger   *slog.Logger
	cancels  map[string]context.CancelFunc
}

// NewConsumer создаёт Kafka consumer.
func NewConsumer(options Options, logger *slog.Logger) *Consumer {
	return &Consumer{
		handlers: make(map[string][]events.DomainEventHandler),
		options:  options,
		logger:   logger,
		cancels:  make(map[string]context.CancelFunc),
	}
}

// Subscribe подписывает обработчик на события ограниченного контекста.
func (c *Consumer) Subscribe(boundedContext string, handler events.DomainEventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[boundedContext] = append(c.handlers[boundedContext], handler)
	c.logger.Debug("subscribed handler", "boundedContext", boundedContext)
}

// Unsubscribe отписывает обработчик от событий ограниченного контекста.
func (c *Consumer) Unsubscribe(boundedContext string, handler events.DomainEventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	handlers := c.handlers[boundedContext]
	for i, h := range handlers {
		if h == handler {
			c.handlers[boundedContext] = append(handlers[:i], handlers[i+1:]...)
			return
		}
	}
}

// StartConsuming запускает consumer loop для всех зарегистрированных контекстов.
func (c *Consumer) StartConsuming(ctx context.Context) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for bc := range c.handlers {
		topic := fmt.Sprintf("%s.%s", c.options.TopicPrefix, strings.ToLower(bc))
		childCtx, cancel := context.WithCancel(ctx)
		c.cancels[bc] = cancel

		go c.consumeLoop(childCtx, topic, bc)
		c.logger.Info("started kafka consumer", "topic", topic)
	}
}

// Stop останавливает все consumer goroutines.
func (c *Consumer) Stop() {
	for bc, cancel := range c.cancels {
		cancel()
		c.logger.Info("stopped kafka consumer", "boundedContext", bc)
	}
}

func (c *Consumer) consumeLoop(ctx context.Context, topic string, boundedContext string) {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{c.options.BootstrapServers},
		Topic:   topic,
		GroupID: c.options.GroupID,
	})
	defer reader.Close()

	c.logger.Info("kafka consumer subscribed", "topic", topic)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return // Штатное завершение
			}
			c.logger.Error("kafka consume error", "topic", topic, "error", err)
			continue
		}

		envelope, err := Deserialize(msg.Value)
		if err != nil {
			c.logger.Error("deserialize error", "topic", topic, "error", err)
			continue
		}

		c.dispatchToHandlers(ctx, boundedContext, envelope)
	}
}

func (c *Consumer) dispatchToHandlers(ctx context.Context, boundedContext string, envelope *Envelope) {
	c.mu.RLock()
	handlers := c.handlers[boundedContext]
	snapshot := make([]events.DomainEventHandler, len(handlers))
	copy(snapshot, handlers)
	c.mu.RUnlock()

	// Конвертируем Envelope обратно в DomainEvent для обработчиков
	event := &events.DomainEvent{
		AggregateID:        envelope.AggregateID,
		Ver:                envelope.Version,
		BoundedContextName: envelope.BoundedContext,
		ContextName:        envelope.BoundedContext,
		Command: events.NewCommandToAggregate(
			envelope.CorrelationToken,
			envelope.CommandName,
			envelope.SubjectName,
			envelope.CommandVersion,
		),
		ResultReason: envelope.Reason,
	}

	for _, handler := range snapshot {
		if err := handler.HandleEvent(ctx, event); err != nil {
			c.logger.Error("error dispatching kafka event", "error", err)
		}
	}
}

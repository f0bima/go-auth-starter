package event

import (
	"database/sql"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	watermillsql "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/f0bima/go-core/telemetry"
)

// NewAMQPPublisher creates a new RabbitMQ Publisher
func NewAMQPPublisher(amqpURI string, logger watermill.LoggerAdapter) (message.Publisher, error) {
	amqpConfig := amqp.NewDurablePubSubConfig(
		amqpURI,
		amqp.GenerateQueueNameTopicNameWithSuffix("events"),
	)

	publisher, err := amqp.NewPublisher(amqpConfig, logger)
	if err != nil {
		return nil, err
	}

	return telemetry.NewTracingPublisherDecorator(publisher, "amqp-publisher"), nil
}

// NewSQLOutboxPublisher creates a publisher that writes to the outbox table in PostgreSQL
func NewSQLOutboxPublisher(db *sql.DB, logger watermill.LoggerAdapter) (message.Publisher, error) {
	sqlPublisher, err := watermillsql.NewPublisher(
		db,
		watermillsql.PublisherConfig{
			SchemaAdapter:        watermillsql.DefaultPostgreSQLSchema{},
			AutoInitializeSchema: false, // Prod-grade: disabled
		},
		logger,
	)
	if err != nil {
		return nil, err
	}
	return forwarder.NewPublisher(sqlPublisher, forwarder.PublisherConfig{
		ForwarderTopic: "events_to_forward",
	}), nil
}

// NewOutboxForwarder creates a forwarder that reads from the SQL outbox table and publishes to AMQP
func NewOutboxForwarder(
	sqlSubscriber *watermillsql.Subscriber,
	amqpPublisher message.Publisher,
	logger watermill.LoggerAdapter,
) (*forwarder.Forwarder, error) {
	return forwarder.NewForwarder(
		sqlSubscriber,
		amqpPublisher,
		logger,
		forwarder.Config{
			ForwarderTopic: "events_to_forward",
			Middlewares:    []message.HandlerMiddleware{},
		},
	)
}

// NewSQLOutboxSubscriber creates a subscriber that reads from the outbox table
func NewSQLOutboxSubscriber(db *sql.DB, logger watermill.LoggerAdapter) (*watermillsql.Subscriber, error) {
	return watermillsql.NewSubscriber(
		db,
		watermillsql.SubscriberConfig{
			SchemaAdapter:    watermillsql.DefaultPostgreSQLSchema{},
			OffsetsAdapter:   watermillsql.DefaultPostgreSQLOffsetsAdapter{},
			InitializeSchema: false, // Prod-grade: disabled
		},
		logger,
	)
}

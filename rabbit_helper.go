package go_rabbit_helper

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	logging "github.com/lordtor/go-logging"
	"github.com/streadway/amqp"
)

var (
	Log = logging.Log
)

func ConnectionURL(con Rabbit) string {
	address := fmt.Sprintf("amqp://%v:%v@%v:%v/%v", con.Username,
		con.Password,
		con.Server,
		con.Port,
		con.Host,
	)
	return address
}

func handleError(err error, msg string) {
	if err != nil {
		Log.Fatalf("%s: %s", msg, err)
	}

}
func PublisherN(config Rabbit, messages []string) {

	// Create connection
	conn, err := amqp.Dial(ConnectionURL(config))
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()
	// Create channel
	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")
	// Close chanel after work
	defer amqpChannel.Close()
	_, queueDeclareErr := amqpChannel.QueueDeclare(
		config.Data.Queue, // name
		true,              // Passive
		false,             //Durable
		false,             //AutoDelete
		false,             //Exclusive
		nil,               //Arguments
	)
	handleError(queueDeclareErr, fmt.Sprintf(
		"Could not declare `%s` queue",
		config.Data.Queue))

	rand.Seed(time.Now().UnixNano())

	// body, err := json.Marshal(message)
	// if err != nil {
	// 	handleError(err, "Error encoding JSON")
	// }
	for messageID := range messages {
		message := []byte(messages[messageID])
		err = amqpChannel.Publish(
			config.Data.Exchange,   // exchange
			config.Data.RoutingKey, // RoutingKey
			false,                  //mandatory - we don't care if there I no queue
			false,                  // immediate - we don't care if there is no consumer on the queue
			amqp.Publishing{
				DeliveryMode: 1,
				ContentType:  "application/json",
				Body:         message,
				Timestamp:    time.Now(),
			})

		if err != nil {
			Log.Error("Error publishing message: %s", err)
		}
	}

}
func Publisher(config Rabbit, message []byte) {

	// Create connection
	conn, err := amqp.Dial(ConnectionURL(config))
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()
	// Create channel
	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")
	// Close chanel after work
	defer amqpChannel.Close()
	_, queueDeclareErr := amqpChannel.QueueDeclare(
		config.Data.Queue, // name
		true,              // Passive
		false,             //Durable
		false,             //AutoDelete
		false,             //Exclusive
		nil,               //Arguments
	)
	handleError(queueDeclareErr, fmt.Sprintf(
		"Could not declare `%s` queue",
		config.Data.Queue))

	rand.Seed(time.Now().UnixNano())

	// body, err := json.Marshal(message)
	// if err != nil {
	// 	handleError(err, "Error encoding JSON")
	// }

	err = amqpChannel.Publish(
		config.Data.Exchange,   // exchange
		config.Data.RoutingKey, // RoutingKey
		false,                  //mandatory - we don't care if there I no queue
		false,                  // immediate - we don't care if there is no consumer on the queue
		amqp.Publishing{
			DeliveryMode: 2,
			ContentType:  "application/json",
			Body:         message,
			Timestamp:    time.Now(),
		})

	if err != nil {
		Log.Error("Error publishing message: %s", err)
	}
}

func Consumer(config Rabbit, handler func(d amqp.Delivery) bool, concurrency int) {
	// Create connection
	conn, err := amqp.Dial(ConnectionURL(config))
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()
	// Create channel
	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")
	// Close chanel after work
	defer amqpChannel.Close()
	// Declare Exchange
	exchangeDeclareErr := amqpChannel.ExchangeDeclare(
		config.Data.Exchange,     // name
		config.Data.ExchangeType, // type
		true,                     // durable
		false,                    // auto-deleted
		false,                    // internal
		false,                    // no-wait
		nil,                      // arguments
	)
	handleError(exchangeDeclareErr, fmt.Sprintf(
		"Can't declare a exchange %s type %s",
		config.Data.Exchange,
		config.Data.ExchangeType))
	_, queueDeclareErr := amqpChannel.QueueDeclare(
		config.Data.Queue, // name
		true,              // Passive
		false,             //Durable
		false,             //AutoDelete
		false,             //Exclusive
		nil,               //Arguments
	)
	handleError(queueDeclareErr, fmt.Sprintf(
		"Could not declare `%s` queue",
		config.Data.Queue))
	// bind the queue to the routing key
	queueBindErr := amqpChannel.QueueBind(
		config.Data.Queue,      // name
		config.Data.RoutingKey, // exchange
		config.Data.Exchange,   // key
		false,                  // noWait
		nil,                    //args
	)
	handleError(queueBindErr, fmt.Sprintf(
		"Could not bind `%s` queue routing key %s",
		config.Data.Queue,
		config.Data.RoutingKey,
	))
	// prefetch 4x as many messages as we can handle at once
	prefetchCount := concurrency * 4
	err = amqpChannel.Qos(prefetchCount, 0, false)
	handleError(err, "Could not configure QoS")

	messageChannel, err := amqpChannel.Consume(
		config.Data.Queue, // queue
		"",                // consumer
		false,             // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	handleError(err, "Could not register consumer")
	stopChan := make(chan bool, concurrency)
	for i := 0; i < concurrency; i++ {
		Log.Info("Processing messages on thread ", i)

		go func() {
			Log.Info("Consumer ready, PID: ", os.Getpid())
			for d := range messageChannel {
				// if tha handler returns true then ACK, else NACK
				// the message back into the rabbit queue for
				// another round of processing
				if handler(d) {
					_ = d.Ack(false)
					//log.Printf("Acknowledged message")
				} else {
					_ = d.Nack(false, true)
					//Log.Error("Error acknowledging message : ", err)
				}
			}
		}()

		// Остановка для завершения программы

	}
	<-stopChan
}

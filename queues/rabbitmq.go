package queues

import (
	"context"
	"encoding/json"
	"log"
	"mail-qq/utils"
	"os"

	"net/smtp"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ() {
	Rabbitize(initExchange)
	Rabbitize(initQueue)
	go Rabbitize(SendEmail)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func initExchange(ch *amqp.Channel) (error, string) {
	err := ch.ExchangeDeclare(
		"app.events", // exchange name
		"topic",      // exchange type
		true,         // durable
		false,        // auto-delete
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	message := ""

	if err != nil {

		message = "Could not create an exchange"
	}

	return err, message
}

func initQueue(ch *amqp.Channel) (error, string) {
	emailQueue, err := ch.QueueDeclare(
		"email_queue",
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return err, "Could not declare a queue"
	}

	err = ch.QueueBind(
		emailQueue.Name, //queue name
		"email.#",       // binding pattern
		"app.events",    // binding name
		false,           // no-wait
		nil,             // table
	)

	if err != nil {
		return err, "Could not bind a queue to the exchange"
	}

	return nil, ""
}

func Rabbitize(c func(*amqp.Channel) (error, string)) {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err, msg := c(ch)
	failOnError(err, msg)
}

func Publish(msg []byte) {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	ctx := context.Background()

	err = ch.PublishWithContext(
		ctx,             // context
		"app.events",    // exchange name
		"email.created", // routing key, follows the binding pattern email.#
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{ // msg
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         msg,
		},
	)
	failOnError(err, "Could not publish a message")

}

func SendEmail(ch *amqp.Channel) (error, string) {

	messages, err := ch.Consume(
		"email_queue",
		"",    // consumer tag
		false, // auto-ack (IMPORTANT: false)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	for message := range messages {

		var mailData utils.SendEmailRequest
		if err := json.Unmarshal(message.Body, &mailData); err != nil {
			log.Println(err)
			continue
		}

		from := os.Getenv("MAILER_EMAIL")
		password := os.Getenv("PASSWORD")
		smtpHost := os.Getenv("SMTP_HOST")
		smtpPort := os.Getenv("SMTP_PORT")

		auth := smtp.PlainAuth("", from, password, smtpHost)

		message := []byte("From: " + from + "\r\n" +
			"To: " + mailData.To + "\r\n" +
			"Subject: " + mailData.Subject + "\r\n" +
			"\r\n" +
			mailData.Body)

		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{mailData.To}, message)
		if err != nil {
			log.Fatalf("Failed to send email: %s", err)
			continue
		}

		log.Println("The email has been sent successfully.")
	}

	return nil, ""
}

package modules

import (
	"fmt"

	"github.com/streadway/amqp"
)

func SendInvoiceToQueue(ch *amqp.Channel, message string, emailTo string, nameTo string) error {
	var mapLog = make(map[string]string)
	mapLog["message"] = message
	mapLog["email"] = emailTo
	mapLog["name"] = nameTo

	jsonMessage := ConvertMapStringToJSON(mapLog)

	errP := ch.Publish(
		"",                                    // exchange
		"EMAIL_INVOICE", // online logging queue name
		false,                                 // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(jsonMessage),
		})

	if errP != nil {
		fmt.Println("Failed to submit to online logging.")

		return errP
	} else {
		fmt.Println("Success to submit to online logging.")

		return nil
	}
}


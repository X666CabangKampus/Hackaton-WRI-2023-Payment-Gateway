package emailservice

import (
	"backend-hacktober/modules"
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/ratelimit"
)

type logData struct {
	Datetime     string `bson:"datetime"`
	LogLevel     string `bson:"loglevel"`
	TraceID      string `bson:"traceid"`
	Application  string `bson:"application"`
	Module       string `bson:"module"`
	Function     string `bson:"function"`
	Identity     string `bson:"identity"`
	RemoteIP     string `bson:"remoteip"`
	Message      string `bson:"message"`
	ErrorMessage string `bson:"errormessage"`
	Status       string `bson:"status"`
}

var cxM context.Context

var connRabbit *amqp.Connection
var channelLog *amqp.Channel
var queueLog = "TRCV_MONGO_LOG"
var tps = 10
var qosCount = 5

func processTheMessage(queueMessage string) {

	mapQueueMessage := modules.ConvertJSONStringToMap("", queueMessage)

	email := modules.GetStringFromMapInterface(mapQueueMessage, "email")
	name := modules.GetStringFromMapInterface(mapQueueMessage, "name")
	message := modules.GetStringFromMapInterface(mapQueueMessage, "message")

	modules.SendActivationMail(email, name, message)
}

func processQueue() {

	theRateLimit := ratelimit.New(tps)

	errQ := channelLog.Qos(qosCount, 0, false)
	if errQ != nil {
		//modules.DoLog("INFO", "", "Transceiver9POINTS", "processQueue",
		//	"Failed to make rabbitmq QOS to "+strconv.Itoa(qosCount), true, errQ)

		panic(errQ)
	}

	// consume
	messageTransmitter, err := channelLog.Consume(
		"EMAIL_INVOICE", // queue
		"",              // consumer
		false,           // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		//modules.DoLog("INFO", "", "Transceiver9POINTS", "processQueue",
		//	"Failed to consume queue "+queueLog+". Error occured.", true, err)
	} else {
		forever := make(chan bool)

		for d := range messageTransmitter {
			//modules.DoLog("INFO", "", "Transceiver9POINTS", "processQueue",
			//	"Receiving message: "+string(d.Body), false, nil)

			// do the process with rateLimit transaction per second
			//modules.DoLog("INFO", "", "Transceiver9POINTS", "processQueue",
			//	"Do processing the incoming message from queue "+queueLog, false, nil)
			theRateLimit.Take()
			theRateLimit.Take()
			theRateLimit.Take()
			theRateLimit.Take()
			theRateLimit.Take()

			queueMessage := string(d.Body)
			go processTheMessage(queueMessage)

			errx := d.Ack(false)

			if errx != nil {
				//modules.DoLog("DEBUG", "", "Transceiver9POINTS", "readQueue",
				//	"Failed to acknowledge manually message: "+string(d.Body)+". STOP the transceiver.", false, nil)

				os.Exit(-1)
			}

			//modules.DoLog("DEBUG", "", "Transceiver9POINTS", "readQueue",
			//"Done Processing queue message "+string(d.Body)+". Sending ack to rabbitmq.", false, nil)
		}

		fmt.Println("[*] Waiting for messages. To exit press CTRL-C")
		<-forever
	}
}

func startReceiver() {

	var errI error
	channelLog, errI = connRabbit.Channel()
	if errI != nil {

		panic(errI)
	}
	defer channelLog.Close()

	// Thread utk check status chIncoming and reconnect if failed. Do it in different treads run forever
	go func() {
		for {
			c := make(chan int)

			theCheckedQueue, errC := channelLog.QueueInspect(queueLog)

			if errC != nil {
				channelLog, _ = connRabbit.Channel()
			} else {
				if theCheckedQueue.Consumers == 0 {

					channelLog, _ = connRabbit.Channel()
					_ = channelLog.Qos(qosCount, 0, false)
				}
			}

			sleepDuration := 10 * time.Minute
			time.Sleep(sleepDuration)

			<-c
		}
	}()

	processQueue()
}

func EmailQueueReceiver() {
	// Load configuration file
	modules.InitiateGlobalVariables(false)
	runtime.GOMAXPROCS(4)

	// Initiate RabbitMQ
	var errRabbit error
	connRabbit, errRabbit = amqp.Dial("amqp://" + modules.MapConfig["rabbitUser"] + ":" + modules.MapConfig["rabbitPass"] + "@" + modules.MapConfig["rabbitHost"] + ":" + modules.MapConfig["rabbitPort"] + "/" + modules.MapConfig["rabbitVHost"])
	if errRabbit != nil {
		// modules.DoLog("INFO", "", "SMPP20", "main",
		// "Failed to connect to RabbitMQ server. Error", true, errRabbit)
		log.Println("Failed to connect to RabbitMQ server. Error : ", errRabbit)
	} else {
		// modules.DoLog("INFO", "", "SMPP20", "main",
		// "Success to connect to RabbitMQ server.", false, nil)
		log.Println("Success to connect to RabbitMQ server.")
	}
	defer connRabbit.Close()

	startReceiver()
}

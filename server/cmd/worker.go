package main

import (
	"encoding/json"
	"fmt"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/event_queue"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/notifier"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
)

func loop(notifierConfig *notifier.Config, eventQueueConfig *event_queue.Config) {
	emknCourseMail := os.Getenv("EMKN_COURSE_MAIL")
	emknCoursePassword := os.Getenv("EMKN_COURSE_PASSWORD")
	rabbitLogin := os.Getenv("RABBIT_LOGIN")
	rabbitPassword := os.Getenv("RABBIT_PASSWORD")

	mailer := notifier.New(notifierConfig, emknCourseMail, emknCoursePassword)

	url := fmt.Sprintf("amqp://%s:%s@%s:%s", rabbitLogin, rabbitPassword, eventQueueConfig.Address, eventQueueConfig.Port)

	conn, _ := amqp.Dial(url)
	amqpChannel, _ := conn.Channel()
	queue, _ := amqpChannel.QueueDeclare("Email", true, false, false, false, nil)
	messageChannel, _ := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	for d := range messageChannel {
		msg := &models.Message{}
		err := json.Unmarshal(d.Body, msg)
		if err != nil {
			logrus.Errorf("Unmarshal has err %s", err)
		}
		if err := mailer.SendEmail(*msg); err != nil {
			logrus.Errorf("Mailer has err %s on input %s", err, msg.Body)
		}
		if err := d.Ack(false); err != nil {
			logrus.Errorf("Error acknowledging message : %s", err)
		} else {
			logrus.Debugf("Acknowledged message")
		}
	}

}

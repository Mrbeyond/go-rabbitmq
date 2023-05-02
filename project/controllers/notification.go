package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"rabbit/rabbit_instance"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
func NotifySingleUser(c *gin.Context) {

	defer func() {
		if r := recover(); r != nil {
			c.String(http.StatusInternalServerError, fmt.Sprint(r))
		}
	}()

	user := c.Param("user")
	conn, err := rabbit_instance.GetConnectionFromPool(os.Getenv("RABBITMQ_URL"))
	panicIfError(err)
	defer rabbit_instance.ReturnConnectionToPool(conn)
	ch, err := conn.Channel()
	panicIfError(err)

	q, err := ch.QueueDeclare(
		user,
		false,
		false,
		false,
		false,
		nil,
	)

	cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = ch.PublishWithContext(
		cntx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("Hi user %s", user)),
		},
	)
	c.String(http.StatusOK, fmt.Sprintf("Notification sent to user %s ", user))
}

func Broadcast(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.String(http.StatusInternalServerError, fmt.Sprint(r))
		}
	}()

	user := c.DefaultQuery("user", "")
	conn, err := rabbit_instance.GetConnectionFromPool(os.Getenv("RABBITMQ_URL"))
	panicIfError(err)
	defer rabbit_instance.ReturnConnectionToPool(conn)
	ch, err := conn.Channel()
	panicIfError(err)

	q, err := ch.QueueDeclare(
		"all",
		false,
		false,
		false,
		false,
		nil,
	)

	cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = ch.PublishWithContext(
		cntx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("Hi user %s, \n the body is %s", user, randomString())),
		},
	)
	c.String(http.StatusOK, fmt.Sprintf("Notification sent to user %s ", user))
}

func randomString() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, 16)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

package rabbit_instance

import (
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	connPool chan *amqp.Connection
	mutex    sync.Mutex
)

// InitConnPool initializes pool of rabbitMQ connections to avoid creation of new
// connection on every http request. One of the connections could be used
// by calling GetConnectionFromPool()
func InitConnPool(numOfWorkers int, amqpAddr string) {
	connPool = make(chan *amqp.Connection, numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		conn, err := amqp.Dial(amqpAddr)
		if err != nil {
			log.Fatalf("Error connecting to rabbitMQ serve >>> %v", err)
		}

		connPool <- conn
	}
}

// GetConnectionFromPool gets one from the avaible connections on the connPool channel.
// If no valid connection is available, a new instance is returned
func GetConnectionFromPool(amqpAddr string) (*amqp.Connection, error) {
	mutex.Lock()
	defer mutex.Unlock()
	select {
	case conn := <-connPool:
		if conn == nil {
			return nil, amqp.ErrClosed
		}
		return conn, nil
	default:
		return amqp.Dial(amqpAddr)
	}
}

// ReturnConnectionToPool returns the current rabbitMQ connection back to the connPool channel
// after the current http request ends.
func ReturnConnectionToPool(conn *amqp.Connection) {
	mutex.Lock()
	connPool <- conn
	mutex.Unlock()
}

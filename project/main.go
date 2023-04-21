package main

import (
	"log"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	log.Println("\n\n\n\n\n\n\n\n\n\n\n\n\n Started \n\n\n\n\n\n\n\n\n ")
	Sender()
	Receiver()
	log.Println("\n\n\n\n\n\n\n\n\n\n\n\n\n Ends \n\n\n\n\n\n\n\n\n ")

}

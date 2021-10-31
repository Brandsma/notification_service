package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/globalsign/mgo"

	ctx "context"

	log "github.com/Sirupsen/logrus"
)

// ListenForGracefulShutdown listens for SIGINT and SIGTERM (and standard os interrupt) and consequently shuts down gracefully.
// It is a function that can be called after starting a HTTP server.
// Calling the same signal again immediately stops the service.
func ListenForGracefulShutdown(srv *http.Server) {
	/* GRACEFUL SHUTDOWN */
	// Make a channel that listens for SIGINT and SIGTERM
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Blocking call that waits until SIGINT or SIGTERM is submitted to the done channel
	<-done
	log.Infof("Graceful shutdown of service started")

	// If another terminating signal is received, kill immediately
	go func() {
		select {
		case s := <-done:
			log.Fatalf("received signal %s: terminating", s)
		}
	}()

	// Timeout the shutdown if it takes longer than 3 seconds
	curContext, cancel := ctx.WithTimeout(ctx.Background(), 3*time.Second)
	// Release resources of context after use
	defer cancel()

	// Do a graceful shutdown
	if err := srv.Shutdown(curContext); err != nil {
		log.Fatalf("Graceful shutdown failed: %+v", err)
	}

	log.Infof("Service exited properly")
}


// ConnectToDatabase returns a database session to the mongo database. Don't forget to defer db.Close()
func ConnectToDatabase(host string, port string) *mgo.Session {
	log.Infof("Mongo database set to " + host + ":" + port)
	// Using retry because it might not succeed immediately
	db, err := Retry(5, time.Second, func() (*mgo.Session, error) {
		log.Infof("Attempting to dial mongoDB server...")
		db, err := mgo.Dial(host + ":" + port)
		if err != nil {
			return nil, err
		}
		return db, nil
	})
	if err != nil {
		log.Errorf("timed out: cannot dial mongo ", err)
		return nil
	}
	log.Infof("Succesful database connection")

	return db
}

// Retry attempts to create a database connection several times. Every attempt the sleep timer is doubled.
func Retry(attempts int, sleep time.Duration, fn func() (*mgo.Session, error)) (*mgo.Session, error) {
	db, err := fn()
	if err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return nil, s.error
		}

		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, fn)
		}

		return nil, err

	}
	if err != nil {
		return nil, err
	}
	return db, nil
}

type stop struct {
	error
}

// SendData takes in any object, transforms it to JSON and sends it on the responseWriter
func SendData(w http.ResponseWriter, data interface{}) *AppError {
	//Marshal or convert user object back to json and write to response
	dataJSON, e := json.Marshal(data)
	if e != nil {
		return AppErrorf(500, e, "JSONParseError", "Could not marshal object for response")
	}

	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/json")
	//Write json response back to response
	w.Write(dataJSON)
	return nil
}


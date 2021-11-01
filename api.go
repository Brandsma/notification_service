package main

import (
	"net/http"
	"os"
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/Sirupsen/logrus"
)

func Start(serverHost string, serverPort string, dbHost string, dbPort string) {
	// Start the Websocket hub
	go h.run()

	// Get a connection to the database
	db := ConnectToDatabase(dbHost, dbPort, ENV_MONGODB_USER, ENV_MONGODB_PASS)
	defer db.Close() // clean up when we're done

	registerHandlers(db)

	/* SETUP SERVER */
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", serverPort),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	/* SERVER START */
	log.Infoln("listening on " + serverHost + ":" + serverPort)
	go log.Fatal(srv.ListenAndServe())

	ListenForGracefulShutdown(srv)
}

func registerHandlers(db *mgo.Session) {
	r := mux.NewRouter()

	// Add a verify email point
	hSendNotification := Adapt(AppHandler(SendNotification), WithDB(db))
	hDeleteNotification := Adapt(AppHandler(SendNotification), WithDB(db))
	hListenForNotification := Adapt(AppHandler(ListenForNotification), WithDB(db))
	hWebsocketNotifications := AppHandler(serveWs)

	r.Handle("/api/notification/notification", context.ClearHandler(hSendNotification)).Methods("POST")
	r.Handle("/api/notification/notification", context.ClearHandler(hListenForNotification)).Methods("GET")
	r.Handle("/api/notification/notification", context.ClearHandler(hDeleteNotification)).Methods("DELETE")
	r.Handle("/api/notification/ws/notification/{userID}", context.ClearHandler(hWebsocketNotifications))

	log.Infof("Serving API end-points")

	// // Set standard response handlers on all routers
	// r.MethodNotAllowedHandler = handler.NotAllowedHandler{}
	// r.NotFoundHandler = handler.NotFoundHandler{}

	// Delegate all of the HTTP routing and serving to the gorilla/mux router.
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
}


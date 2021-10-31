package main

import (
	"net/http"
	"github.com/gorilla/context"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

type Notification struct {
	Recipient string `json:"recipient" bson:"recipient"`
	Label string `json:"label" bson:"label"`
	Message string `json:"message" bson:"message"`
}

// SendNotification waits for a notification to be received and adds it to the message queue
func SendNotification(w http.ResponseWriter, r *http.Request) *AppError {
	// Get the database reference
	db := context.Get(r, "database").(*mgo.Session)

	// Get message from body
	notification := &Notification{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&notification); err != nil {
		return AppErrorf(http.StatusInternalServerError, err, "JSONParseError", "Could not decode request into notification struct")
	}

	// Add message to queue
	log.Debugf("%v", notification)

	// Get collection reference
	nc := db.DB("schrodat").C("notifications")

	// TODO: Insert in queue
	// TODO: Send to websocket of user

	// Insert notification into database if it cannot be send to the user right now
	if err := nc.Insert(&notification); err != nil {
		return AppErrorf(http.StatusInternalServerError, err, "DatabaseError", "Could not insert notification into database: %v", err)
	}


	return nil
}

type NotificationRequest struct {
	Recipient string `json:"recipient" bson:"recipient"`
}

// ListenForNotification is a websocket endpoint that the frontend can access to listen for old notifications
func ListenForNotification(w http.ResponseWriter, r *http.Request) *AppError {
	// Get the database reference
	db := context.Get(r, "database").(*mgo.Session)

	notificationRequest := &NotificationRequest{}
	if err := json.NewDecoder(r.Body).Decode(&notificationRequest); err != nil {
		return AppErrorf(http.StatusInternalServerError, err, "JSONParseError", "Could not decode request into notification request struct")
	}

	// Retrieve stored notifications from server for specific user
	notificationList := []Notification{}
	// TODO: Get variables from env
	if err := db.DB("schrodat").C("notifications").Find(bson.M{"recipient": notificationRequest.Recipient}).All(&notificationList); err != nil {
		return AppErrorf(404, err, "NotificationDoesNotExistError", "Notification has not been found in database")
	}

	if err := SendData(w, notificationList); err != nil {
		return err
	}

	return nil
}

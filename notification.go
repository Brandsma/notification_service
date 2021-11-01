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
	// db := context.Get(r, "database").(*mgo.Session)

	// Get message from body
	notification := &Notification{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&notification); err != nil {
		return AppErrorf(http.StatusInternalServerError, err, "JSONParseError", "Could not decode request into notification struct")
	}

	// Add message to queue
	log.Debugf("%v", notification)


	// TODO: Maybe insert in message queue
	// Send to websocket of user if possible
	if err := NewNotification(notification.Recipient, []byte(notification.Message)); err != nil {
		return AppErrorf(http.StatusInternalServerError, err, "SendNotificationError", "Unable to send notifaction to %s", notification.Recipient)
	}

	// Get collection reference
	// nc := db.DB(ENV_MONGODB_DB_NAME).C(ENV_MONGODB_NOTIFICATION_COLLECTION)
	// TODO: Insert the notification into a database like mongodb if the user cannot be reached right now
	// // Insert notification into database if it cannot be send to the user right now
	// if err := nc.Insert(&notification); err != nil {
	// 	return AppErrorf(http.StatusInternalServerError, err, "DatabaseError", "Could not insert notification into database: %v", err)
	// }


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

	// Get collection reference
	nc := db.DB(ENV_MONGODB_DB_NAME).C(ENV_MONGODB_NOTIFICATION_COLLECTION)

	// Retrieve stored notifications from server for specific user
	notificationList := []Notification{}
	// TODO: Get variables from env
	if err := nc.Find(bson.M{"recipient": notificationRequest.Recipient}).All(&notificationList); err != nil {
		return AppErrorf(404, err, "NotificationDoesNotExistError", "Notification has not been found in database")
	}

	if err := SendData(w, notificationList); err != nil {
		return err
	}

	return nil
}

// DeleteNotifications will remove all notifications for a specific recipient. This can be called after the notifications have been received
func DeleteNotifications(w http.ResponseWriter, r *http.Request) *AppError {
	// Get the database reference
	db := context.Get(r, "database").(*mgo.Session)

	notificationRequest := &NotificationRequest{}
	if err := json.NewDecoder(r.Body).Decode(&notificationRequest); err != nil {
		return AppErrorf(http.StatusInternalServerError, err, "JSONParseError", "Could not decode request into notification request struct")
	}

	// Get collection reference
	nc := db.DB(ENV_MONGODB_DB_NAME).C(ENV_MONGODB_NOTIFICATION_COLLECTION)

	// Remove all notifications
	_, _ = nc.RemoveAll(bson.M{"recipient": notificationRequest.Recipient});

	return nil
}

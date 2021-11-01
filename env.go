package main

import (
	"os"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

const (
	MONGODB_USER = "MONGODB_USER"
	MONGODB_PASS="MONGODB_PASS"
	MONGODB_HOST="MONGODB_HOST"
	MONGODB_PORT="MONGODB_PORT"
	MONGODB_DB_NAME="MONGODB_DB_NAME"
	MONGODB_NOTIFICATION_COLLECTION="MONGODB_NOTIFICATION_COLLECTION"

	NOTIFICATION_HOST="NOTIFICATION_HOST"
	NOTIFICATION_PORT="NOTIFICATION_PORT"
)

var ENV_MONGODB_USER= os.Getenv(MONGODB_USER)
var ENV_MONGODB_PASS= os.Getenv(MONGODB_PASS)
var ENV_MONGODB_HOST= os.Getenv(MONGODB_HOST)
var ENV_MONGODB_PORT= os.Getenv(MONGODB_PORT)
var ENV_MONGODB_DB_NAME= os.Getenv(MONGODB_DB_NAME)
var ENV_MONGODB_NOTIFICATION_COLLECTION= os.Getenv(MONGODB_NOTIFICATION_COLLECTION)
var ENV_NOTIFICATION_HOST= os.Getenv(NOTIFICATION_HOST)
var ENV_NOTIFICATION_PORT= os.Getenv(NOTIFICATION_PORT)

// VerifyProfileEnvs checks whether each of the environment variables returned a non-empty value
func VerifyProfileEnvs() error {
	if ENV_MONGODB_USER=="" {
		return ErrEnvironment(MONGODB_USER, ENV_MONGODB_USER)
	} else if ENV_MONGODB_PASS=="" {
		return ErrEnvironment(MONGODB_PASS, ENV_MONGODB_PASS)
	} else if ENV_MONGODB_HOST=="" {
		return ErrEnvironment(MONGODB_HOST, ENV_MONGODB_HOST)
	} else if ENV_MONGODB_PORT=="" {
		return ErrEnvironment(MONGODB_PORT, ENV_MONGODB_PORT)
	} else if ENV_MONGODB_DB_NAME=="" {
		return ErrEnvironment(MONGODB_DB_NAME, ENV_MONGODB_DB_NAME)
	} else if ENV_MONGODB_NOTIFICATION_COLLECTION=="" {
		return ErrEnvironment(MONGODB_NOTIFICATION_COLLECTION, ENV_MONGODB_NOTIFICATION_COLLECTION)
	} else if ENV_NOTIFICATION_HOST=="" {
		return ErrEnvironment(NOTIFICATION_HOST, ENV_NOTIFICATION_HOST)
	} else if ENV_NOTIFICATION_PORT=="" {
		return ErrEnvironment(NOTIFICATION_PORT, ENV_NOTIFICATION_PORT)
	}
	return nil
}

func init() {
	errProfile := VerifyProfileEnvs()

	err := ReturnFirstErr(errProfile)
	if err != nil {
		log.Fatalln(err)
	}
}

// EnvironmentError is raised when an environment variable that is verified turns up empty
type EnvironmentError struct {
	Variable string
	Value    string
}

func (err *EnvironmentError) Error() string {
	return fmt.Sprintf("EnvironmentError: %v with value '%v' is non-valid", err.Variable, err.Value)
}

func (err *EnvironmentError) Unwrap() error {
	return nil
}

// ErrEnvironment is used to raise a EnvironmentError
func ErrEnvironment(variable, value string) *EnvironmentError {
	return &EnvironmentError{Variable: variable, Value: value}
}

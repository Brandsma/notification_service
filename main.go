package main

func main() {
	// Start the http server
	// ServerHost, ServerPort, DBHost, DBPort
	Start(ENV_NOTIFICATION_HOST, ENV_NOTIFICATION_PORT, ENV_MONGODB_HOST, ENV_MONGODB_PORT)
}

############################
# STEP 1 build executable binary
############################
# golang alpine 1.12
FROM golang:1.17-alpine as builder

WORKDIR $GOPATH/src/notification_service/src
COPY . .

# Fetch dependencies.
RUN go get -d -v

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/notification .

############################
# STEP 2 build a small image
############################
FROM scratch

# Copy our static executable
COPY --from=builder /go/bin/notification /go/bin/notification

# Open this docker on port 9003
EXPOSE 9003

# Run the hello binary.
ENTRYPOINT ["/go/bin/notification"]

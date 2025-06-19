############################
# STEP 1 build executable binary
############################
FROM golang:1.24-alpine AS buildergo
ARG BUILD_VERSION
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates
# Create appuser
RUN adduser -D -g '' appuser
# Copy the go source
COPY ./cmd/ $GOPATH/src/github.com/stevenweathers/hatch-messaging-service/cmd/
COPY ./domain/ $GOPATH/src/github.com/stevenweathers/hatch-messaging-service/domain/
COPY ./internal/ $GOPATH/src/github.com/stevenweathers/hatch-messaging-service/internal/
COPY ./*.go $GOPATH/src/github.com/stevenweathers/hatch-messaging-service/
COPY ./go.mod $GOPATH/src/github.com/stevenweathers/hatch-messaging-service/
COPY ./go.sum $GOPATH/src/github.com/stevenweathers/hatch-messaging-service/
# Set working dir
WORKDIR $GOPATH/src/github.com/stevenweathers/hatch-messaging-service/
# Fetch dependencies.
RUN go mod download
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -ldflags "-X main.version=$BUILD_VERSION" -o /go/bin/hms
############################
# STEP 2 build a small image
############################
FROM scratch
# Import from builder.
COPY --from=buildergo /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=buildergo /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=buildergo /go/bin/hms /go/bin/hms
# Use an unprivileged user.
USER appuser

# Run the hatch messaging service binary.
ENTRYPOINT ["/go/bin/hms", "serve"]
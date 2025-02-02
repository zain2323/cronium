FROM golang:1.23.5-alpine3.21 AS build
LABEL authors="zain-siddiqui"

WORKDIR /app

# copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy the rest of the application
COPY . .

# building the binary
RUN go build -o userservice-exec main.go

FROM alpine:3.21 AS runtime

RUN apk add --no-cache curl

WORKDIR /app

# copying the build files from previous stage
COPY --from=build /app .

RUN curl -fsSL \
        https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
        sh

EXPOSE 8081

# run the binary
ENTRYPOINT ["sh", "-c", "goose postgres postgres://user-service-user:password@user-service-db:5432/user-service-db up --dir sql/schema/ && /app/userservice-exec"]
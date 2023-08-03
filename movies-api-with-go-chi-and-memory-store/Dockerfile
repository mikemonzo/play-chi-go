FROM golang:1.20.7

RUN export PATH="$PATH:$(go env GOPATH)/bin"

RUN go install "github.com/kelseyhightower/envconfig@latest"
RUN go install "github.com/go-chi/chi"
RUN go install "github.com/go-chi/chi/v5"
RUN go install "github.com/go-chi/render"
RUN go install "github.com/google/uuid"
RUN go install "github.com/joho/godotenv"


RUN export PATH="$PATH:$(go env GOPATH)/bin"

WORKDIR /usr/src/app
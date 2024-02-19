FROM golang:1.21 AS BUILD

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

FROM build AS final

COPY . .
RUN go build -v -o tag-onboarding-api cmd/ps-tag-onboarding/main.go

CMD ["./tag-onboarding-api"]
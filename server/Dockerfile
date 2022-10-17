FROM golang:1.18.2-alpine

RUN apk update && apk upgrade && apk add --no-cache gcc musl-dev

WORKDIR /app
COPY . .
RUN go build

EXPOSE 3001 

# CMD ["go", "run", "main.go"]
CMD ["./calorie-tracker"]
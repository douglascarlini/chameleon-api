FROM golang:1.20-alpine
RUN mkdir /app
COPY src /app
WORKDIR /app
RUN go build -o main .
CMD [ "/app/main" ]
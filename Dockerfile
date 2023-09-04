FROM golang:alpine3.18 as build
WORKDIR /server
COPY . /server
RUN go build -o /torch-server

FROM alpine
COPY --from=build ./torch-server ./
COPY --from=build ./server/.env ./
EXPOSE 8080
ENTRYPOINT ["/torch-server"]
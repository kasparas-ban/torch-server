# FROM golang:alpine3.18 as build
# WORKDIR /server
# COPY . /server
# RUN go build -o /torch-server

# FROM alpine
# COPY --from=build ./torch-server ./
# COPY --from=build ./server/.env ./
# EXPOSE 3001

# Installing certificates for AWS RDS connection
# RUN wget https://truststore.pki.rds.amazonaws.com/eu-north-1/eu-north-1-bundle.pem
# RUN mv eu-north-1-bundle.pem /usr/local/share/ca-certificates/eu-north-1-bundle.crt
# RUN update-ca-certificates

# ENTRYPOINT ["/torch-server", "-prod"]

# =====================
# For Fly.io
# =====================

FROM golang:latest as build
WORKDIR /server
COPY . /server
RUN export GOPROXY=https://goproxy.io,direct
RUN export GOPROXY=https://goproxy.io,direct && go build -o ./torch-server
EXPOSE 3003
ENTRYPOINT ["./torch-server", "-prod"]
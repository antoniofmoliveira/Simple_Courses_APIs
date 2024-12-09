# Simple User, Course and Category API

## ***in progress***

using Fullcycle's Goexpert Course example to expand the API

very dificult to find information about flatbuffers vectors and grpc streams

## Several ways to implement the API

- using golang workspace to share code between sub projects
- use of mariadb or sqlite databases

- generating certificates
- authentication and authorization
- use jwttoken
- multiple middlewares
- json api
- grpc api
- flatbuffer api
- graphql api
- webserver ? TODO!
- gui client ? TODO!
- test apis ? TODO!

## Generate self-signed certificates

```bash
openssl req -x509 -newkey rsa:4096 -sha256 -days 30 -nodes -keyout server_key.pem -out server_cert.pem -subj '/C=XX/ST=State/L=City/O=Organization/OU=Section/CN=hostname' -addext 'subjectAltName=DNS:hostname,DNS:localhost'
```

put `server_key.pem` and `server_cert.pem` in `x509` folder on **server**

copy `server_cert.pem` to **client** as `ca_cert.pem` in `x509` folder

update `hostname` in `client.go`

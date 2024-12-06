# gRPC Client and Server with TLS Authentication in Golang using streams

## Generate self-signed certificates

```bash
openssl req -x509 -newkey rsa:4096 -sha256 -days 30 -nodes -keyout server_key.pem -out server_cert.pem -subj '/C=XX/ST=State/L=City/O=Organization/OU=Section/CN=hostname' -addext 'subjectAltName=DNS:hostname,DNS:localhost'
```

put `server_key.pem` and `server_cert.pem` in `x509` folder on **server**

copy `server_cert.pem` to **client** as `ca_cert.pem` in `x509` folder

update `hostname` in `client.go`

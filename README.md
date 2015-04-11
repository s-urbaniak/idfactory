# idfactory

## Installation

```
$ go get github.com/s-urbaniak/idfactory
```

## Usage

### Server

To start a server on port 8080 and a secret `password`:

```
$ ./idfactory -addr=":8080" -secret="password"
```

All generated UUIDs will be signed using the given secret.

### Client

To create a signed UUID, send an empty POST request:

```
$ curl -v -XPOST localhost:8080 | python -mjson.tool
> POST / HTTP/1.1

< HTTP/1.1 201 Created
< Content-Type: application/json
< Location: /633b5398-7233-4f53-975e-65e0ac39dbe6:YezlAZEYbbrfHSRHSUy7Vy1N/JEEMQtATzsQRCxnkSI=

{
    "id": "633b5398-7233-4f53-975e-65e0ac39dbe6",
    "signed": "633b5398-7233-4f53-975e-65e0ac39dbe6:YezlAZEYbbrfHSRHSUy7Vy1N/JEEMQtATzsQRCxnkSI="
}
```

The server will respond with a JSON response having an `id` field
containing the actual and `signed` filed containing the signed UUID.
The server will also respond with a `Location` header
containing a relative link for future validation.
Note that the returned `Location` header is relative
(see https://tools.ietf.org/html/rfc7231#section-7.1.2).

To validate a given signed UUID send a GET request:

```
$ curl -v localhost:8080/633b5398-7233-4f53-975e-65e0ac39dbe6:YezlAZEYbbrfHSRHSUy7Vy1N/JEEMQtATzsQRCxnkSI=
> GET /633b5398-7233-4f53-975e-65e0ac39dbe6:YezlAZEYbbrfHSRHSUy7Vy1N/JEEMQtATzsQRCxnkSI= HTTP/1.1
> Host: localhost:8080

< HTTP/1.1 204 No Content
```

The server will respond with a 204 status code and an empty body
or with a status code 412 with an empty body
if the signed UUID could not be validated:

```
$ curl -v localhost:8080/foo:bar
> GET /foo:bar HTTP/1.1

< HTTP/1.1 412 Precondition Failed
```

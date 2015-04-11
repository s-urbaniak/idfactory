# idfactory

idfactory is a small http service for generating signed UUIDs.

A common (anti)pattern is to provide endpoints
which accept POST requests to create new resources
like new entries in a database.

This has the severe disadvantage of not being idempotent
and can cause duplicate created resources.
An alternative is to provide idempotent PUT endpoints
with the cost of a second round-trip.

This service generates V4 UUIDs according to RFC4122
which can be used as primary keys or unique IDs.

A client first invokes this service
to create a cryptographically signed unique ID:

1. Request
```
POST /
```
1. Response
```
HTTP/1.1 201 Created
Content-Type: application/json
Location: /633b5398-7233-4f53-975e-65e0ac39dbe6:YezlAZEYbbrfHSRHSUy7Vy1N/JEEMQtATzsQRCxnkSI=

{
    "id": "633b5398-7233-4f53-975e-65e0ac39dbe6",
    "signed": "633b5398-7233-4f53-975e-65e0ac39dbe6:YezlAZEYbbrfHSRHSUy7Vy1N/JEEMQtATzsQRCxnkSI="
}

```

Then the client calls the actual target service
with the pregenerated key:

2. Request
```
PUT /service/entity/633b5398-7233-4f53-975e-65e0ac39dbe6:YezlAZEYbbrfHSRHSUy7Vy1N/JEEMQtATzsQRCxnkSI=
...
```

The target service validates the generated UUID in situ
or can alternatively validate the transmitted key
by invoking a GET request towards idfactory.

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
containing the actual UUID and `signed` field containing the signed UUID.
The server will also respond with a `Location` header
containing a link for future validation.
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

# RBH AxiomXa Driver

Driver for the AxiomXa product version, produced by RBH. The server is an access control system which consists of an
"Axiom Server" which provides the vendor specific UI for managing the estate. This server also hosts any API access
available.

Axiom itself has the concept of cards and card holders, controllers (which they call networks) and card readers (which
they call devices).

AxiomXa provides three different methods of integration, each with benefits and drawbacks

1. HTTP API
2. Message Ports
3. Direct DB access

## HTTP API

The HTTP API has to be enabled by AxiomXa contractors with special licencing. By default it is HTTP only but is hosted
within IIS and it is the responsibility of the site administrator to manage the certificates and setup TLS for the API.
The API is quite basic and the documentation is NDA-walled, it's consists of what looks like generated API docs from
some internal representation, edited and placed into PDF (likely via Word).

The API is not quite restful, but does use JSON as the payload. While GET and POST verbs are used, you'll find urls like
this to get a single item:

```http request
# Get a single access level
POST /v1/accesslevel/one

{ "ID": "string" }
```

The API is authenticated using a simple `/v1/login` API which returns a bearer token to include in the Authorization
header of future requests.

There is no event or notification mechanism using this API.

## Message Ports

Invented by RBH, think _outgoing webhooks_. These are configured per AxiomXa install, typically by the people
maintaining that system. They are quite flexible but also limited. They use TCP not HTTP, they have a payload limit of
50 characters but can notify you of changes and events that happen in all sorts of areas within the AxiomXa system.

The payload of the messages is not defined at all. During commissioning the AxiomXa engineer sets up templates for the
payload that resembles `EVENTTIME CARDID EVENTTYPE ALLOWED` in the AxiomXa UI, they set up a trigger condition and a
target to send it to and the system replaces the template keywords and sends the data down the socket.

Here are some of the available template placeholder values (called inserts):

| Insert           | Description                                                                                         | Type                              |
|------------------|-----------------------------------------------------------------------------------------------------|-----------------------------------|
| `TIMESTAMP`      | Date & Time of the event, acquired from the event message.                                          | `string`, `"dd/mm/yyyy hh:MM:ss"` |
| `EVENTID`        | Identification number associated with the event.                                                    | `uint`                            |
| `EVENTDESC`      | Description of the event, acquired from the event message.                                          | `string`                          |
| `NETWORKID`      | Identification number associated with the network of the event.                                     | integer                           |
| `NETWORKDESC`    | Description of the network, associated with the event message.                                      | `string`                          |
| `NC100ID`        | Identification number associated with the NC100 of the event.                                       | integer                           |
| `NC100DESC`      | Description of the NC100, associated with the event message.                                        | `string`                          |
| `DEVICEID`       | Identification number associated with the device (RC2, IOC16, or SafeSuiteTM panel) of the event.   | integer                           |
| `DEVICEDESC`     | The description of the device (RC2, IOC16, or SafeSuiteTM panel) associated with the event message. | `string`                          |
| `CARDID`         | Identification number associated with the Card                                                      | integer                           |
| `CARDNUMBER`     | Card number associated with the event.                                                              | `uint64`                          |
| `CARDHOLDERDESC` | Name of the cardholder associated with the event.                                                   | `string`                          |
| `USAGECOUNT`     | Usage count assigned to a card.                                                                     | ?                                 |

Message Ports cannot be written to.

A _network_ in Axiom terminology refers to a card reader or door. An example `NETWORKDESC` is `"ACU-06 Comms intake"`.

## Direct DB Access

We try to avoid this as it's both the least secure and least robust, it does give us access to all the information
available to the AxiomXa system though.


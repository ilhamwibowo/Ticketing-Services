# Ticket App

## API Docs

Dokumentasi ini tidak strict, silahkan ubah sesuai yang diinginkan, yang penting ada terdokumentasi interface dari masing2 app

### HTTP APIs

| HTTP Method | Endpoint                                | Description                                       |
| ----------- | --------------------------------------- | ------------------------------------------------- |
| GET         | /seats                                  | List all seats                                    |
| POST        | /seats                                  | Create new seat                                   |
| GET         | /seats/<event_id>/<seat_number>/status  | Status check for selected seat                    |
| GET         | /events                                 | Lists all events                                  |
| GET         | /events/<event_id>/empty-seats          | Lists empty seats within an event                 |
| POST        | /book/<event_id>/<seat_number>          | Book seat <seat_number> in event <event_id>       |
| POST        | /webhook/payment                        | Webhook for payment service                       |

## How To Start

Jelaskan step by step cara menjalankan kode dari service ini, misal:

1. Ensure port X, Y, Z is not used and exposed
2. Run `docker-compose up`


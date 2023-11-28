# Client App

## API Docs

Dokumentasi ini tidak strict, silahkan ubah sesuai yang diinginkan, yang penting ada terdokumentasi interface dari masing2 app

### HTTP APIs

| HTTP Method | Endpoint                                     | Description                          |
| ----------- | -------------------------------------------- | ------------------------------------ |
| GET         | /health/                                     | Get service health check             |
| GET         | /book/                                       | Get list of bookings                 |
| GET         | /book/api/bookings/                          | Get list of bookings (API)           |
| GET         | /book/api/bookings/<str:booking_id>/refresh/ | Refresh booking status (API)         |
| GET         | /book/book/                                  | Book a seat (Web)                    |
| POST        | /book/book/                                  | Book a seat (API)                    |
| GET         | /book/api/events/                            | Get available events (API)           |
| GET         | /book/api/events/<str:event_id>/chairs/      | Get chairs status for an event (API) |
| POST        | /book/api/invoices/create/                   | Create invoice (API)                 |
| GET         | /book/api/invoices/<str:invoice_id>/         | Get invoice file (API)               |

## How To Start

Jelaskan step by step cara menjalankan kode dari service ini, misal:

1. Ensure port X, Y, Z is not used and exposed
2. Run `docker-compose up --build`
3. Hit http://localhost:X/health and see if it returns properly

## References

1.  [Quickstart: Compose and Django](https://github.com/docker/awesome-compose/tree/master/official-documentation-samples/django/#readme)
2.  [Django](https://www.djangoproject.com)

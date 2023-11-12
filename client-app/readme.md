# Client App

## API Docs

Dokumentasi ini tidak strict, silahkan ubah sesuai yang diinginkan, yang penting ada terdokumentasi interface dari masing2 app

### HTTP APIs

| HTTP Method | Endpoint   | Description              |
| ----------- | ---------- | ------------------------ |
| GET         | /v1/health | Get service health check |
| GET         | /v1/todos  | Get list of todos        |

### GRPC APIs

| Method                               | Return  | Description       |
| ------------------------------------ | ------- | ----------------- |
| GetAllTodos(page int, search string) | [ ]Todo | Get list of todos |

## How To Start

Jelaskan step by step cara menjalankan kode dari service ini, misal:

1. Ensure port X, Y, Z is not used and exposed
2. Run `docker-compose up --build`
3. Hit http://localhost:X/health and see if it returns properly

## References

 1. [Quickstart: Compose and Django](https://github.com/docker/awesome-compose/tree/master/official-documentation-samples/django/#readme)
 2. [Django](https://www.djangoproject.com)
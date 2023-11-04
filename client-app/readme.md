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
2. Run `docker-compose -f docker-compose-dev.yml up`
3. Hit http://localhost:X/health and see if it returns properly

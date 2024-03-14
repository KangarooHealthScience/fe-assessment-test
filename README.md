# Simple TODO List Web Service

## Guide

Create a simple TODO list application using nextjs. The data will be stored on the BE, and we have provided the API for that.

On the front end side, the application must be built under the following requirements:

- using nextjs
- use redux or other state management tool
- it has login form and consume the login endpoint on the backend side to authenticate
- and after user successfully logged in, user should see the list of TODO items, along with ability to create, edit, and delete TODO item.
- the TODO items will be stored on the web service. therefore the front end needs to communicate with the BE.
- the FE must have unit test. For the API test, it is optional.

There is no specific due date when the app must be completed.
The interview invitation link has been sent separately via Wellfound chat.
Later on, during the interview, you must have the ability to run the app and share the screen, we will have code review session together.

Do not hesitate to let me know if you have any questions.
Good luck!

Thank you.

## Web Service API

The web service is available on GHCR:

```sh
docker pull ghcr.io/kangaroohealthscience/fe-assessment-test:latest
docker run -it -p 3000:3000 --rm ghcr.io/kangaroohealthscience/fe-assessment-test
```

With above command, the web service shall be accessible via `http://localhost:3000`

## Endpoints

### Auth Login

User must login to get the access token, and use it to access the protected TODO list endpoints.

It uses the standard basic authentication. Use `kangaroohealth` as username and `the magnificent chicken` as the password. Example curl:

Endpoint:

```
POST /api/login
```

Example:

```bash
$ curl -X POST -u "kangaroohealth:the magnificent chicken" http://localhost:3000/api/login
{"string":"ok","data":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImJlZTIxMmM4LTcxMTQtNGZlZS1hMzk2LTM4YTkwOWY2Mjk3MSIsInVzZXJuYW1lIjoia2FuZ2Fyb29oZWFsdGgifQ.ytc23_WsRefZDZd5JniK68PgqRKok9heMHYhZHoS__k"}
```

Return value: access token

### Add TODO List

Endpoint:

```
POST /api/todo
```

Payload:

JSON object data with three fields: `name` (string), `details` (string), and `done` (boolean).

Example:

```bash
$ curl -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImJlZTIxMmM4LTcxMTQtNGZlZS1hMzk2LTM4YTkwOWY2Mjk3MSIsInVzZXJuYW1lIjoia2FuZ2Fyb29oZWFsdGgifQ.ytc23_WsRefZDZd5JniK68PgqRKok9heMHYhZHoS__k" -H "Content-type: application/json" -d "{ \"name\": \"do math homework\", \"details\": \"due date is tomorrow! do not forget\", \"done\": false }" http://localhost:3000/api/todo
```

Return value: list of TODO items.

### Get TODO List

Endpoint:

```
GET /api/todo
```

Example:

```bash
$ curl -X GET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImJlZTIxMmM4LTcxMTQtNGZlZS1hMzk2LTM4YTkwOWY2Mjk3MSIsInVzZXJuYW1lIjoia2FuZ2Fyb29oZWFsdGgifQ.ytc23_WsRefZDZd5JniK68PgqRKok9heMHYhZHoS__k" http://localhost:3000/api/todo
```

Return value: list of TODO items.

### Update TODO List

Endpoint:

```
PUT /api/todo/{todoID}
```

Payload:

JSON object data with three fields: `name` (string), `details` (string), and `done` (boolean).

Example:

```bash
$ curl -X PUT -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImJlZTIxMmM4LTcxMTQtNGZlZS1hMzk2LTM4YTkwOWY2Mjk3MSIsInVzZXJuYW1lIjoia2FuZ2Fyb29oZWFsdGgifQ.ytc23_WsRefZDZd5JniK68PgqRKok9heMHYhZHoS__k" -H "Content-type: application/json" -d "{ \"name\": \"do math homework\", \"details\": \"updated details\", \"done\": true }" http://localhost:3000/api/todo/97871bea-ea19-44d2-9778-4dc93183a1fe
```

Return value: list of TODO items.

### Delete TODO List

Endpoint:

```
DELETE /api/todo/{todoID}
```

```bash
$ curl -X DELETE -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImJlZTIxMmM4LTcxMTQtNGZlZS1hMzk2LTM4YTkwOWY2Mjk3MSIsInVzZXJuYW1lIjoia2FuZ2Fyb29oZWFsdGgifQ.ytc23_WsRefZDZd5JniK68PgqRKok9heMHYhZHoS__k" http://localhost:3000/api/todo/97871bea-ea19-44d2-9778-4dc93183a1fe
```

Return value: list of TODO items

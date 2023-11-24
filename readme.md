# GOLANG WEBSOCKET WITH SVELTEKIT

## Nugraha Agung Pratama || TalentHub batch 11

Code Repository: https://github.com/NugrahaAP/sveltekitgowebsocket

### FrontEnd

```console
// Install dep
npm install

// Start dev server default at localhost:5173
npm run dev

```

### BackEnd

```console
// Start dev server, default at localhost:1437
go run main.go
```

### feature

- Jwt Auth
- Register,Login
- Personal chat room (2 user)
- Group chat room
- Leave group

### Default user creds

```console

email: admin@admin.com
pass:  password1234

email: admin2@admin.com
pass:  password1234

email: joel@email.com
pass:  password1234

```

### API route

- GET /health_check
- POST /auth/login
- POST /auth/register
- POST /auth/refresh
- GET /backend/api/v1/rahasia
- POST /backend/api/v1/message
- DELETE /backend/api/v1/message
- PATCH /backend/api/v1/message
- POST /backend/api/v1/group_chat_room
- GET /backend/api/v1/group_chat_room?gcrid=
- GET /backend/api/v1/group_chat_room
- DELETE /backend/api/v1/group_chat_room
- PATCH /backend/api/v1/group_chat_room
- POST /backend/api/v1/chat_room
- GET /backend/api/v1/chat_room?crid=
- GET /backend/api/v1/chat_room
- DELETE /backend/api/v1/chat_room
- GET /backend/api/v1/user?userId=
- GET /backend/api/v1/user?listUser=
- GET /backend/api/v1/checkEmail?email=

### Websocket route

- /ws/:chat_room_id

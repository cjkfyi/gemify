# `gemify`

> *__not another coding assistant__?* 

## ⁉️

After being fed-up with having to Copy and Paste my code, back-and-forth a million times for a more current understanding every turn... `gemify` was suddenly born.

After noticing some traffic for this repo in the first week of prototyping, I've decided to focus a bit on documentation. This repo could have been private and monetized. I'm not going to sit here and dream of turning this into some SaaS operation. My goal is to keep this project open-source; preventing commercialization, data-collection and additional bottlenecks. 

Help a fella out, PRs are always welcomed.

## 🌱

### Dev

```sh
# first, cd somewhere awesome ~ or just ~
git clone git@github.com:cjkfyi/gemify.git
cd gemify
# from:  aistudio.google.com
echo 'API_KEY="..."' >> .env
yarn
cd extension
yarn
cd ../backend
cat scripts/gen.sh
./scripts/gen.sh
go run ./cmd/main.go
# ezpz, help me out in anyway possible :)
```

##### Tip:

ℹ️: Open or Add `extension/` into a VSCode workspace, then hit `F5`.

## 🥞

#### Backend (Go)

- Well-orchestrated H/2 Proxy & gRPC server.
- Utilizing `bitcask` & `bitio` for our storage.
- Custom logger and improved err handling. (🔜™️)

```go
require (
    go-chi/chi           v5.0.12
    google/genai         v0.10.0
    google/grpc          v1.62.1
    google/protobuf      v1.33.0 
    gorilla/websocket    v1.5.1
    joho/godotenv        v1.5.1
    prologic/bitcask     v2.0.3 
)
```

#### Frontend (JS)

- Custom, and soon to be responsive UI. 
- Webview API for integration with VSCode.

```json
  "devDependencies": {
    "esbuild": "^0.20.2",
    "eslint": "^8.56.0"
  },
  "dependencies": {
    "ws": "^8.16.0",
    "zustand": "^4.5.2"
  }
```

## 📚

### newConvo 

_r.Post("/chat", newConvo)_

`POST` @ `127.0.0.1:8000/chat`

#### Response:

```json
{
  "command": "execNewConvo",
  "status": "success",
  "data": {
    "convoID": "6e450f653022ae70"
  }
}
```

---

### getConvos 

_r.Get("/chat/list", getConvos)_

`GET` @ `127.0.0.1:8000/chat/list`

#### Response:

```json
{
  "command": "convoList",
  "status": "success",
  "data": {
    "conversations": [
      {
        "id": "6e450f653022ae70",
        "title": "VSCode extension `gemify`",
        "lastModified": "2024-04-04T23:24:24.160944195-07:00",
        "firstCreated": "2024-04-04T23:24:24.160944195-07:00",
      },
      {
        "id": "34b6a8995d1e1d39",
        "title": "Bubbletea TUI app `backstacc`",
        "lastModified": "2024-04-04T11:33:17.5810942-07:00",
        "firstCreated": "2024-04-04T11:33:17.5810942-07:00",
      },
    ]
  }
}
```

---

### newMessage

_r.Get("/ws/chat/{id}", newMessage)_

`WS` @ `ws://127.0.0.1:8000/ws/chat/{convoID}`

#### Message:

```json
{
  "message": "MD: List 3 trees, and 2 facts for each of them?"
}
```

#### Response:

```
> Connection Terminated       
23:31:26

> {"content":"EOF"}
23:31:26

> {"content":" years"}
23:31:26

> {"content":"ces maple syrup, a popular sweetener\n* Leaves turn vibrant shades of red, orange, and yellow in the fall\n\n**Tree 3: Redwood Tree**\n* Among the tallest trees in the world, reaching heights of up to 379 feet\n* Can live for over 2,000"}
23:31:26

> {"content":"1,000 years\n* Produces acorns, which are a valuable food source for wildlife\n\n**Tree 2: Maple Tree**\n* Produ"}
23:31:25

> {"content":"**Tree 1: Oak Tree**\n* Can live for up to "}
23:31:25

> { "message": "MD: List 3 trees, and 2 facts for each of them?" }
23:31:24

> Connected to ws://localhost:8000/ws/chat/34b6a8995d1e1d39
23:30:56
```


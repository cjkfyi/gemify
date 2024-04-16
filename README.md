# `gemify`

> *__not another coding assistant__?*

## 🤔

I'm prototyping a multi-turn, conversational extension for VSCode. After being fed-up with having to Copy and Paste my code, back-and-forth a million times, for a more current understanding every turn. Thus, `gemify` was suddenly born. 

## ✨

Hereby my pledge to maintain this repo. Planning to scale, especially as goog figures out the unlimited token ordeal! Simply imagine the near-future.

## 📸

:D

## 🧰

#### Backend (Go)

- Well-orchestrated H/2 Proxy & gRPC server.
- Utilizing `bitcask` & `bitio` for our storage.
- Custom logger and improved err handling. (🔜™️)
s
```go
require (
	github.com/go-chi/chi/v5               v5.0.12
	github.com/google/generative-ai-go     v0.10.0
	github.com/gorilla/websocket           v1.5.1
	github.com/icza/bitio                  v1.1.0
	github.com/joho/godotenv               v1.5.1
	go.mills.io/bitcask/v2                 v2.0.3
	golang.org/x/sync                      v0.7.0
	google.golang.org/api                  v0.172.0
	google.golang.org/grpc                 v1.63.2
	google.golang.org/protobuf             v1.33.0
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
Add `extension/` into a VSCode workspace, then hit `F5`⁉️

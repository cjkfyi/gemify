# `gemify`

> *__not another coding assistant__?* 

## ⁉️

After being fed-up with having to Copy and Paste my code, back-and-forth a million times, for a more current understanding every turn... `gemify` was suddenly born.

I've noticed some traffic for this repo in the first week of prototyping. Which lead me to start planning ahead some. Expect continuous development and information, eventually an entire community (if possible)? With all of that in-mind, came the realization that many would (or could) try to capitalize on a mere fork of this single, random idea of mine?

Help a fella out, PRs are always welcomed.

## 🥞

#### Frontend (JS)

- Custom, and soon to be responsive UI. 
- Webview API for integration with VSCode.

#### Backend (Go)

- Well-orchestrated H/2 Proxy & gRPC server.
- Utilizing `bitcask` & `bitio` for our storage.
- Custom logger and improved err handling. (🔜™️)

## 💕

### Dev

```sh
# first, cd somewhere awesome ~ or just ~
git clone git@github.com:cjkfyi/gemify.git
cd gemify
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
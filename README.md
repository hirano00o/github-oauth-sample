# github-oauth-sample
githubの認証用APIサーバのサンプル(WIP)

## frontend
認証用APIサーバを使ってログインを試すWebサーバ

#### 主な処理

- トークンを発行し、クライアントのCookieに保存する
- `backend`に対して、GithubAPIを利用する処理の依頼をする  
  処理の依頼時にはトークンを渡す

## backend
認証用APIサーバ

#### 主な処理

- `frontend`から初めて送信されたトークンに対して期限、リフレッシュトークン、リフレッシュトークンの期限を作成、保存、返却
- `frontend`から送信されたトークンに対して、Webが発行したトークンか、期限内か検証
- トークン検証後、問題なければ`frontend`から依頼された処理を実行

#### start server

```
$ cd backend
$ export DB_PASSWORD=root DB_USER=test DB_NAME=oauth DB_HOST=mysql DB_PORT=13306
$ make mysql.start
$ export GITHUB_CLIENT_ID=xxxxxxxx GITHUB_CLIENT_SECRET=yyyyyyyy SERVER_HOST=example.com
$ make start
```



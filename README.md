# github-oauth-sample
githubの認証用APIサーバのサンプル(WIP)

## 構成
frontend -> backend -> DB  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;└------> github client <- DB  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;└---------> github API

## frontend
認証用APIサーバを使ってログインを試すWebサーバ

#### 主な処理

- `backend`から受け取ったトークン等をクライアントのCookieに保存する
- `backend`に対して、GithubAPIを利用する処理の依頼をする  
  処理の依頼時にはCookie内のトークンを渡す

## backend
認証用APIサーバ

#### 主な処理

- `frontend`からログイン用URLを依頼されたら、トークンやState、トークン期限を作成し、URLとともに返却
- `frontend`から送信されたトークンが正しいか、期限内か検証
- トークン検証後、問題なければ`frontend`から依頼された処理を実行

#### start server

```
$ cd backend
$ export DB_PASSWORD=root DB_USER=test DB_NAME=oauth DB_HOST=mysql DB_PORT=13306
$ make mysql.start
$ export GITHUB_CLIENT_ID=xxxxxxxx GITHUB_CLIENT_SECRET=yyyyyyyy SERVER_HOST=example.com
$ make start
```

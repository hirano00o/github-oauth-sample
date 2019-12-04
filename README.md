# github-oauth-sample
githubの認証用APIサーバのサンプル(WIP)

## backend
認証用APIサーバ

#### start server

```
$ cd backend
$ export DB_PASSWORD=root DB_USER=test DB_NAME=oauth DB_HOST=mysql DB_PORT=13306
$ make mysql.start
$ export GITHUB_CLIENT_ID=xxxxxxxx GITHUB_CLIENT_SECRET=yyyyyyyy SERVER_HOST=example.com
$ make start
```


## frontend
認証用APIサーバを使ってログインを試すWebサーバ

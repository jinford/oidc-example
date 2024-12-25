## 手順

### Hydra を立ち上げる

```
$ cd idp
$ docker compose up -d
$ docker compose logs -f
```

マイグレーションが完了するまで待機する。

## クライアントIDの払い出し

```
$ cd idp
$ docker compose exec hydra \
    hydra create client \
    --endpoint http://127.0.0.1:4445 \
    --redirect-uri http://127.0.0.1:13000/callback
```

client_id と client_secret をメモしておく。

### IdP のフロントエンドを立ち上げる

```
$ cd idp/webui
$ npm run dev
```

ポート 14000 番でフロントの開発サーバが立ち上がる。

### IdP のバックエンドを立ち上げる

```
$ cd idp/api
$ go run .
```

ポート 14001 番でAPIサーバが立ち上がる。


### RP のフロントエンドを立ち上げる

```
$ cd rp/webui
$ npm run dev
```

ポート 13000 番でフロントの開発サーバが立ち上がる。

### RP のバックエンドを立ち上げる

メモした CLIENT_ID と CLIENT_SECRET を環境変数に設定した上で、サーバーを起動する。
CALLBACK_URL には RP のフロントエンドのURLを指定する。

```
$ export CLIENT_ID=421e65d4-7c90-415e-8ea2-7cab0541db35
$ export CLIENT_SECRET=e_GY0zABa4_Lai7hztGsm4sF~d
$ export CALLBACK_URL=http://127.0.0.1:13000/callback
$ cd rp/api
$ go run .
```

ポート 13001 番でAPIサーバが立ち上がる。


### ログインを試す

ブラウザに http://{ip}/13000/login にアクセスする。

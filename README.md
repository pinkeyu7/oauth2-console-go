# oauth2-console-go

這是一個 Oauth2 Server 的管理平台，管理 Client 以及 Scope。

## Related tables

### Oauth Client

紀錄 Client app 的資料，包含名稱、密鑰、授權的 api、允許的授權的網域、app icon。

| Field          |     Type     |      Comment       |
| -------------- | :----------: | :----------------: |
| id             | VARCHAR(255) |         id         |
| sys_account_id |   int(11)    |     manager id     |
| name           | VARCHAR(255) |        name        |
| secret         | VARCHAR(255) |   client secret    |
| domain         | VARCHAR(255) |       domain       |
| scope          | VARCHAR(255) |     allow apis     |
| icon_path      | VARCHAR(255) |   app icon path    |
| data           |     TEXT     | for oauth2 library |
| created_at     |   datetime   |                    |
| updated_at     |   datetime   |                    |

### Oauth Scope

紀錄開放 api 的資料，包含名稱、路徑、方法。

| Field       |     Type     |        Comment        |
| ----------- | :----------: | :-------------------: |
| id          |   int(11)    |          id           |
| scope       | VARCHAR(100) |  scope label(unique)  |
| path        | varchar(100) | api path(casbin rule) |
| method      | VARCHAR(20)  |      http method      |
| name        | VARCHAR(100) |     display name      |
| description | VARCHAR(255) |      description      |
| is_disable  |  tinyint(4)  |                       |
| created_at  |   datetime   |                       |
| updated_at  |   datetime   |                       |

## Oauth Scope Handling

假定目前有一個 app 想要取得 使用者資料 以及 聯絡人資料，但是並沒有新增的權限，scope 的處理方式如下。

1. 從 oauth_client 取得 client app 的資料，並從 scope 取得列表。

   | Scope                     | Authorized |
   | ------------------------- | :--------: |
   | user.profile_get          |    Yes     |
   | address-book.list_get     |    Yes     |
   | address-book.contact_post |     No     |
   | address-book.contact_get  |    Yes     |

2. 從 oauth_scope 取得所有的 scope 列表，建成樹狀表單，搭配從 client app 取得的授權資料，生成 client app 的 scope tree。

   Scope List

   - user
     - user.profile_get
   - address-book
     - address-book.list_get
     - address-book.contact_post
     - address-book.contact_get

   Scope list in JSON

   ```JSON
   {
    "address-book": {
        "name": "address-book",
        "items": {
            "contact_get": {
                "name": "contact_get",
                "is_auth": true
            },
            "contact_post": {
                "name": "contact_post",
                "is_auth": false
            },
            "list_get": {
                "name": "list_get",
                "is_auth": true
            }
        },
        "is_auth": false
    },
    "user": {
        "name": "user",
        "items": {
            "profile_get": {
                "name": "profile_get",
                "is_auth": true
            }
        },
        "is_auth": false
    }
   }
   ```

3. Scope validation flow

   ```
    Oauth2Middleware
        1. check token, get claims
        2. check scope
            a. find scope by path
            b. verify scope is authorized
                - get scope tree from redis
                - if scope tree not found from redis, build one by mysql data)
                - check scope
   ```

## GO

### 套件管理 Go Module

專案目前使用 Go Module 進行管理，Go 1.11 版本以上才有支援。

#### Go Module

先下指令 `go env` 確認 go module 環境變數是否為 `on`

如果不等於 `on` 的話，下指令

```
export GO111MODULE=on
```

即可打開 go module 的功能。

原則上專案編譯時會自行安裝相關套件，

但也可以先執行下列指令，安裝 module 套件。

```
go mod tidy
```

### How to set up environment?

安裝 docker 後，使用 `docker-compose.yml` 來建立 mysql、redis

```bash
docker-compose up -d
```

並參考 `.env.sample` 來設置 `.env`

### How to do DB migration?

我們使用 sql-migrate 套件實作 DB migration 功能，

先進行 cmd 安裝

```bash
go get -v github.com/rubenv/sql-migrate/...
```

指令如下：

執行 Migration

```bash
make migrate-up
```

Rollback Migration

```bash
make migrate-down
```

套件連結： Please refer [sql-migrate](https://github.com/rubenv/sql-migrate)

### How to develop?

```shell
go install
go run main.go
```

## Swagger API Doc

先進行 cmd 安裝

```bash
go get -u github.com/swaggo/swag/cmd/swag
```

1. 產生文件

可以經由 Makefile 執行

```bash
make doc
```

或者原生指令

```bash
swag init
```

2. And then you take [Swagger document](http://localhost:8080/swagger/index.html)

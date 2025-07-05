# Mock APIサーバー

MCP Bridge機能をテストするための設定可能なMock REST APIサーバーです。

## 概要

Mock APIサーバーは、設定可能なエンドポイント、認証、データを備えた柔軟なテスト環境を提供します。異なるAPIシナリオ用の複数の設定ファイルをサポートしています。

## クイックスタート

```bash
# Mock APIサーバーのビルド
go build -o bin/mock-api ./cmd/mock-api

# デフォルト設定で起動（ユーザーAPI）
./bin/mock-api

# 商品設定で起動
MOCK_CONFIG=configs/mock/products.json ./bin/mock-api

# 特定のポートで起動
PORT=8081 ./bin/mock-api

# または直接実行
go run ./cmd/mock-api
```

## 設定

### 利用可能な設定

Mock APIサーバーは `configs/mock/` ディレクトリの設定ファイルを使用します：

- `configs/mock/users.json` - ユーザーAPI（デフォルト、ポート8080）
- `configs/mock/products.json` - 商品API（ポート8081）

### 設定構造

```json
{
  "server": {
    "port": "8080",
    "name": "Mock API Server"
  },
  "auth": {
    "enabled": false,
    "username": "admin",
    "password": "password"
  },
  "resources": [
    {
      "name": "users",
      "path": "/users",
      "enabled": true,
      "data": [...ユーザーオブジェクト...],
      "methods": ["GET", "POST", "PUT", "DELETE"],
      "supportsId": true
    }
  ],
  "endpoints": [
    {
      "path": "/health",
      "method": "GET",
      "enabled": true,
      "response": {"status": "healthy"}
    }
  ]
}
```

### 環境変数

- `MOCK_CONFIG` - 設定ファイルのパス（デフォルト: `configs/mock/users.json`）
- `PORT` - サーバーポート（設定ファイルを上書き）
- `AUTH_ENABLED` - Basic認証を有効化（`true`/`false`）
- `AUTH_USERNAME` - 認証用ユーザー名（デフォルト: `admin`）
- `AUTH_PASSWORD` - 認証用パスワード（デフォルト: `password`）

## Basic認証

Mock APIサーバーは環境変数によるBasic認証をサポートしています：

```bash
# Basic認証を有効にして起動
AUTH_ENABLED=true AUTH_USERNAME=admin AUTH_PASSWORD=secret PORT=8081 go run ./cmd/mock-api

# カスタム認証情報で起動
AUTH_ENABLED=true AUTH_USERNAME=myuser AUTH_PASSWORD=mypass PORT=8081 go run ./cmd/mock-api

# 認証付きエンドポイントのテスト
curl -u admin:secret http://localhost:8081/users

# またはAuthorizationヘッダーで認証
curl -H "Authorization: Basic YWRtaW46c2VjcmV0" http://localhost:8081/users
```

Basic認証が有効な場合、すべてのエンドポイントで有効な認証情報が必要です。認証なしでアクセスすると `401 Unauthorized` が返されます。

## 利用可能なエンドポイント

### ユーザー設定使用時（デフォルト）

- `GET /health` - ヘルスチェック
- `GET /users` - 全ユーザー取得  
- `POST /users` - ユーザー作成
- `GET /users/{id}` - 特定ユーザー取得
- `PUT /users/{id}` - ユーザー更新
- `DELETE /users/{id}` - ユーザー削除

### 商品設定使用時

- `GET /health` - ヘルスチェック
- `GET /products` - 全商品取得
- `GET /products/{id}` - 特定商品取得

## 使用例

### 異なるサービスの起動

```bash
# ユーザーAPI（デフォルト）
./bin/mock-api

# 商品API
MOCK_CONFIG=configs/mock/products.json ./bin/mock-api

# 認証付きユーザーAPI
AUTH_ENABLED=true AUTH_USERNAME=admin AUTH_PASSWORD=secret ./bin/mock-api

# カスタムポートで商品API
MOCK_CONFIG=configs/mock/products.json PORT=9000 ./bin/mock-api
```

### エンドポイントのテスト

```bash
# ヘルスチェック
curl http://localhost:8080/health

# 全ユーザー取得
curl http://localhost:8080/users

# 新しいユーザー作成
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "田中太郎", "email": "tanaka@example.com"}'

# 特定ユーザー取得
curl http://localhost:8080/users/1

# ユーザー更新
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "田中二郎", "email": "tanaka2@example.com"}'

# ユーザー削除
curl -X DELETE http://localhost:8080/users/1

# 認証付きテスト
curl -u admin:secret http://localhost:8081/users
```

## カスタム設定の作成

`configs/mock/` ディレクトリにカスタム設定ファイルを作成できます：

```json
{
  "server": {
    "port": "8082",
    "name": "カスタムMock API"
  },
  "auth": {
    "enabled": true,
    "username": "custom",
    "password": "secret"
  },
  "resources": [
    {
      "name": "orders",
      "path": "/orders",
      "enabled": true,
      "data": [
        {
          "id": 1,
          "customerId": 1,
          "total": 99.99,
          "status": "pending"
        }
      ],
      "methods": ["GET", "POST"],
      "supportsId": true
    }
  ],
  "endpoints": [
    {
      "path": "/status",
      "method": "GET",
      "enabled": true,
      "response": {
        "service": "orders",
        "status": "running"
      }
    }
  ]
}
```

起動方法：
```bash
MOCK_CONFIG=configs/mock/custom.json ./bin/mock-api
```

## CORSサポート

Mock APIサーバーは自動的にクロスオリジンリクエスト用のCORSヘッダーを含みます：

- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization`

## 機能

- **動的設定**: 環境変数による異なるAPI設定の読み込み
- **リソース管理**: 設定されたリソースのCRUD操作
- **カスタムエンドポイント**: 静的レスポンスエンドポイントの定義
- **認証**: オプションのBasic認証
- **CORSサポート**: 組み込みクロスオリジンリクエストサポート
- **タイムスタンプテンプレート**: レスポンス内で `{{timestamp}}` を使用した動的タイムスタンプ
- **柔軟なデータ型**: 様々なJSONデータ構造のサポート
# MCP Bridge

REST APIをMCPサーバーとして利用するためのプロキシサーバーです。

## 概要

MCP Bridge は、既存の REST API を Model Context Protocol (MCP) サーバーとして利用できるようにするプロキシサーバーです。これにより、REST API を MCP クライアント（Claude Code など）から直接利用できるようになります。

## プロジェクトの意義

### 技術的価値
- **プロトコル統一**: 既存のREST APIをMCPプロトコルに変換し、AIツールとの統合を標準化
- **アダプター層**: 異なるAPI形式を統一インターフェースで扱えるブリッジ機能を提供
- **型安全性**: JSON-RPC 2.0準拠でスキーマベースの型チェック

### 実用的価値
- **既存資産活用**: 新しいAPIを作らずに、既存REST APIをMCPクライアント（Claude Code等）で直接利用可能
- **開発効率向上**: 各APIの個別実装不要で、設定ファイルだけで新しいAPIを追加
- **認証・セキュリティ**: 統一されたヘッダー管理とエラーハンドリング

### エコシステム貢献
- **MCP普及**: REST APIをMCP対応にすることで、MCPエコシステムの拡張に貢献
- **相互運用性**: 異なるサービス間の連携を促進する標準的な方法を提供

## 特徴

- **REST API to MCP変換**: REST APIエンドポイントをMCPツールとして自動変換
- **JSON-RPC 2.0準拠**: MCPプロトコルに完全準拠
- **複数トランスポート対応**: 標準入出力とHTTP通信の両方をサポート
- **設定可能**: 設定ファイルによる柔軟なカスタマイズ
- **Mock APIサーバー**: テスト用のシンプルなREST APIサーバーを内蔵

## プロジェクト構成

```
mcp-bridge/
├── cmd/
│   ├── mcp-server-stdio/  # 標準入出力トランスポート対応MCPサーバー
│   ├── mcp-server-http/   # HTTP通信トランスポート対応MCPサーバー
│   └── mock-api/          # 設定可能なMock APIサーバー
├── internal/
│   ├── mcp/              # MCP実装
│   ├── bridge/           # REST API変換ロジック
│   ├── transport/        # トランスポート層（stdio/HTTP）
│   └── config/           # 設定管理
├── pkg/
│   └── types/            # 共通型定義
├── go.mod
├── go.sum
├── README.md
└── README-ja.md          # Japanese version
```

## インストール・ビルド

### 依存関係
- Go 1.21以上

### ビルド

```bash
# 標準入出力トランスポート対応MCPサーバーのビルド
go build -o bin/mcp-server-stdio ./cmd/mcp-server-stdio

# HTTP通信トランスポート対応MCPサーバーのビルド
go build -o bin/mcp-server-http ./cmd/mcp-server-http

# Mock APIサーバーのビルド
go build -o bin/mock-api ./cmd/mock-api
```

## 使用方法

### 1. Mock APIサーバーの起動

設定可能なMock REST APIサーバーを起動します：

```bash
# デフォルト設定で起動（ユーザーAPI）
./bin/mock-api

# 商品設定で起動
MOCK_CONFIG=configs/mock/products.json ./bin/mock-api

# または直接実行
go run ./cmd/mock-api
```

Mock APIの詳細なドキュメント、設定オプション、使用例については、**[Mock APIドキュメント](docs/MOCK-API.ja.md)** をご覧ください。

### 2. MCPサーバーの起動

MCPブリッジサーバーを起動します：

#### 標準入出力トランスポート（従来版）

```bash
./bin/mcp-server-stdio

# または設定ファイルを指定
./bin/mcp-server-stdio -config ./config.json

# またはAPIベースURLを直接指定
./bin/mcp-server-stdio -api-url http://localhost:8080

# 詳細ログを有効にする場合
./bin/mcp-server-stdio -verbose
```

#### HTTP通信トランスポート（新機能）

```bash
./bin/mcp-server-http

# HTTP設定を指定
./bin/mcp-server-http -port 8080 -host localhost -cors

# 設定ファイルとHTTP設定を指定
./bin/mcp-server-http -config ./example-config.json -port 8080

# 詳細ログを有効にする場合
./bin/mcp-server-http -verbose
```

HTTPサーバーは `http://localhost:8080/mcp` でMCP JSON-RPCリクエストを受け付けます。

### 3. 設定ファイル

設定ファイル例（`example-config.json`）：

```json
{
  "apis": [
    {
      "name": "users-api",
      "baseUrl": "http://localhost:8081",
      "timeout": 30,
      "headers": {
        "X-API-Key": "your-api-key-here"
      },
      "auth": {
        "type": "basic",
        "basic": {
          "username": "admin",
          "password": "password"
        }
      },
      "endpoints": [
        {
          "name": "health",
          "description": "ヘルスチェックエンドポイント",
          "method": "GET",
          "path": "/health",
          "parameters": []
        },
        {
          "name": "get_users",
          "description": "全ユーザー取得",
          "method": "GET",
          "path": "/users",
          "parameters": []
        },
        {
          "name": "create_user",
          "description": "新しいユーザー作成",
          "method": "POST",
          "path": "/users",
          "parameters": [
            {
              "name": "email",
              "type": "string",
              "required": true,
              "description": "ユーザーメールアドレス",
              "in": "body"
            },
            {
              "name": "name",
              "type": "string",
              "required": true,
              "description": "ユーザー名",
              "in": "body"
            }
          ]
        },
        {
          "name": "get_user",
          "description": "特定ユーザーをIDで取得",
          "method": "GET",
          "path": "/users/{id}",
          "parameters": [
            {
              "name": "id",
              "type": "integer",
              "required": true,
              "description": "ユーザーID",
              "in": "path"
            }
          ]
        },
        {
          "name": "update_user",
          "description": "特定ユーザーをIDで更新",
          "method": "PUT",
          "path": "/users/{id}",
          "parameters": [
            {
              "name": "id",
              "type": "integer",
              "required": true,
              "description": "ユーザーID",
              "in": "path"
            },
            {
              "name": "email",
              "type": "string",
              "required": true,
              "description": "ユーザーメールアドレス",
              "in": "body"
            },
            {
              "name": "name",
              "type": "string",
              "required": true,
              "description": "ユーザー名",
              "in": "body"
            }
          ]
        },
        {
          "name": "delete_user",
          "description": "特定ユーザーをIDで削除",
          "method": "DELETE",
          "path": "/users/{id}",
          "parameters": [
            {
              "name": "id",
              "type": "integer",
              "required": true,
              "description": "ユーザーID",
              "in": "path"
            }
          ]
        }
      ]
    },
    {
      "name": "products-api",
      "baseUrl": "http://localhost:8082",
      "timeout": 30,
      "headers": {
        "Authorization": "Bearer your-token-here"
      },
      "endpoints": [
        {
          "name": "get_products",
          "description": "全商品取得",
          "method": "GET",
          "path": "/products",
          "parameters": []
        },
        {
          "name": "get_product",
          "description": "特定商品をIDで取得",
          "method": "GET",
          "path": "/products/{id}",
          "parameters": [
            {
              "name": "id",
              "type": "integer",
              "required": true,
              "description": "商品ID",
              "in": "path"
            }
          ]
        }
      ]
    }
  ],
  "server": {
    "name": "mcp-bridge",
    "version": "1.0.0",
    "description": "REST API to MCP Bridge Server"
  },
  "headers": {
    "Content-Type": "application/json"
  }
}
```

**注意**: 設定ファイルは標準入出力版とHTTP版で共通です。トランスポートの種類は使用するサーバー実行ファイル（`mcp-server-stdio` または `mcp-server-http`）によって決まります。

HTTP通信の場合は、コマンドライン引数でサーバー設定を指定できます：

```bash
# カスタム設定でHTTP通信を使用
go run ./cmd/mcp-server-http --config ./example-config.json --port 8080 --host localhost --cors
```

### 4. Claude Codeでの利用

Claude Codeで使用する場合の設定例：

#### 標準入出力トランスポート
```json
{
  "mcpServers": {
    "mcp-bridge": {
      "command": "go",
      "args": ["run", "./cmd/mcp-server-stdio", "--config", "./example-config.json"]
    }
  }
}
```

#### HTTP通信トランスポート
```json
{
  "mcpServers": {
    "mcp-bridge-http": {
      "transport": {
        "type": "http",
        "url": "http://localhost:8080/mcp"
      }
    }
  }
}
```

注意：HTTP MCPサーバーは別途起動する必要があります：
```bash
go run ./cmd/mcp-server-http --config ./example-config.json --port 8080
```

## 利用可能なツール

MCPブリッジサーバーが提供するツールの一覧：

### デフォルトツール（設定例を使用した場合）

- `health` - ヘルスチェックエンドポイント
- `get_users` - 全ユーザー取得
- `create_user` - ユーザー作成
- `get_user` - 特定ユーザー取得
- `update_user` - ユーザー更新
- `delete_user` - ユーザー削除
- `get_products` - 全商品取得
- `get_product` - 特定商品取得

### 利用例

```javascript
// ヘルスチェック
await callTool("health", {});

// ユーザー一覧取得
await callTool("get_users", {});

// 新しいユーザー作成
await callTool("create_user", {
  name: "田中太郎",
  email: "tanaka@example.com"
});

// 特定ユーザー取得
await callTool("get_user", {
  id: 1
});

// 全商品取得
await callTool("get_products", {});

// 特定商品取得
await callTool("get_product", {
  id: 1
});
```

## リソース

MCPサーバーは以下のリソースを提供します：

- `rest-api://docs` - REST APIの仕様書（JSON形式）

## 開発

### テスト実行

```bash
# Mock APIサーバーの起動
go run ./cmd/mock-api &

# MCPサーバーのテスト（標準入出力）
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0.0"}}}' | go run ./cmd/mcp-server-stdio

# MCPサーバーのテスト（HTTP通信）
go run ./cmd/mcp-server-http -port 8080 &
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "id": 1, "method": "ping"}'
```

### カスタムエンドポイントの追加

設定ファイルの `endpoints` セクションに新しいエンドポイントを追加することで、カスタムAPIエンドポイントを利用できます。

## ライセンス

MIT License

## 貢献

プルリクエストやイシューの報告を歓迎します。
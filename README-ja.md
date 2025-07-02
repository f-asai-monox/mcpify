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
- **設定可能**: 設定ファイルによる柔軟なカスタマイズ
- **Mock APIサーバー**: テスト用のシンプルなREST APIサーバーを内蔵

## プロジェクト構成

```
mcp-bridge/
├── cmd/
│   ├── mcp-server/     # MCPサーバー実行ファイル
│   └── mock-api/       # REST API mockサーバー
├── internal/
│   ├── mcp/           # MCP実装
│   ├── bridge/        # REST API変換ロジック
│   └── config/        # 設定管理
├── pkg/
│   └── types/         # 共通型定義
├── go.mod
├── go.sum
└── README.md
```

## インストール・ビルド

### 依存関係
- Go 1.21以上

### ビルド

```bash
# MCPサーバーのビルド
go build -o bin/mcp-server ./cmd/mcp-server

# Mock APIサーバーのビルド
go build -o bin/mock-api ./cmd/mock-api
```

## 使用方法

### 1. Mock APIサーバーの起動

テスト用のREST APIサーバーを起動します：

```bash
./bin/mock-api

# または直接実行
go run ./cmd/mock-api
```

APIサーバーは `http://localhost:8080` で起動し、以下のエンドポイントが利用できます：

- `GET /health` - ヘルスチェック
- `GET /users` - 全ユーザー取得
- `POST /users` - ユーザー作成
- `GET /users/{id}` - 特定ユーザー取得
- `PUT /users/{id}` - ユーザー更新
- `DELETE /users/{id}` - ユーザー削除

### 2. MCPサーバーの起動

MCPブリッジサーバーを起動します：

```bash
./bin/mcp-server

# または設定ファイルを指定
./bin/mcp-server -config ./config.json

# またはAPIベースURLを直接指定
./bin/mcp-server -api-url http://localhost:8080

# 詳細ログを有効にする場合
./bin/mcp-server -verbose
```

### 3. 設定ファイル

設定ファイル例（`config.json`）：

```json
{
  "api": {
    "baseUrl": "http://localhost:8080",
    "timeout": 30
  },
  "server": {
    "name": "mcp-bridge",
    "version": "1.0.0",
    "description": "REST API to MCP Bridge Server"
  },
  "headers": {
    "Content-Type": "application/json",
    "Authorization": "Bearer your-token-here"
  },
  "endpoints": [
    {
      "name": "custom_endpoint",
      "description": "カスタムエンドポイント",
      "method": "GET",
      "path": "/api/custom",
      "parameters": [
        {
          "name": "param1",
          "type": "string",
          "required": true,
          "description": "パラメータ1",
          "in": "query"
        }
      ]
    }
  ]
}
```

### 4. Claude Codeでの利用

Claude Codeで使用する場合の設定例：

```json
{
  "mcpServers": {
    "rest-api-bridge": {
      "command": "/path/to/mcp-server",
      "args": ["-api-url", "http://localhost:8080"]
    }
  }
}
```

## 利用可能なツール

MCPブリッジサーバーが提供するツールの一覧：

### デフォルトツール（Mock APIサーバー使用時）

- `get_users` - 全ユーザー取得
- `create_user` - ユーザー作成
- `get_user` - 特定ユーザー取得
- `update_user` - ユーザー更新
- `delete_user` - ユーザー削除
- `get_products` - 商品取得
- `create_product` - 商品作成
- `health_check` - ヘルスチェック

### 利用例

```javascript
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

// 商品をカテゴリで絞り込み
await callTool("get_products", {
  category: "Electronics"
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

# MCPサーバーのテスト
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0.0"}}}' | go run ./cmd/mcp-server
```

### カスタムエンドポイントの追加

設定ファイルの `endpoints` セクションに新しいエンドポイントを追加することで、カスタムAPIエンドポイントを利用できます。

## ライセンス

MIT License

## 貢献

プルリクエストやイシューの報告を歓迎します。
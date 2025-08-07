# mcpify

REST APIをMCPサーバーとして利用するためのプロキシサーバーです。

## 特徴

- **REST API to MCP変換**: REST APIエンドポイントをMCPツールとして自動変換
- **複数トランスポート対応**: 標準入出力とHTTP通信の両方をサポート
- **JSON-RPC 2.0準拠**: MCPプロトコルに完全準拠
- **設定可能**: 設定ファイルによる柔軟なカスタマイズ
- **Mock APIサーバー**: テスト用のシンプルなREST APIサーバーを内蔵

## インストール

### オプション1: ビルド済みバイナリをダウンロード（推奨）

[Releasesページ](https://github.com/f-asai-monox/mcpify/releases)から、お使いのプラットフォーム用の最新リリースをダウンロードしてください。

#### Linux/macOS

curlを使用:
```bash
curl -sSL https://raw.githubusercontent.com/f-asai-monox/mcpify/main/install.sh | bash
```

wgetを使用:
```bash
wget -qO- https://raw.githubusercontent.com/f-asai-monox/mcpify/main/install.sh | bash
```

#### Windows

PowerShellを使用（推奨）:
```powershell
iwr -useb https://raw.githubusercontent.com/f-asai-monox/mcpify/main/install.ps1 | iex
```

コマンドプロンプトを使用:
```cmd
curl -sSL https://raw.githubusercontent.com/f-asai-monox/mcpify/main/install.bat -o install.bat && install.bat
```

#### 手動ダウンロード
1. [Releases](https://github.com/f-asai-monox/mcpify/releases)にアクセス
2. お使いのOSとアーキテクチャに合ったアーカイブをダウンロード
3. 展開してバイナリをPATHに移動

### オプション2: go installでインストール
```bash
go install github.com/f-asai-monox/mcpify/cmd/mcp-server-stdio@latest
```

### オプション3: ソースからビルド
```bash
# リポジトリをクローン
git clone https://github.com/f-asai-monox/mcpify.git
cd mcpify

# MCPサーバーのビルド
go build -o bin/mcp-server-stdio ./cmd/mcp-server-stdio
go build -o bin/mcp-server-http ./cmd/mcp-server-http

# テスト用Mock APIのビルド
go build -o bin/mock-api ./cmd/mock-api
```

## クイックスタート

### 1. Mock APIの起動（テスト用）
```bash
# ソースからビルドした場合
./bin/mock-api
```

### 2. MCPサーバーの起動

#### 標準入出力方式（Claude Code用）:
```bash
# バイナリでインストールした場合
mcp-server-stdio

# 設定ファイルを指定
mcp-server-stdio -config ./example-config.json

# API URLを指定
mcp-server-stdio -api-url http://localhost:8080

# ソースからビルドした場合
./bin/mcp-server-stdio -config ./example-config.json
```

#### HTTP方式:
```bash
# バイナリでインストールした場合
mcp-server-http -port 8080

# 設定ファイルを指定
mcp-server-http -config ./example-config.json -port 8080

# ソースからビルドした場合
./bin/mcp-server-http -port 8080
```

## 基本的な使用方法

### 設定例
`config.json`ファイルを作成：

```json
{
  "apis": [
    {
      "name": "users-api",
      "baseUrl": "http://localhost:8081",
      "endpoints": [
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
              "name": "name",
              "type": "string",
              "required": true,
              "description": "ユーザー名",
              "in": "body"
            },
            {
              "name": "email",
              "type": "string",
              "required": true,
              "description": "メールアドレス",
              "in": "body"
            }
          ]
        }
      ]
    }
  ]
}
```

### Claude Codeでの使用
```json
{
  "mcpServers": {
    "mcp-bridge": {
      "command": "mcp-server-stdio",
      "args": ["-config", "./config.json"]
    }
  }
}
```

## 利用可能なツール

設定例では以下のツールが利用可能：
- `get_users` - 全ユーザー取得
- `create_user` - 新しいユーザー作成
- `get_user` - 特定ユーザー取得
- `update_user` - ユーザー情報更新
- `delete_user` - ユーザー削除

## HTTP通信

標準入出力の代わりにHTTP通信を使用する場合：

```bash
# HTTPサーバーの起動
go build -o bin/mcp-server-http ./cmd/mcp-server-http
./bin/mcp-server-http -port 8080

# Claude Codeの設定
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

## ドキュメント

- **[アーキテクチャ](docs/ARCHITECTURE.md)** - プロジェクト構造と技術詳細
- **[設定ガイド](docs/CONFIGURATION.md)** - 完全な設定ガイド
- **[開発ガイド](docs/DEVELOPMENT.md)** - 開発・テストガイド
- **[APIリファレンス](docs/API-REFERENCE.md)** - 利用可能ツールと使用例

## ライセンス

MIT License

## 貢献

プルリクエストやイシューの報告を歓迎します。
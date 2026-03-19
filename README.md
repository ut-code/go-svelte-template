# web-go-template

Go (Chi + sqlc + PostgreSQL) + SvelteKit のWebアプリテンプレート。

Wordle風ミニゲームをサンプルとして含む。

## 構成

- `backend/` — Go API サーバー (Chi router, sqlc, pgx)
- `frontend/` — SvelteKit (adapter-node, SSR, Tailwind CSS, daisyUI)

## 開発

### 環境構築

以下のいずれかの方法でツールをインストール。

**Nix** — すべてのツールが `flake.nix` で提供される。

```bash
direnv allow
# もしくは: nix develop
```

**mise** — `.mise.toml` の `[tools]` セクションをアンコメントしてからインストール。

```bash
mise install

# sqlc, golangci-lint, air は別途インストール
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/air-verse/air@latest
# https://golangci-lint.run/welcome/install/
```

**手動** — 以下を個別にインストール。

- [Docker](https://docs.docker.com/get-docker/) (PostgreSQL コンテナ用)
- [Go](https://go.dev/dl/) (1.25+)
- [Bun](https://bun.sh/) (1.3+)
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)
- [golangci-lint](https://golangci-lint.run/welcome/install/)

### 起動

```bash
# 初回セットアップ: DB 起動 + スキーマ適用 + シードデータ投入
mise run db-setup

# 開発サーバー起動 (DB + Backend + Frontend)
mise run dev
```

- Frontend: http://localhost:5173
- Backend: http://localhost:8080

### mise タスク

`.mise.toml` にタスクが定義されている。`mise run <タスク名>` で実行。

```bash
mise run lint       # golangci-lint + Biome
mise run format     # goimports/gofmt + Biome
```

### 環境変数

| 変数 | 説明 | デフォルト |
|------|------|-----------|
| `DATABASE_URL` | PostgreSQL 接続文字列 | `postgres://postgres:postgres@localhost:5432/wordle?sslmode=disable` |
| `PORT` | Backend リッスンポート | `8080` |
| `API_URL` | Backend URL (SSR サーバーから) | `http://localhost:8080` |

### Docker Compose (動作確認用)

[Docker](https://docs.docker.com/get-docker/) のみで動く。ホットリロードは対応していない。

```bash
docker compose up
```

- Frontend: http://localhost:3000
- Backend: http://localhost:8080

## デプロイ (Coolify)

### Docker Compose

`docker-compose.yaml` にはポートマッピングを含めていない（Coolify/Traefik が管理するため）。ローカルでは `docker-compose.override.yaml` が自動で読み込まれ、ポートが公開される。

1. Coolifyで新規プロジェクト作成 → **Docker Compose** を選択
2. GitHubリポジトリ接続
3. Compose ファイルに `docker-compose.yaml` を指定
4. 各サービスのポートとドメインを Coolify の UI で設定
5. 環境変数で `API_URL` を設定（frontend → backend の内部通信用）
6. デプロイ

### 個別サービス

Backend, Frontend, DB をそれぞれ別サービスとして追加。個別にスケール・再デプロイ可能。

**Backend**

1. GitHubリポジトリ接続、Base Directory を `backend` に設定
2. ビルドパック: **Dockerfile**
3. ポート: `8080`
4. 環境変数に `DATABASE_URL` を設定

**Frontend**

1. GitHubリポジトリ接続、Base Directory を `frontend` に設定
2. ビルドパック: **Dockerfile**
3. ポート: `3000`
4. 環境変数に `API_URL` を設定（バックエンドの内部URL）

**DB**

CoolifyのUIから PostgreSQL をワンクリックで追加可能。

## このテンプレートの再現手順

このリポジトリを作成した際のコマンド。テンプレートとして使う場合はこのセクションは不要。

```bash
# backend
mkdir -p backend && cd backend
go mod init web-go-template
go get github.com/go-chi/chi/v5 github.com/jackc/pgx/v5
mkdir -p sql
# sql/schema.sql, sql/queries.sql, sqlc.yaml を作成後:
sqlc generate

# frontend
cd ..
bunx sv create frontend --template minimal --types ts
cd frontend
bun add -d @sveltejs/adapter-node @biomejs/biome tailwindcss @tailwindcss/vite daisyui@beta
# svelte.config.js で adapter-node に変更
# vite.config.ts に tailwindcss plugin を追加
```

## 参考

- [Chi](https://github.com/go-chi/chi)
- [sqlc](https://sqlc.dev/)
- [pgx](https://github.com/jackc/pgx)
- [Air](https://github.com/air-verse/air)
- [SvelteKit](https://svelte.dev/docs/kit)
- [Tailwind CSS](https://tailwindcss.com/)
- [daisyUI](https://daisyui.com/)
- [Biome](https://biomejs.dev/)
- [Coolify](https://coolify.io/docs/)

<div align="center">
    <a href="https://rk2.uk">
        <img src="./app/static/images/logo-horizontal.png" />
    </a>
    <a href="./LICENSE">
        <img src="https://img.shields.io/badge/license-MIT-blue.svg?style=flat">
    </a>
    <img src="https://github.com/Outtech105k/ShortUrlServer/actions/workflows/test.yml/badge.svg" />
    <a href="https://codecov.io/gh/Outtech105k/ShortUrlServer">
        <img src="https://codecov.io/gh/Outtech105k/ShortUrlServer/main/graph/badge.svg" />
    </a>
</div>

## Overview

自分の使いたい形式に合わせて、カスタムURLを作成できます。

ここから利用できます。
[https://rk2.uk](https://rk2.uk)

個々の要望に応じたURLを職人のように生成したい、という想いで **"URL Forge"** と命名しました。

## Rules & Constraints

URL生成に関するルールや制約については、以下をご確認ください。

[URL生成ルールと制約](/docs/rules.md)

## REST API Usage

URL生成サービスは、REST APIに対応しています。
GUIアプリと機能は同じです。

[POST `/api/set`](/docs/api/POSTset.md): Set Custom URL

## Usage

### Preconfigure

`config.sample.env` を `config.env` にコピーし、環境に合わせて各値を設定してください。

#### 設定項目

| 環境変数 | 説明 | デフォルト値 |
| :--- | :--- | :--- |
| `ENDPOINT` | **(必須)** 公開サーバーのURL (例: `https://example.com`) | - |
| `PORT` | サーバーの待受ポート番号 | `8080` |
| `REDIS_ADDR` | Redisサーバーのアドレス | `redis:6379` |
| `REDIS_PASSWORD` | Redisのパスワード | (空) |
| `REDIS_DB` | RedisのDB番号 | `0` |
| `APP_NAME` | アプリケーション名 (HTMLタイトル等に反映) | `URL Forge` |
| `ALLOW_ORIGINS` | CORS許可オリジン (カンマ区切り、または `*`) | `*` |
| `SHUTDOWN_TIMEOUT` | シャットダウン時の待機時間 | `5s` |
| `OGP_FETCH_TIMEOUT` | OGP取得時のタイムアウト時間 | `5s` |
| `DEFAULT_ID_LENGTH` | 自動生成されるIDのデフォルト長 | `6` |
| `MAX_ID_LENGTH` | IDの最大長 | `100` |
| `MAX_RETRY_COUNT` | ID自動生成時の衝突回避リトライ回数 | `10` |
| `BOT_USER_AGENTS` | ボット判定に使用するUser-Agentのリスト (カンマ区切り) | (主要なボット) |

### Startup

1. 開発環境では

```bash
docker compose -f compose.dev.yml up --build
```

[Air](https://github.com/air-verse/air) を利用してホットリロード開発ができます。(適宜`-d`オプションを付加してください)

2. デプロイ環境では

```bash
docker compose -f compose.prod.yml up -d --build
```

マルチステージングにより、バイナリにビルドした後に [Alpineコンテナ](https://hub.docker.com/_/alpine)で実行されます。

## Thanks

ロゴ作成ツールには [Shopify ロゴメーカー](https://www.shopify.com/jp/tools/logo-maker) を使用しました。
素晴らしいサービスの提供者に感謝申し上げます。

## Contact

Plat (プラット)

<a href="https://github.com/Outtech105k">
    <img src="https://img.shields.io/badge/-@Outtech105k-000000.svg?logo=github&style=flat">
</a>
<a href="https://x.com/105techno">
    <img src="https://img.shields.io/badge/-@105techno-000000.svg?logo=x&style=flat">
</a>
<a href="mailto:techno510tk@gmail.com">
    <img src="https://img.shields.io/badge/-techno510tk@gmail.com-000000.svg?logo=gmail&style=flat">
</a>

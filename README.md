# ShortUrlServer

<a href="./LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg?style=flat">
</a>

![build workflow](https://github.com/Outtech105k/ShortUrlServer/actions/workflows/test.yml/badge.svg)
[![codecov](https://codecov.io/gh/Outtech105k/ShortUrlServer/main/graph/badge.svg)](https://codecov.io/gh/Outtech105k/ShortUrlServer)


## Overview

**URLをカスタムで設定できるサービス**です。

ここで利用できます。
[https://rk2.uk](https://rk2.uk)

## REST API Usage

URL生成サービスは、REST APIに対応しています。
GUIアプリと機能は同じです。

[POST `/api/set`](/docs/api/POSTset): Set Custom URL

## Usage

### Preconfigure

[/config.sample.env](/config.sample.env) 内の`ENDPOINT` を、運用するサーバのエンドポイントに合わせて設定し、`config.env`にリネームします。

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

## Contact

<a href="https://github.com/Outtech105k">
    <img src="https://img.shields.io/badge/-@Outtech105k-000000.svg?logo=github&style=flat">
</a>
<a href="https://x.com/105techno">
    <img src="https://img.shields.io/badge/-@105techno-000000.svg?logo=x&style=flat">
</a>
<a href="mailto:techno510tk@gmail.com">
    <img src="https://img.shields.io/badge/-techno510tk@gmail.com-000000.svg?logo=gmail&style=flat">
</a>

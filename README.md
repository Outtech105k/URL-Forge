# ShortUrlServer

![build workflow](https://github.com/Outtech105k/ShortUrlServer/actions/workflows/test.yml/badge.svg)

[![codecov](https://codecov.io/gh/Outtech105k/ShortUrlServer/main/graph/badge.svg)](https://codecov.io/gh/Outtech105k/ShortUrlServer)

## Overview

**URLをカスタムで設定できるサービス**です。

ここで利用できます。
[https://rk2.uk](https://rk2.uk)

## REST API Usage

URL生成サービスは、REST APIに対応しています。
GUIアプリと機能は同じです。

### Requests

`"expire_in"`キーは、正規表現`` `[0-9]+[smhd]` ``を受け付けます。
`"10m", "30h", "2d"`のように指定してください。
指定しない場合、有効期限は設定されません。

`"base_url"`のみの指定の場合、ランダムIDがセットされます。

| key | 説明 | 必須/デフォルト値 | 競合する値 |
| :-- | :-- | :-- | :-- |
| `base_url` | リダイレクト先URL | 必須 | |
| `use_uppercase`| ランダムIDに英大文字を含めるか | `false` | `custom_id` |
| `use_lowercase`| ランダムIDに英小文字を含めるか | `true` | `custom_id` |
| `use_numbers`| ランダムIDに数字を含めるか | `true` | `custom_id` |
| `id_length`| ランダムIDの文字数 | `6` | `custom_id` |
| `custom_id`| 設定するカスタムID<br>(最大文字数100文字) | ランダムIDを採用 | `use_uppercase`, `use_lowercase`, `use_numbers`, `id_length` |
| `expire_in`| リンクの有効期間 | 無期限 | |
| `sand_cushion`| クッションページを使用するか | `false` | |

1. ランダムIDでURL生成する例

```JSON
{
    "base_url": "https://example.com",
    "use_uppercase": true,
    "use_lowercase": false,
    "use_numbers": true,
    "id_length": 5,
    "expire_in": "10h",
    "sand_cushion": true
}
```

2. カスタムIDでURL生成する例

```JSON
{
    "base_url": "https://example.com",
    "custom_id": "example",
    "expire_in": "10h",
    "sand_cushion": true
}
```

### Responces

1. 200 OK Responce

```JSON
{
    "base_url": "https://example.com",
    "short_url": "https://rk2.uk/example"
}
```

2. 400 Bad Request Responce
    - Varidation Error
    ```JSON
    {
        "type": "varidation_error",
        "details": [
            {
                "field": "base_url",
                "message": "base_url is required"
            }
        ]
    }
    ```
    必要なパラメータがない、もしくは入力制約に違反した場合に返されます。

    - Invalid Request
    ```JSON
    {
        "type": "invalid_request",
        "message": "Empty JSON body"
    }
    ```
    JSON Bodyに対する問題や、入力制約に違反した場合に返されます。

    - Parameter Conflict 
    ```JSON
    {
        "type": "parameter_conflict",
        "message": "*** cannot be used together with ***"
    }
    ```
    競合関係にある（同時に設定できない）JSONパラメータを設定した場合に返されます。

3. 409 Conflict Responce
    ```JSON
    {
        "type": "conflict",
        "message": "custom_id is already used."
    }
    ```
    `"custom_id"`を設定した場合、すでにそのカスタムIDが存在していて登録できない場合に返されます。

4. 500 InternalServerError Responce
    ```JSON
    {
        "type": "internal_error",
        "message": "An unexpected error occurred. Please try again later."
    }
    ```
    サーバー側の問題が発生しています。一時的な問題の可能性もありますが、継続的に発生する場合はお問合せください。

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

Outtech105k

[Email](owner@outtech105.com)

[Twitter(X)](https://x.com/105techno)

# Set custom URL

`/api/set` にPOSTリクエストすることで、カスタムURLが生成できます。

## Requests

- `"expire_in"`キーは、正規表現`` `[0-9]+[smhd]` ``を受け付けます。
- `"10m", "30h", "2d"`のように指定してください。
指定しない場合、有効期限は設定されません。
- `"base_url"`のみの指定の場合、ランダムIDがセットされます。

### Rules & Constraints

URL生成に関する共通のルールや制約（禁止文字、予約語、有効期限の形式など）については、[URL生成ルールと制約](../rules.md) を確認してください。

| key | 説明 | 必須 | 競合する値 |
| :-- | :-- | :-- | :-- |
| `base_url` | リダイレクト先URL | 必須 | |
| `use_uppercase`| ランダムIDに英大文字を含めるか | | `custom_id` |
| `use_lowercase`| ランダムIDに英小文字を含めるか | | `custom_id` |
| `use_numbers`| ランダムIDに数字を含めるか | | `custom_id` |
| `id_length`| ランダムIDの文字数 | | `custom_id` |
| `custom_id`| 設定するカスタムID | | `use_uppercase`, `use_lowercase`, `use_numbers`, `id_length` |
| `expire_in`| リンクの有効期間 | | |
| `sand_cushion`| クッションページを使用するか | | |
| `public_ctrl`| カスタムURLの操作ページを公開するか | | |

## Request examples

1. ランダムIDでURL生成する例

```JSON
{
    "base_url": "https://example.com",
    "use_uppercase": true,
    "use_lowercase": false,
    "use_numbers": true,
    "id_length": 5,
    "expire_in": "10h",
    "public_ctrl": true
}
```

2. カスタムIDでURL生成する例

```JSON
{
    "base_url": "https://example.com",
    "custom_id": "example",
    "expire_in": "10h",
    "sand_cushion": true,
    "public_ctrl": false
}
```

### Response examples

1. 200 OK

```JSON
{
    "base_url": "https://example.com",
    "short_url": "https://rk2.uk/example",
    "warnings": ["custom_id contains multibyte characters; behavior may be unstable in some environments."]
}
```
正常にURLが生成されたものの、その過程で要注意事項がある場合 `warnings` がレスポンスに含まれます。

2. 400 Bad Request
    - Validation Error
    ```JSON
    {
        "type": "validation_error",
        "details": [
            {
                "field": "base_url",
                "message": "base_url is required."
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

3. 409 Conflict
    ```JSON
    {
        "type": "conflict",
        "message": "custom_id is already used."
    }
    ```
    `"custom_id"`を設定した場合、すでにそのカスタムIDが存在していて登録できない場合に返されます。

4. 500 InternalServerError
    ```JSON
    {
        "type": "internal_error",
        "message": "An unexpected error occurred. Please try again later."
    }
    ```
    サーバー側の問題が発生しています。発生したエラーはサーバー側のログで確認できます。

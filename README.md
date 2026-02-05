# grpc-connect-envoy

Connect (gRPC-Web/HTTP) と Envoy を活用したマイクロサービスアーキテクチャのデモプロジェクトです。
Gateway パターンと Sidecar パターンを組み合わせ、JWT による認証フローを実装しています。

## 特徴

- **Connect**: Go 言語による gRPC 互換のサービス実装
- **Envoy Proxy**:
  - **Gateway**: `envoy-client` がエッジプロキシとしてリクエストを集約
  - **Sidecar**: 各サービスの前段に配置され、通信制御や認証を担当
- **Security (Zero Trust)**:
  - **Sidecar Authentication**: Gateway ではなく、各サービスの直前 (`envoy-greet`) で認証を実施
  - **JWKS Endpoint**: `auth` サービスが公開鍵を配信し、Envoy が動的に取得してトークンを検証

## アーキテクチャ

```
                                     Docker Network
┌──────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                          │
│  ┌─────────┐      ┌──────────────┐      ┌─────────────┐      ┌─────────┐                 │
│  │ client  │      │ client-envoy │      │ envoy-auth  ├─────►│  auth   │                 │
│  │   app   ├─────►│   (Gateway)  ├─────►│   sidecar   │      │ service │                 │
│  └─────────┘      └──────┬───────┘      └──────▲──────┘      └─────────┘                 │
│                          │                     │                                         │
│                          │                     │ Fetch JWKS                              │
│                          │                     │                                         │
│                          │                ┌────┴────────┐ Verify Token                   │
│                          ├───────────────►│ envoy-greet │                                │
│                          │                │   sidecar   │                                │
│                          │                └──────┬──────┘                                │
│                          │                       │                                       │
│                          │                       ▼                                       │
│                          │                  ┌─────────┐                                  │
│                          │                  │  greet  │                                  │
│                          │                  │ service │                                  │
│                          │                  └─────────┘                                  │
│                          │                                                               │
│                          │                ┌─────────────┐      ┌─────────┐               │
│                          └───────────────►│ envoy-thanks├─────►│  thanks │               │
│                                           │   sidecar   │      │ service │               │
│                                           └─────────────┘      └─────────┘               │
│                                                                                          │
└──────────────────────────────────────────────────────────────────────────────────────────┘
```

## サービス構成

| サービス | ポート(内部) | 役割 | 認証 |
| --- | --- | --- | --- |
| **auth** | 9090 | ユーザー認証、トークン発行 (JWT)、JWKS配信 | 不要 |
| **greet** | 9090 | 挨拶を返すサービス | **必要** (Sidecarで検証) |
| **thanks** | 9090 | 感謝を返すサービス | 不要 |
| **client** | - | 動作確認用の CLI ツール | - |

- 全てのサービスへのアクセスは `envoy-client` (port 8080) を経由します。

## 実行方法

Docker Compose を使用して全サービスを起動します。

```bash
docker compose up --build
```

起動すると、`client` コンテナが自動的に以下のフローを実行し、ログに出力します。

1. `auth` サービスでログインし、JWT トークンを取得
2. トークンを使用して `greet` サービスへリクエスト（成功）
3. トークンなしで `thanks` サービスへリクエスト（成功）
4. トークンのリフレッシュ

### ログの確認

```bash
docker compose logs -f client
```

### 手動での動作確認 (curl)

ホストマシンから `localhost:8080` に対してリクエストを送ることで動作確認が可能です。

1. **トークン取得**
   ```bash
   TOKEN=$(curl -s -X POST http://localhost:8080/auth.v1.AuthService/Authenticate \
     -H "Content-Type: application/json" \
     -d '{"username":"user", "password":"password"}' | jq -r .token)
   echo $TOKEN
   ```

2. **Greet (認証あり)**
   ```bash
   curl -X POST http://localhost:8080/greet.v1.GreetService/Greet \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"name": "World"}'
   ```

3. **Thanks (認証なし)**
   ```bash
   curl -X POST http://localhost:8080/thanks.v1.ThanksService/Thanks \
     -H "Content-Type: application/json" \
     -d '{"name": "World"}'
   ```

## 技術詳細

### 認証フローの変更点

以前は Gateway (`envoy-client`) で一括して認証を行っていましたが、ゼロトラストの観点から **Sidecar (`envoy-greet`) での認証** に変更しました。

1. `auth` サービスに `/.well-known/jwks.json` エンドポイントを追加。
2. `envoy-greet` は `remote_jwks` 設定により、`auth` サービスから動的に公開鍵を取得。
3. リクエストが `envoy-greet` に到達した時点で JWT の検証が行われます。

## 参考リンク

- [Connect](https://connectrpc.com/)
- [Envoy Proxy](https://www.envoyproxy.io/)

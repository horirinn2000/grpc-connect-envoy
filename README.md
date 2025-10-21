# grpc-connect-envoy

## 概要

connectとenvoyを使ったweb-grpc実装

## 構成図
            ┌─────────┐   ┌──────────────────────────┐
            │ client  │   │      (Docker Network)    │
            │   app   │   │                          │
            └────┬────┘   │ ┌─────────┐  ┌─────────┐ │
                 │        │ │ greet   │  │ envoy   │ │
            ┌────▼────┐   │ │ service │<-│ sidecar │ │
            │ client  │   │ └─────────┘  └────▲────┘ │
            │  envoy  │───┼───────────────────┤      │
            └─────────┘   │ ┌─────────┐  ┌────▼────┐ │
                          │ │ thanks  │  │ envoy   │ │
                          │ │ service │<-│ sidecar │ │
                          │ └─────────┘  └─────────┘ │
                          │                          │
                          └──────────────────────────┘

## 参考URL

https://connectrpc.com/docs/introduction  
https://www.envoyproxy.io/

# Prokect: Apple App Store Server Golang SDK

## 1. Overview
The App Store Server Library is an open source library from Apple, available in four languages. It makes adopting the App Store Server API and working with JSON Web Signature (JWS) transactions easier. Find the App Store Server Library for each language in the following GitHub repositories:

    Swift: App Store Server Swift Library

    Java: App Store Server Java Library

    Python: App Store Server Python Library

    Node: App Store Server Node Library

官方SDK没有支持golang，故我们需要参考Python: App Store Server Python Library来完成一个golang SDK

## 2. 开发规范
你是简洁优雅的工程师，你的架构和编码要遵循go规范，要使用业界常用的编码方式
- **Clean Architecture**：分层设计，依赖单向流动。
- **DRY/KISS/YAGNI**：避免重复代码，保持简单，只实现必要功能。
- **代码可维护性**：模块化设计，清晰的包结构和函数命名。
- **易用性**：你开发的sdk易用使用
- **OSS**: 这是一个开源项目，尽可能使用标准库，引入第三方库需要向我申请

## 4. Project Layout
├── app-store-server-library-python # App Store Server Python Library
├── appstoreserver
│   └── v1 // apple api，本次开发应该在此目录中
│       ├── constants.go
│       └── datatypes.go
├── appstoreservernotifications 
│   └── v2 // apple 通知
│       ├── constants.go
│       └── datatypes.go
├── go.mod
├── go.sum
├── internal
│   ├── jwsutil # jws目录
│   └── jwtutil # jwt目录
│       └── datatypes.go
├── LICENSE
├── README.md
└── requirment.md # 需求文档

## 3. 根据apple文档，这是目前支持的接口
- /inApps/v2/history/{transactionId}
- /inApps/v1/transactions/{transactionId}
- /inApps/v1/subscriptions/{transactionId}
- /inApps/v1/transactions/{originalTransactionId}/appAccountToken
- /inApps/v1/transactions/consumption/{transactionId}
- /inApps/v1/lookup/{orderId}
- /inApps/v2/refund/lookup/{transactionId}
- /inApps/v1/subscriptions/extend/{originalTransactionId}
- /inApps/v1/subscriptions/extend/mass
- /inApps/v1/subscriptions/extend/mass/{productId}/{requestIdentifier}
- /inApps/v1/notifications/history
- /inApps/v1/notifications/test
- /inApps/v1/notifications/test/{testNotificationToken}

## 4. 开发流程
1. 了解app-store-server-library-python/appstoreserverlibrary的项目构成
2. 参考app-store-server-library-python/appstoreserverlibrary/signed_data_verifier.py构建jws工具包
3. app-store-server-library-python/appstoreserverlibrary/receipt_utility.py 已经弃用，不需要
4. app-store-server-library-python/appstoreserverlibrary/promotional_offer.py 不需要
5. 参考app-store-server-library-python/appstoreserverlibrary/api_client.py:构建jwt工具包，internal/jwtutil/datatypes.go已经有Claims
6. 开发接口，不清楚的可以和我索要Apple接口文档

## 限制
未经过容许，不可以随意引进第三方包，尽可能使用标准库完成开发
如有需要可以请求我
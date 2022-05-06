# Authorize-Api文档

## Summary

`Authorize`接口是CeylonPlatform提供OAuth服务的接口，主要包括下述内容：

+ 创建用户
+ 创建客户端
+ 申请token
+ 验证token
+ 刷新token

若需要在程序内调用oauth的实现，请转到 [oauth.md](oauth.md)

## 创建用户

此接口可以创建多个用户，考虑到本系统均为内部使用，暂时未限制本接口的调用

URL

    POST /authorize/users

**Body示例**

```json
{
    "user_list": [
        {
            "username": "user1",
            "password": "123456",
            "scope": 1
        },
        {
            "username": "user2",
            "password": "123456",
            "scope": 4
        }
    ]
}
```

**响应示例**

```json
{
    "user_list": [
        {
            "user_id": "986fae42ac817f2de0418a480bdb6b39",
            "username": "user1",
            "scope": 1,
            "create_at": "2022-05-06T17:17:02.588452+08:00"
        },
        {
            "user_id": "185b905f6186ce93dfc610bac5199ff0",
            "username": "user2",
            "scope": 4,
            "create_at": "2022-05-06T17:17:02.669297+08:00"
        }
    ]
}
```

## 创建客户端

此接口可以创建一个客户端，考虑到本系统均为内部使用，暂时未限制本接口的调用

**URL**

    POST /authorize/client

**Body示例**

```json
{
    "client_name": "client",
    "client_domain": "client.com",
    "client_scope": 1,
    "client_auth_type": "client"
}
```

**响应示例**

```json
{
    "client_id": "fc701433dfa60ca78cd7b43e06a6e092",
    "client_key": "6cdaa304287254d53c3bb807f6b2437b",
    "client_secret": "6552a1757d12baed02d7bd3a7f6dde19",
    "client_domain": "client.com",
    "scope": 2,
    "auth_type": "client",
    "create_at": "2028-05-06T12:22:15.805302+08:00"
}
```

## 申请token

关于四种申请 oauth-token 的模式，请转到 [oauth.md](oauth.md)

此接口可以使用四种方式中的一种申请 access-token 和 refresh-token，此接口进行了权限校验

**URL**

    POST /authorize/oauth

**Params**

    auth_type   - 必须
    client_id   - 必须
    scope       - 必须

**Body示例**

URL: `/authorize/oauth?auth_type=client`

```json
{
    "client_id": "fc701433dfa60ca78cd7b43e06a6e092",
    "client_secret": "6552a1757d12baed02d7bd3a7f6dde19",
    "scope": 1
}
```

**响应示例**

```json
{
    "token": "b9f669376d289d53cd13c124463a0092",
    "refresh_token": "681259555913376084183f830164e314",
    "expire_at": "2022-05-06T18:24:52.481842+08:00"
}
```

## 验证token

此接口可以验证一个 access-token 是否具有指定 scope 的资源访问权限

**URL**

    GET /authorize/introspect

**Params**

    auth_type   - 必须

**Body示例**

```json
{
    "token": "6a53ef8aa2a4120754ae688ab6f3f731",
    "client_id": "fc701433dfa60ca78cd7b43e06a6e092",
    "scope": 1
}
```

**响应示例**
```json
{
    "status": true
}
```
# OAuth中间件

此处是CeylonPlatform的OAuth模块文档，用于实现Client/Manager与Backend之间的通信。

## 约定

在开始之前，我们需要明确一些内容：

+ Client: 一个Client是运行在学生计算机上的程序，用于访问CeylonSystem的数字资源
+ Manager: 一个Manager是运行在教师计算机上的程序，用于管理(增、删、改)CeylonSystem的数字资源
+ Backend: 一个Backend是运行在服务器上的程序，本CeylonPlatform是Backend的一个实现
+ Repository: 一个Repository是运行在服务器上的程序，用于托管具体资源，本系统中Repository的实现是gitea

由于本项目使用C/S架构，Client和Manager是运行在Windows上的.Net Framework客户端，所以我们在Client上直接使用password授权模式，在Manager上使用client授权模式，implicit和code模式将在后续使用

## 流程

在使用Client的情况下，OAuth中间件使用密码方式(`PasswordAuth`)进行授权，整个访问的过程如下：

1. Client在开发时向Backend申请注册，获取app-key与app-secret，将获取方法或获取值写进Client的业务逻辑中
2. 学生打开Client，使用Client向Backend发起授权请求
3. 学生授权Client使用其权限内(Scope)的资源，Backend将授权与Token返回给Client
4. Client使用自身的`Client_Token`与学生的`User_Token`访问Backend，Backend鉴权后将Repository的资源路径返回给Client
5. Client访问资源并提供给学生使用

## 接口

OAuth中间件中，默认的Authenticator实现提供如下方法：

### 基础entity的CURD

+ CreateClient() (client *Client, err error)

    ```go
    // @description 生成并返回一个空白的Client，此方法会写数据库
    // @return{client} 生成后的client数据
    // @return{err} 生成时产生的错误，若为空则为无错误
    func (a Authenticator)CreateClient() (client *Client, err error)
    ```

+ CreateClientWith(opts *ClientOptions) (client *Client, err error)

    ```go
    // @description 使用options生成并返回Client，此方法会写数据库
    // @param{opts} 提供的client基础信息
    // @return{client} 生成后的client数据
    // @return{err} 生成时产生的错误，若为空则为无错误
    func (a Authenticator)CreateClientWith (opts *ClientOptions) (client *Client, err error)
    ```

+ UpdateClientWith(opts *ClientOptions) (client *Client, err error)

    ```go
    // @description 使用options更新并返回Client，此方法会写数据库
    // @param{opts} 提供的client基础信息
    // @return{client} 更新后的client数据
    // @return{err} 更新时产生的错误，若为空则为无错误
    func (a Authenticator)UpdateClientWith(opts *ClientOptions) (client *Client, err error)
    ```
  
+ DeleteClient(opts *ClientOptions) (err error)

    ```go
    // @description 删除Client，此方法并不会真正删除Client记录，只会修改一条记录
    // @param{opts} 删除的client的信息，全部满足才会删除，为空则不删除
    // @return{ok} 删除状态，true则为删除成功
    // @return{err} 删除时产生的错误，若为空则为无错误
    func (a Authenticator)DeleteClient(opts *ClientOptions) (ok bool, err error)
    ```

+ CreateUser() (user *User, err error)

    ```go
    // @description 生成并返回一个空白的User，此方法会写数据库
    // @return{user} 生成后的user数据
    // @return{err} 生成时产生的错误，若为空则为无错误
    func (a Authenticator) CreateUser() (user *User, err error)
    ```

+ CreateUserWith(opts *UserOptions) (user *User, err error)

    ```go
    // @description 使用options生成并返回User，此方法会写数据库
    // @params{opts} 提供的user基础信息
    // @return{user} 生成后的user数据
    // @return{err} 生成时产生的错误，若为空则为无错误
    func (a Authenticator) CreateUserWith(opts *UserOptions) (user *User, err error)
    ```

+ UpdateUserWith(opts *UserOptions) (user *User, err error)

    ```go
    // @description 使用options更新并返回User，此方法会写数据库
    // @params{opts} 提供的user基础信息
    // @return{user} 更新后的user数据
    // @return{err} 更新时产生的错误，若为空则为无错误
    func (a Authenticator) UpdateUserWith(opts *UserOptions) (user *User, err error)
    ```

+ DeleteUser(opts *UserOptions) (ok bool, err error)

    ```go
    // @description 删除User，此方法并不会真正删除User记录，只会修改一条记录
    // @param{opts} 删除的user信息，全部满足才会删除，为空则不删除
    // @return{ok} 删除状态，true则为删除成功
    // @return{err} 删除时产生的错误，若为空则为无错误
    func (a Authenticator)DeleteUser(opts *UserOptions) (ok bool, err error)
    ```

+ CreateAccessToken() (token *AccessToken, err error)

    ```go
    // @description 生成并返回一个空白的AccessToken，此方法会写数据库
    // @return{token} 生成后的access-token数据
    // @return{err} 生成时产生的错误，若为空则为无错误
    func (a Authenticator) CreateAccessToken() (token *AccessToken, err error)
    ```


+ CreateAccessTokenWith(opts *AccessTokenOptions) (token *AccessToken, err error)

    ```go
    // @description 使用options生成并返回AccessToken，此方法会写数据库
    // @param{opts} 提供的access-token数据
    // @return{token} 生成后的access-token数据
    // @return{err} 生成时产生的错误，若为空则为无错误
    func (a Authenticator) CreateAccessTokenWith(opts *AccessTokenOptions) (token *AccessToken, err error)
    ```

+ UpdateAccessTokenWith(opts *AccessTokenOptions) (token *AccessToken, err error)

    ```go
    // @description 使用options更新并返回AccessToken，此方法会写数据库
    // @param{opts} 提供的access-token数据
    // @return{token} 更新后的access-token数据
    // @return{err} 更新时产生的错误，若为空则为无错误
    func (a Authenticator) UpdateAccessTokenWith(opts *AccessTokenOptions) (token *AccessToken, err error)
    ```

+ DeleteAccessToken(opts *AccessTokenOptions) (ok bool, err error)

    ```go
    // @description 删除AccessToken，此方法并不会真正删除AccessToken记录，只会修改一条记录
    // @param{opts} 删除的access-token的信息，全部满足才会删除，为空则不删除
    // @return{ok} 删除记录，true则为删除成功
    // @return{err} 删除时产生的错误，若为空则为无错误 
    func (a Authenticator) DeleteAccessToken(opts *AccessTokenOptions) (ok bool, err error)
    ```

+ CreateRefreshToken() (token *RefreshToken, err error)

    ```go
    // @description 创建并返回一个空白的RefreshToken，此方法会写数据库
    // @return{token} 生成后的refresh-token数据
    // @return{err} 生成时产生的错误，若为空则为无错误
    func (a Authenticator) CreateRefreshToken() (token *RefreshToken, err error)
    ```

+ CreateRefreshTokenWith(opts *RefreshTokenOptions) (token *RefreshToken, err error)

    ```go
    // @description 使用options生成并返回RefreshToken，此方法会写数据库
    // @param{opts} 提供的refresh-token数据
    // @return{token} 生成的refresh-token数据
    // @return{err} 生成时产生的错误，若为空则为无错误
    func (a Authenticator) CreateRefreshTokenWith(opts *RefreshTokenOptions) (token *RefreshToken, err error)
    ```

+ UpdateRefreshTokenWith(opts *RefreshTokenOptions) (token *RefreshToken, err error)

    ```go
    // @description 使用options更新并返回RefreshToken，此方法会写数据库
    // @param{opts} 提供的refresh-token数据
    // @return{token} 更新后的refresh-token数据
    // @return{err} 更新时产生的错误，若为空则为无错误
    func (a Authenticator) UpdateRefreshTokenWith(opts *RefreshTokenOptions) (token *RefreshToken, err error)
    ```

+ DeleteRefreshTokenWith(opts *RefreshTokenOptions) (ok bool, err error)

    ```go
    // @description 删除RefreshToken，此方法并不会真正删除RefreshToken记录，只会修改一条记录
    // @param{opts} 删除的refresh-token的信息，全部满足才会删除，为空则不删除
    // @return{ok} 删除记录，true则为删除成功
    // @return{err} 删除时产生的错误，若为空则为无错误 
    func (a Authenticator) DeleteRefreshTokenWith(opts *RefreshTokenOptions) (ok bool, err error)
    ```

+ CreateAuthorizationCode() (token *AuthorizationCode, err error)

    ```go
    // @description 创建并返回一个空白的AuthorizationCode，此方法会写数据库
    // @return{token} 生成后的authorization-code数据
    // @return{err} 生成时产生的错误，若为空则为无错误
    func (a Authenticator) CreateAuthorizationCode() (token *AuthorizationCode, err error)
    ```

+ CreateAuthorizationCodeWith(opts *AuthorizationCodeOptions) (token *AuthorizationCode, err error)

    ```go
    // @description 使用options生成并返回AuthorizationCode，此方法会写数据库
    // @param{opts} 提供的authorization-code数据
    // @return{token} 生成的authorization-code数据
    // @return{err} 生成时产生的错误，若为空则为无错误
    func (a Authenticator) CreateAuthorizationCodeWith(opts *AuthorizationCodeOptions) (token *AuthorizationCode, err error)
    ```

+ UpdateAuthorizationCodeWith(opts *AuthorizationCodeOptions) (token *AuthorizationCode, err error)

    ```go
    // @description 使用options更新并返回AuthorizationCode，此方法会写数据库
    // @param{opts} 提供的authorization-code数据
    // @return{token} 更新后的authorization-code数据
    // @return{err} 更新时产生的错误，若为空则为无错误
    func (a Authenticator) UpdateAuthorizationCodeWith(opts *AuthorizationCodeOptions) (token *AuthorizationCode, err error)
    ```

+ DeleteAuthorizationCodeWith(opts *AuthorizationCodeOptions) (ok bool, err error)

    ```go
    // @description 删除AuthorizationCode，此方法并不会真正删除AuthorizationCode记录，只会修改一条记录
    // @param{opts} 删除的authorization-code的信息，全部满足才会删除，为空则不删除
    // @return{ok} 删除记录，true则为删除成功
    // @return{err} 删除时产生的错误，若为空则为无错误 
    func (a Authenticator) DeleteAuthorizationCodeWith(opts *AuthorizationCodeOptions) (ok bool, err error)
    ```

### OAuth的业务封装

#### 使用password授权模式

+ PasswordAuth(userID, clientID, password string, scope ScopeType) (token *AccessToken, refreshToken *RefreshToken, err error)

    ```go
    // @description 使用password模式进行认证
    // @params{userID} user的ID
    // @params{clientID} client的ID
    // @params{password} 用户的密码
    // @params{scope} 授权的范围
    // @return{token} 认证成功的access-token，认证失败则为空
    // @return{refreshToken} 认证成功的refresh-token，认证失败则为空
    // @return{err} 认证过程中产生的错误，若为空则为无错误
    func (a Authenticator) PasswordAuth(userID, clientID, password string, scope ScopeType) (token *AccessToken, refreshToken *RefreshToken, err error)
    ```

#### 使用client授权模式

+ ClientAuth(clientID, clientKey, clientSecret string, scope ScopeType) (token *AccessToken, refreshToken *RefreshToken, err error)

    ```go
    // @description 使用client模式进行认证
    // @params{clientID} client的ID
    // @params{clientSecret} client的Secret
    // @params{scope} 授权的范围
    // @return{token} 认证成功的access-token，认证失败则为空
    // @return{refreshToken} 认证成功的refresh-token，认证失败则为空
    // @return{err} 认证过程中产生的错误，若为空则为无错误
    func (a Authenticator) ClientAuth(clientID, clientSecret string, scope ScopeType) (token *AccessToken, refreshToken *RefreshToken, err error)
    ```

#### 使用implicit授权模式

+ ImplicitAuth(clientID string, scope ScopeType) (redirectUri string, err error)

    ```go
    // @description 使用implicit模式进行认证
    // @params{clientID} client的ID
    // @params{redirectUri} 重定向URI
    // @params{scope} 授权的范围
    // @return{uri} 拼接完成的重定向URI
    // @return{err} 认证过程中产生的错误，若为空则为无错误
    func (a Authenticator) ImplicitAuth(clientID, redirectUri string, scope ScopeType) (uri string, err error)
    ```
  
#### 使用code授权模式

+ CodeAuth(clientID, redirectUri string, scope ScopeType) (uri string, err error)

    ```go
    // @description 在code模式进行认证，生成code
    // @params{clientID} client的ID
    // @params{redirectUri} 重定向URI
    // @params{scope} 授权的范围
    // @return{uri} 拼接完成的重定向URI
    // @return{err} 认证过程中产生的错误，若为空则为无错误
    func (a Authenticator) CodeAuth(clientID, redirectUri string, scope ScopeType) (uri string, err error)
    ```

+ CodeToToken(clientID, clientSecret, code string, scope ScopeType) (token *AccessToken, refreshToken *RefreshToken, err error)

    ```go
    // @description 在code模式进行认证，根据code生成token
    // @params{clientID} client的ID
    // @params{clientSecret} client的Secret
    // @params{code} 上一个步骤中生成的code
    // @params{scope} 授权的范围
    // @return{token} 认证成功的access-token，认证失败则为空
    // @return{refreshToken} 认证成功的refresh-token，认证失败则为空
    // @return{err} 认证过程中产生的错误，若为空则为无错误
    func (a Authenticator) CodeToToken(clientID, clientSecret, code string, scope ScopeType) (token *AccessToken, refreshToken *RefreshToken, err error)
    ```
  
### 通用业务接口

#### Token相关

+ AuthToken(authType AuthType, opts *TokenOptions) (ok bool, err error)

    ```go
    // @description 根据options验证token的正确性
    // @params{authType} 认证的模式
    // @params{opts} 提供的token的信息
    // @return{ok} token是否正确，true则为正确
    // @return{err} 认证过程中产生的错误，若为空则为无错误
    func (a Authenticator) AuthToken(authType AuthType, opts *TokenOptions) (ok bool, err error)
    ```

+ RefreshToken(authType AuthType, opts *TokenOptions) (token *AccessToken, refreshToken *RefreshToken, err error)

    ```go
    // @description 根据options刷新token
    // @params{authType} 认证的模式
    // @params{opts} 提供的token的信息
    // @return{token} 认证成功的access-token，认证失败则为空
    // @return{refreshToken} 认证成功的refresh-token，认证失败则为空
    // @return{err} 认证过程中产生的错误，若为空则为无错误
    func (a Authenticator) RefreshToken(authType AuthType, opts *TokenOptions) (token *AccessToken, refreshToken *RefreshToken, err error)
    ```
  

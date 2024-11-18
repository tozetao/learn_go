### TDD

先写测试、再写实现。



为什么要用TDD

TDD是站在用户的角度，站在使用接口方的角度上来看待接口该如何设计，如何实现，有什么样的输入和输出。

- 通过撰写测试，理清楚接口该如何定义，体会用户使用是否合适。

- 通过撰写测试用例，理清楚整个功能要考虑的主流程、异常流程。

TDD专注某个功能的实现。



大体流程：

1. 定于接口

2. 定义测试模板

   定义测试模板需要考虑有什么输入，有什么输出。

3. 写测试用例

   一个测试用例就对应一个流程，TDD的测试用例更关注你的实现细节。如果你是从集成测试出发，它真实的对应到你的业务流程。

4. 提供或修改代码

5. 执行测试

   测试通过说明流程没有问题，可以进行下一个测试，考虑下一个场景。如果未通过需要回去修改代码，直到测试通过。

通过不断的增加测试用例，把所有流程都考虑完毕，并且所有测试用例都执行通过后，这个功能也就完善了。



TDD核心循环：

1. 根据对需求的理解，初步定义接口。

   不需要害怕定义的接口不合适。

2. 根据接口定义测试。

   即根据测试模板，先把测试的框架写出来。

3. 执行核心循环

   增加测试用例 => 提供/修改实现 => 执行测试用例。





### 测试

- 测试文件以xxx_test命名。
- 测试方法以Test开头。



依赖的包：

- mock包：https://github.com/uber-go/mock
- sqlmock: https://github.com/DATA-DOG/go-sqlmock



```
-- 安装命令行工具
go install go.uber.org/mock/mockgen@latest

-- 生成mock文件
mockgen -source=.../internal/service/user.go -package=svcmocks -destination=webook/internel/service/mocks/user.mock.go

-- 安装依赖包
go mod tidy
```







```
# service
mockgen -source=./internal/service/user.go -package=svcmocks -destination=./internal/service/mocks/user.mock.go

mockgen -source=./internal/service/article.go -package=svcmocks -destination=./internal/service/mocks/article.mock.go

mockgen -source=./internal/repository/article/article_author.go -package=artrepomocks -destination=./internal/repository/mocks/article/article_author.mock.go

mockgen -source=./internal/repository/article/article_reader.go -package=artrepomocks -destination=./internal/repository/mocks/article/article_reader.mock.go

# repository
mockgen -source=./internal/repository/user.go -package=repomocks -destination=./internal/repository/mocks/user.mock.go

mockgen -source=./internal/repository/dao/user.go -package=daomocks -destination=./internal/repository/dao/mocks/user.mock.go

mockgen -source=./internal/repository/cache/user.go -package=cachemocks -destination=./internal/repository/cache/mocks/user.mock.go


# redis
mockgen -package=redismocks -destination=./internal/repository/cache/redismocks/redis.mock.go github.com/redis/go-redis/v9 Cmdable
```


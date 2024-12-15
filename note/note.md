https://gitee.com/geektime-geekbang_admin/geektime-basic-go



#### 进度

```
完成榜单模型的计算
完成可观测性章节的代码
```





#### 提示

安装go包的可执行文件

```
git tag
git checkout -b install v.x.x
```





#### 问题

````


### kafka相关
- 怎么知道一个消费者组中有多少个消费者，每个消费者在消费主题的哪个分区？
- 测试kafka分片消费消息的顺序

- kafka如何设置数据保存的时间？
- sarama的偏移量设置，最新和最旧有什么区别?
  OffsetNewest、OffsetOldest

- 如果多个goroutine返回错误，那么errGroup.Wait()究竟返回的是哪个错误？





!(a && b)等价于什么
!(a || b)等价于什么






资源：

https://github.com/ecodeclub/ekit



目标：
基础：切片的辅助方法、map的辅助方法，用内置map封装一个set
中级：设计List、普通队列、HashMap
高级：基于树形结构衍生出来的类型、基于跳表衍生出来的类型、ben copier机制。


实现切片的删除操作
- 考虑高性能操作
- 改造成泛型方法
- 支持缩容。

切片辅助方法
- 求和
- 求最大值、最小值
- 添加、删除、查找、过滤、Map Reduce。
- 集合运算：交集、并集、差集
````


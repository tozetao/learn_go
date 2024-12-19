在go中，context.Context承担了俩大职责：

- 在goroutine中内传递数据，类似别的语言中的thread local。

  ```go
  // 设置键值对，并且返回一个新的context
  context.WithValue(key, value)
  ```

- 用于超时控制，或者取消执行

  ```
  // 三者都返回一个可取消的context示例和可取消的函数
  context.WithCancel
  context.WithDeadline
  context.WithTimeout
  ```



**Context接口**

核心API有4个：

- Deadline() (deadline time.Time, ok bool)

  返回过期时间，如果ok为false表示没有设置过期时间。

- Done() <-chan struct{}

  返回一个channel。一般用于监听context实例的信号。比如过期，或者正常关闭。

- Err()

  返回一个错误用于表示context发生了什么。Canceled错误表示主动关闭，就是调用了对应的cancel；DeadlineExceeded表示过期超时。

- Value

  取值



context的实例之间存在父子关系：

- 当父亲取消或者超时，所有派生出的context都会被取消或者超时。控制是从上到下的。

- 在寻找key的时候，子context找不到会向上去祖先里面找。查找是从下往上查找的。





> context是怎么做到线程安全的？
>
> context是一个不可变对象，当一个context被创建后将不可修改，只能从当前的context派生出子context。




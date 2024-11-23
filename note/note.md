https://gitee.com/geektime-geekbang_admin/geektime-basic-go









进度：

```
新建并发表
编辑后首次发表
编辑后发表
```

安装go包的可执行文件

```
git tag
git checkout -b install v.x.x
```



问题

```
为什么SQL语句的where查询条件会重复？
err := dao.db.WithContext(ctx).Model(&article).
    Where("id = ?", article.ID).
    Updates(map[string]any{
        "title":   article.Title,
        "content": article.Content,
        "u_time":  article.Utime,
    }).Error
return err


插入时如果id存在则进行更新，下面的代码会是什么SQL语句。

// clause.OnConflict: 指定发生冲突时的行为。在这里是当列id发生冲突时，应该更新title和content。
return dao.db.Clauses(clause.OnConflict{
    Columns: []clause.Column{{Name: "id"}},
    DoUpdates: clause.Assignments(map[string]interface{}{
        "title":   article.Title,
        "content": article.Content,
    }),
}).Create(&article).Error

INSERT INTO your_table (id, title, content)
VALUES (1, 'New Title', 'New Content')
ON DUPLICATE KEY UPDATE
title = VALUES(title),
content = VALUES(content);
```



````
问题

- 如何让发送的jwt token过期？

- 限流：根据ip进行限流，限流速率的大小该如何调整？





### week2



14：12





### week1



资源：

https://github.com/ecodeclub/ekit



目标：
基础：切片的辅助方法、map的辅助方法，用内置map封装一个set
中级：设计List、普通队列、HashMap
高级：基于树形结构衍生出来的类型、基于跳表衍生出来的类型、ben copier机制。



!(a && b)等价于什么
!(a || b)等价于什么



作业：

```
实现切片的删除操作
- 考虑高性能操作
- 改造成泛型方法
- 支持缩容。

切片辅助方法
- 求和
- 求最大值、最小值
- 添加、删除、查找、过滤、Map Reduce。
- 集合运算：交集、并集、差集
```


````


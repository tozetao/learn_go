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

```




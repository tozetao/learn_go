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




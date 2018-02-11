# gosql

golang sql builder tool, similar to mybatis.

follow the example of <a href="http://github.com/jmoiron/sqlx/">github.com/jmoiron/sqlx</a>

```
go get github.com/Xu-Rui/gosql
```
---
### example

1. \# represents string splicing
2. $  represents prepare statement placeholder

 test sql
```sql
select * from #tablename where driver_id = $driver_id #sort
```

test table info (MySQL 5.7.18)
```sql
--table info
CREATE TABLE `driver_info` (
  `id` int(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `driver_id` int(20) NOT NULL DEFAULT '0' COMMENT '司机ID',
  `name` varchar(20) COLLATE utf8_bin NOT NULL DEFAULT '0' COMMENT '司机姓名',
  `age` int(20) NOT NULL DEFAULT '0' COMMENT '司机年龄',
  PRIMARY KEY (`id`),
  KEY `inx_driver_id_age` (`driver_id`,`age`)
) ENGINE=InnoDB AUTO_INCREMENT=124 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

--row info
INSERT INTO `driver_info` (`id`, `driver_id`, `name`, `age`)
VALUES
	(123,456,X'7875727569',12),
	(124,457,X'7875727569',13),
	(125,458,X'7875727569',14),
	(126,459,X'7875727569',15);
```
[demo](https://github.com/Xu-Rui/gosql/blob/master/executor_test.go)

---

![](http://photo.rhymecode.com/%E6%8F%92%E5%9B%BE/icon.png)

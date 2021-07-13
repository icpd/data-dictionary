## 生成MySQL数据字典

### 方式一、通过SQL生成数据字典
```sql
USE information_schema;

SELECT
	T.TABLE_SCHEMA AS '数据库名称',
	T.TABLE_NAME AS '表名称',
	T.TABLE_COMMENT AS '表说明',
	T.TABLE_TYPE AS '表类型',
	T.ENGINE AS '数据库引擎',
	C.ORDINAL_POSITION AS '序号',
	C.COLUMN_NAME AS '字段名',
	C.COLUMN_TYPE AS '数据类型',
	C.IS_NULLABLE AS '允许为空',
	C.EXTRA AS '自增属性',
	C.CHARACTER_SET_NAME AS '编码名称',
	C.COLUMN_DEFAULT AS '默认值',
	C.COLUMN_COMMENT AS '字段说明' 
FROM
	COLUMNS C
	INNER JOIN TABLES T ON C.TABLE_SCHEMA = T.TABLE_SCHEMA 
	AND C.TABLE_NAME = T.TABLE_NAME 
WHERE
	T.TABLE_SCHEMA = 'test';
```

### 方式二、使用该项目命令生成markdown表格
```shell
./main -d='username:password@tcp(127.0.0.1:3306)/information_schema' -s=database
```
完成后在目录下会生成一个`db.md`的文件

### 参考
- https://github.com/jayknoxqu/data-dictionary
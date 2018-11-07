# TopURL

从超大日志文件中获取出现频率 Top-K 的 URL

# 需求描述

现有 100GB url 文件

仅有 1GB 内存情况下（硬盘无大小限制）

需计算出现次数 top100 的 url 和出现的次数

# 解决方案

使用 map reduce 的思路，将问题分解

通过对 url 字符串计算哈希的方式，将大文件中 url 写入约 100 个小文件中

这样可以保证同一个 url 只出现在一个小文件中

然后对每个小文件中的 url 进行 hashmap 计数，并依次将 url 加入一个限制大小的小根堆中

遍历完所有小文件之后得到的堆中的 url 便是结果

# 使用方式

- 单元测试

```bash
go test
```

- 运行程序

假设 url 日志文件（/path/to/urls）格式如下：

```
http://www.qq.com/
http://www.baidu.com/
... ...
http://ecl.me/

```

运行如下命令，结果会直接输出到 stdout

```bash
go run . --file=/path/to/urls
```

# 性能测试

CPU: 2.5 GHz Intel Core i7

Memory: 16 GB 1600 MHz DDR3

## 1GB Test

```bash
./generator.py 1 test1gb
go run . --file=test1gb
// program took 10m21.245904647s
```

## 10GB Test

```bash
./generator.py 10 test10gb
go run . --file=test10gb
// program took 1h50m45.67196501s
```

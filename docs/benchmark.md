# verax 性能基准报告

对照对象: ozzo-validation v4 (本地源码, 同机同链路)。
基准代码位于 `bench/`(独立子模块, 通过 replace 引用两个库的本地副本)。

## 测试环境

| 项目 | 值 |
|------|-----|
| OS / 架构 | linux / amd64 |
| CPU | 12th Gen Intel Core i5-12400 (12 线程) |
| Go | 1.26.0 |
| 方法 | `go test -bench=. -benchmem`, 场景均为热路径稳态 |

## 场景设计

| 场景 | 规则链 | 数据规模 |
|------|--------|---------|
| Value | Required + Length(5,100) + EmailFormat | 单字符串 |
| Struct | 3 字段(Name/Email/Age), 各 2 条规则 | 单对象, 每次 new |
| Slice | Required + Length(1,20) 逐元素 | 100 元素切片 |

每个场景分 Pass(全部通过)与 Fail(中途失败)两条路径。

## 结果

| 基准 | verax ns/op | ozzo ns/op | 速度倍数 | verax allocs | ozzo allocs |
|------|------------:|-----------:|---------:|-------------:|------------:|
| ValuePass    |        625 |       1247 | 2.0x |  4 |   7 |
| ValueFail    |        332 |        400 | 1.2x |  4 |   6 |
| StructPass   |       1375 |       1984 | 1.4x | 11 |  17 |
| StructFail   |       1237 |       2248 | 1.8x | 12 |  21 |
| SlicePass    |        715 |       7651 | 10.7x |  1 | 102 |
| SliceFail    |        393 |       5859 | 14.9x |  0 | 104 |

## 分析

1. 切片场景优势最大(10 到 15 倍):
   verax 直接按类型遍历切片并复用预实例化的规则函数,
   ozzo 对每个元素做 reflect 取值与接口分发, 产生上百次分配。
2. struct 场景 1.4 到 1.8 倍:
   Field 绑定字段指针, 校验期零反射;
   ozzo 需要反射解析字段与规则匹配。
   StructPass 中 verax 的单次字节量略高(B/op +423),
   来自聚合错误容器的一次性预分配, 但分配次数更少且总吞吐更高。
3. 失败路径同样全面领先:
   错误对象在规则构造期预格式化, 失败时只做查表与装箱,
   不存在消息模板的重复拼接。
4. 所有场景分配次数均少于 ozzo, 高并发服务下可显著降低 GC 压力。

## 复现方式

```
cd bench
go test -bench=. -benchmem -run=^$
```

注意: 基准数字受机器负载影响, 关注相对倍数量级而非绝对值。

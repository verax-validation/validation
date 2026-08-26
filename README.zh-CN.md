# verax

[English](README.md) | **简体中文**

[![CI](https://github.com/verax-validation/validation/actions/workflows/ci.yml/badge.svg)](https://github.com/verax-validation/validation/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/verax-validation/validation.svg)](https://pkg.go.dev/github.com/verax-validation/validation)
[![Go Version](https://img.shields.io/github/go-mod/go-version/verax-validation/validation)](https://go.dev/dl/)
[![Release](https://img.shields.io/github/v/release/verax-validation/validation)](https://github.com/verax-validation/validation/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/verax-validation/validation)](https://goreportcard.com/report/github.com/verax-validation/validation)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

基于 Go 1.26 泛型的数据验证库。规则即函数, 类型安全由编译器保证。

## 目录

- [特性](#特性)
- [安装](#安装)
- [快速上手](#快速上手)
- [内置规则速览](#内置规则速览)
- [使用注意事项](#使用注意事项)
- [错误处理](#错误处理)
- [国际化](#国际化)
- [与其他库对比](#与其他库对比)
- [版本与稳定性](#版本与稳定性)
- [贡献](#贡献)
- [License](#license)

## 特性

- **全泛型 API**: `Rule[T]` 让规则与值类型在编译期绑定, 传错规则无法通过编译
- **链式字段配置**: `Field[T]()` 通过 option 模式绑定指针、追加规则、设置标签与自定义错误
- **跨字段比较**: `WithFieldEq` / `WithFieldGte` 等内置方法, 一行表达"确认密码一致""结束时间 ≥ 开始时间"
- **条件必填**: `rules.RequiredIf` 让"选了支付方式就必须填卡号"这类约束声明式表达
- **高性能**: 切片逐元素校验比反射式实现快一个数量级, 详见 [基准报告](docs/benchmark.md)
- **显式可选语义**: 规则默认无条件执行, 可选字段用 `Optional` 显式声明, 无隐藏行为
- **进程级国际化**: 内置六种语言, 语言模板保留占位符, 规则参数在渲染时插值
- **零第三方依赖**: 全部基于标准库实现
- **结构化错误**: 错误码 + 消息模板插值 + 字段上下文, 支持 `errors.AsType` 提取

## 安装

```
go get github.com/verax-validation/validation
```

要求 Go 1.26 或以上。

## 快速上手

### 校验简单值

```go
import (
    "github.com/verax-validation/validation"
    "github.com/verax-validation/validation/rules"
)

err := verax.Validate("alice@example.com",
    rules.Required,
    rules.Length[string](5, 100),
)
// 裸泛型规则在调用处自动推断类型, 例如 rules.Required 无需写成 rules.Required[string]
```

### 校验 struct

```go
type User struct {
    Name  string
    Email string
    Age   int
    Bio   string
}

func (u *User) Validate() error {
    return verax.ValidateStruct(u,
        verax.Field[string]().
            WithField(&u.Name).
            WithRules(rules.Required, rules.Length[string](2, 64)),
        verax.Field[string]().
            WithField(&u.Email).
            WithRules(rules.Required, is.EmailFormat).
            WithLabel("邮箱"),
        // 可选字段: 空值跳过后续规则
        verax.Field[string]().
            WithField(&u.Bio).
            WithRules(rules.Optional(rules.Length[string](0, 200))),
    )
}
```

字段通过 option 模式链式配置: `WithField` 绑定指针, `WithRules` 追加规则,
`WithLabel` 设置标签(错误前缀, 如 "邮箱不能为空"), `WithErr` 自定义错误,
`WithCheckFn` 无参自定义校验。元素类型与规则值类型不匹配时编译报错。

### 跨字段比较

```go
type RegisterForm struct {
    Password string
    Confirm  string
    StartAt  int
    EndAt    int
}

err := verax.ValidateStruct(&form,
    // 确认密码必须与密码一致
    verax.Field[string]().WithField(&form.Confirm).WithFieldEq(&form.Password),
    // 结束时间不能早于开始时间
    verax.Field[int]().WithField(&form.EndAt).WithFieldGte(&form.StartAt),
)
// 失败信息归属到当前字段并引用对方字段名: "confirm: 必须等于 password"
```

`WithFieldEq` / `WithFieldNe` / `WithFieldGt` / `WithFieldGte` / `WithFieldLt` / `WithFieldLte`
六个比较方法覆盖相等与有序比较, 同样支持 `WithLabel` / `WithErr` 组合。

### 条件必填

```go
type Order struct {
    PaymentMethod string
    CardNumber    string
}

err := verax.ValidateStruct(&order,
    verax.Field[string]().
        WithField(&order.CardNumber).
        WithRules(
            // 选择了在线支付时, 卡号必填
            rules.RequiredIf[string](func() bool { return order.PaymentMethod == "card" }),
            is.CreditCard,
        ),
)
```

### 嵌套对象

```go
type Address struct {
    City string
}

func (a Address) Validate() error {
    return verax.Validate(a.City, rules.Required)
}

err := verax.ValidateStruct(user,
    verax.Field[Address]().WithField(&user.Address).WithRules(verax.Valid),
)
```

### 校验集合

```go
import "github.com/verax-validation/validation/collections"

err := verax.Validate(tags,
    collections.SliceLen[string](1, 10),            // 数量约束
    collections.Unique[string](),                   // 元素互不重复
    collections.Slice[string](                      // 逐元素校验
        rules.Required,
        rules.Length[string](1, 20),
    ),
)
// 元素错误按位置聚合: Errors{"1": cannot be blank}
```

`collections` 提供 `Slice` / `Each`(iter.Seq) / `Map` 逐元素校验, 以及
`SliceLen` / `MapLen` 数量与 `Unique` 唯一性约束。

### 条件规则

```go
// isAdmin 为业务侧已确定的条件, 在构造规则时求值
rule := verax.When(isAdmin, rules.Min(-1))
```

### 完整示例

完整的可运行示例见 [examples/basic](examples/basic), 覆盖 struct 校验、跨字段比较、
集合校验、中文国际化与错误处理, 可直接运行:

```
go run ./examples/basic
```

## 内置规则速览

### rules 包 — 通用值规则

适用于任意类型, 通过 `verax.Validate` 或 `Field.WithRules` 使用。

**必填与可选**

| 规则 | 说明 |
|------|------|
| `Required` | 值不是其类型的零值(空串 / 0 / false / nil 指针均判空) |
| `RequiredIf(cond)` | 条件为真时值必填 |
| `Optional(rules...)` | 零值直接放行, 非零值才执行后续规则 |

**长度**(字符串与字节切片)

| 规则 | 说明 |
|------|------|
| `Length(min, max)` | 字节长度落在闭区间 `[min, max]` |
| `Len(n)` | 字节长度精确等于 n |

**数值与区间**(有序类型: 整数 / 浮点 / 字符串 / time 等)

| 规则 | 说明 |
|------|------|
| `Min(min)` | 大于等于 min |
| `Max(max)` | 小于等于 max |
| `Between(lo, hi)` | 落在闭区间 `[lo, hi]` |
| `Gt(min)` | 严格大于 min(与闭区间的 Min 互补) |
| `Lt(max)` | 严格小于 max(与闭区间的 Max 互补) |
| `MultipleOf(step)` | 是 step 的整数倍 |

**等值与枚举**

| 规则 | 说明 |
|------|------|
| `Eq(target)` | 等于目标值 |
| `Ne(target)` | 不等于目标值 |
| `In(values...)` | 命中合法值列表(构造时哈希为集合, O(1) 查找) |
| `NotIn(values...)` | 不命中非法值列表 |

**字符串**

| 规则 | 说明 |
|------|------|
| `Match(pattern)` | 匹配正则(构造时编译, 非法表达式立即 panic) |
| `Date(layout)` | 符合 time.Parse 布局, 如 "2006-01-02" |
| `Contains(substr)` | 包含指定子串 |
| `ContainsAny(chars)` | 包含任一指定字符(密码复杂度等场景) |
| `Excludes(substr)` | 不包含指定子串 |
| `StartWith(prefix)` | 以指定前缀开头 |
| `EndWith(suffix)` | 以指定后缀结尾 |
| `TrimSpace(inner)` | 先去首尾空白再套用 inner, 容忍输入前后空格 |

**指针**

| 规则 | 说明 |
|------|------|
| `NotNil` | 指针非 nil(用于 Validate 直接校验指针, 或集合元素为指针时逐元素校验) |

### is 包 — 格式断言(71 条)

全部作用于字符串, 严格模式(空串一律不合法), 可选字段用 `rules.Optional` 包装。

**邮箱**

| 规则 | 说明 |
|------|------|
| `Email` | 格式合法且域名存在 MX 记录(需 DNS 查询) |
| `EmailFormat` | 仅格式合法, 不做 DNS 查询 |

**字符集**

| 规则 | 说明 |
|------|------|
| `Alpha` / `Alphanumeric` | 仅英文字母 / 仅英文字母与数字 |
| `Digit` | 仅 ASCII 数字 |
| `UTFLetter` / `UTFDigit` / `UTFNumeric` / `UTFLetterNumeric` | unicode 字母 / 十进制数字 / 数值字符 / 字母或数值 |
| `LowerCase` / `UpperCase` | 全小写 / 全大写 |
| `ASCII` / `PrintableASCII` | 仅 ASCII / 仅可打印 ASCII |
| `Hexadecimal` | 十六进制数字序列 |

**数值与布尔**

| 规则 | 说明 |
|------|------|
| `Int` | 整数字面量(允许正负号) |
| `Float` | 十进制浮点字面量(允许科学计数法) |
| `Boolean` | true / false / t / f / 1 / 0, 大小写不敏感 |

**颜色**

| 规则 | 说明 |
|------|------|
| `HexColor` | #RGB 或 #RRGGBB |
| `RGBColor` | rgb(r, g, b), 分量 0-255 |
| `RGBA` | rgba(r, g, b, a), a 为 0-1 小数或 0-100% |
| `HSL` / `HSLA` | hsl / hsla, h 0-360, s/l 0-100 |

**标识符**

| 规则 | 说明 |
|------|------|
| `UUID` / `UUIDv3` / `UUIDv4` / `UUIDv5` / `UUIDv7` | 通用 / 指定版本 UUID |
| `ULID` | 26 字符 Crockford Base32 |
| `MongoID` | 24 位十六进制 MongoDB ObjectID |

**文本**

| 规则 | 说明 |
|------|------|
| `JSON` | 语法完整的 JSON 文档 |
| `Base64` | 标准 base64(RFC 4648, 含填充) |
| `Base64URL` | URL 安全 base64(不含填充) |
| `JWT` | 三段式 JWT(各段均可 base64url 解码) |
| `DataURI` | RFC 2397 data URI(含 base64 载荷解码校验) |

**网络**

| 规则 | 说明 |
|------|------|
| `URL` | 绝对 URL(协议 + 主机) |
| `RequestURL` / `RequestURI` | 请求 URL / 请求 URI |
| `IP` / `IPv4` / `IPv6` | IP 地址 |
| `MAC` | MAC 地址(冒号 / 连字符 / 点分隔) |
| `Host` | IP 地址或 DNS 名称 |
| `Port` | 端口 1-65535 |
| `DialString` | host:port 拨号串 |
| `E164` | E.164 国际电话号码 |
| `CIDR` / `CIDRv4` / `CIDRv6` | CIDR 网段 |
| `Subdomain` | 单个子域标签(1-63 字符) |
| `Domain` | 完整域名(至少两段) |
| `DNSName` | DNS 名称(允许末尾根点) |

**地理**

| 规则 | 说明 |
|------|------|
| `Latitude` | 纬度 [-90, 90] |
| `Longitude` | 经度 [-180, 180] |

**卡 / 书号**

| 规则 | 说明 |
|------|------|
| `CreditCard` | 通过 Luhn 校验的卡号 |
| `ISBN10` / `ISBN13` / `ISBN` | ISBN-10 / ISBN-13 / 两者兼容 |

**业务格式**

| 规则 | 说明 |
|------|------|
| `SSN` | 美国社保号 |
| `Semver` | 语义化版本 2.0.0 |
| `Origin` | CORS 来源地址(http(s) + 主机 + 可选端口) |
| `TimeZone` | 系统可识别的时区名称, 如 Asia/Shanghai |

**字符宽度**

| 规则 | 说明 |
|------|------|
| `Multibyte` | 至少含一个多字节字符 |
| `FullWidth` / `HalfWidth` / `VariableWidth` | 含全角 / 含半角 / 两者兼具 |

**中文本地化**

| 规则 | 说明 |
|------|------|
| `MobileCN` | 中国大陆手机号(11 位, 1 开头, 第二位 3-9) |
| `IDCardCN` | 18 位身份证号(校验出生日期与 GB 11643 校验码) |
| `PostalCodeCN` | 6 位邮政编码 |

**国际标准码**

| 规则 | 说明 |
|------|------|
| `CountryCode` | ISO 3166-1 alpha-2 国家代码(完整集合) |
| `CurrencyCode` | ISO 4217 货币代码(完整集合) |
| `LanguageCode` | ISO 639-1 语言代码(完整集合) |

### collections 包 — 集合验证

| 规则 | 说明 |
|------|------|
| `Slice(rules...)` | 切片逐元素校验, 失败按下标聚合 |
| `Each(rules...)` | iter.Seq 序列逐元素校验 |
| `Map(rules...)` | map 值逐元素校验, 失败按键名聚合 |
| `SliceLen(min, max)` | 切片长度落在闭区间 |
| `MapLen(min, max)` | map 键值对数量落在闭区间 |
| `Unique()` | 元素互不重复 |

## 使用注意事项

以下行为为设计使然, 使用前请留意, 避免踩坑:

- **nil 指针字段被直接跳过**: 绑定到 `FieldBuilder` 的指针字段值为 nil 时,
  不会执行任何规则, **包括 `Required`**。必填语义请绑定非指针字段, 或先解引用到值类型。
- **零尺寸字段无法绑定**: `struct{}` 等零尺寸字段因所有实例共享同一地址, 无法可靠定位,
  绑定会触发 panic。请勿对零尺寸字段做校验。
- **`is.Email` 依赖 DNS**: 会查询域名 MX 记录, 离线或单元测试环境请改用 `is.EmailFormat`。
- **`rules.MultipleOf` 经 float64 计算**: 超大整数(超过 2^53)可能损失精度。
- **跨字段有序比较仅支持有序类型**: `WithFieldGt/Gte/Lt/Lte` 对非有序类型(如 struct)会 panic。
- **字段名推导**: `ValidateStruct` 按 json tag 推导对外字段名, 缺省或为 `-` 时回退为蛇形 Go 字段名。
- **`rules.NotNil` 不适用于 nil 指针字段**: 与第一条同理, 指针字段为 nil 时已被跳过,
  `NotNil` 请用于 `Validate` 直接校验指针, 或集合元素为指针时逐元素校验。

## 错误处理

错误码以常量形式提供(如 `rules.CodeLength`), 拼写错误在编译期暴露。
规则失败时通过 `NewMessage` 按当前生效语言渲染模板参数, 生成带错误码的 `*Error`:

```go
// 规则内部(以 Length 为例): 失败时按当前语言渲染
return verax.NewMessage(rules.CodeLength, map[string]string{"min": "5", "max": "100"})

// 应用自定义错误: 直接构造完整消息
err := verax.NewError("app.token.expired", "登录已过期")
```

```go
// 提取结构化信息(Go 1.26 errors.AsType)
if e, ok := errors.AsType[*verax.Error](err); ok {
    fmt.Println(e.Code)               // validation.length
    fmt.Println(e.Message)            // the length must be between 5 and 100
    fmt.Println(e.Field)              // name (经 ValidateStruct 校验时填充)
}
```

渲染顺序: 生效语言表 → 内置英文表 → 错误码本身。
聚合错误的文本形态为 `"name: cannot be blank; age: must be no less than 18"`。

## 国际化

内置 en / zh-CN / zh-TW / ja / fr / de 六种语言, 一行注入:

```go
verax.RegisterZhCN()  // 注册简体中文并立即成为当前生效语言
// 之后校验失败产生的错误自动为中文: "name: 不能为空"
```

- 内置语言: `RegisterEn` / `RegisterZhCN` / `RegisterZhTW` / `RegisterJa` / `RegisterFr` / `RegisterDe`
- 注册即生效: `RegisterLocale(locale, table)` 直接替换当前生效语言表
- 语言表模板保留 `{{.占位符}}`(如 `"不能小于 {{.min}}"`), 由规则参数在渲染时插值,
  所有语言都能拿到精确数值
- 自定义语言可调用 `verax.RegisterLocale(locale, messages)` 扩展,
  用 `verax.Codes()` 取得全部错误码清单校验覆盖完整

## 与其他库对比

| 维度 | verax | go-playground/validator | ozzo-validation |
|------|-------|------------------------|-----------------|
| 声明方式 | 类型安全的链式 API | struct tag | 链式 API |
| 类型安全 | 规则与值类型编译期绑定 | 无(反射解析字符串 tag) | 有限泛型 |
| 反射开销 | 仅字段定位 | 每次校验全反射 | 较少 |
| 可选语义 | `Optional` 显式声明 | `omitempty` tag | `NilOrNotEmpty` |
| 跨字段比较 | `WithFieldEq/Gte` 等内置 | `eqfield` 等 tag | 需手写 |
| 条件必填 | `RequiredIf` | `required_if` 等 tag | 需手写 |
| 集合校验 | 逐元素 + 唯一 + 长度 | 需第三方或手写 | `Each` |
| 国际化 | 内置 6 种语言 | 需第三方库 | 需第三方 |
| 依赖 | 零第三方依赖 | 零 | 零 |

**为什么不用 struct tag?** tag 方案把校验规则写进字符串, 靠反射在运行期解析,
拼写错误只能在运行时暴露, 且无法表达泛型规则与复杂条件。
verax 让每条规则都是 `func(T) error`, 类型由编译器保证, 规则可自由组合、复用与单测。

## 版本与稳定性

- 遵循 [语义化版本](https://semver.org/lang/zh-CN/) SemVer: `X.Y.Z` 分别对应不兼容 / 新功能 / 修复。
- **错误码是对外契约**: `internal/codes` 中的错误码值一经发布不可变更, 新增规则只增不减。
- `0.x` 阶段允许不兼容变更(此时不兼容仅递增次版本号)。
- 变更记录见 [CHANGELOG.zh-CN.md](CHANGELOG.zh-CN.md)。

## 贡献

欢迎提交 Issue 与 PR。新增一条规则的完整流程:

1. 在 `internal/codes/codes.go` 定义错误码常量;
2. 在 `internal/messages/` 的 6 种语言表中补翻译(`TestBuiltinLocalesCompleteness` 会校验逐码覆盖);
3. 在对应包(`rules` / `is` / `collections`)的 `codes.go` 转导出错误码;
4. 实现规则并补充测试(通过与失败两条路径)。

本地验证:

```
go test ./...
go vet ./...
```

提交信息遵循 Conventional Commits(如 `feat(rules): 新增 ContainsAny 规则`)。

## License

[MIT](LICENSE)

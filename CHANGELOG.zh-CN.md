# 变更日志

[English](CHANGELOG.md) | **简体中文**

本项目的所有显著变更都记录在此文件, 格式遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/), 版本遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## 0.1.0 - 2026-08-26

### 新增

- 通用规则: `RequiredIf` / `Optional` / `Length` / `Len` / `Min` / `Max` / `Between` / `Gt` / `Lt` / `Eq` / `Ne` / `In` / `NotIn` / `Match` / `Date` / `Contains` / `ContainsAny` / `Excludes` / `StartWith` / `EndWith` / `MultipleOf` / `TrimSpace` / `NotNil`
- 跨字段比较: `WithFieldEq` / `WithFieldNe` / `WithFieldGt` / `WithFieldGte` / `WithFieldLt` / `WithFieldLte`
- 集合验证: `Slice` / `Each` / `Map` / `SliceLen` / `MapLen` / `Unique`
- 格式断言: 邮箱 / 字符集 / 数值 / 颜色 / UUID / ULID / MongoID / 文本(JSON/Base64/Base64URL/JWT/DataURI) / 网络(URL/IP/MAC/CIDR/Domain/Host/Port/E164/DialString) / 地理(Latitude/Longitude) / 卡( CreditCard ) / 书号(ISBN) / 业务格式(SSN/Semver/Origin/TimeZone/Boolean)
- 中文本地化: `MobileCN` / `IDCardCN` / `PostalCodeCN`
- 国际标准码: `CountryCode` / `CurrencyCode` / `LanguageCode`
- 国际化: 内置 en / zh-CN / zh-TW / ja / fr / de 六种语言, 全部错误码逐码覆盖

### 修复

- `is.Latitude` / `is.Longitude` 放行 `"NaN"` 输入
- `TestMember` 测试数据与断言不一致
- `WithLabel` 对 `WithCheckFn` 返回的裸错误不生效
- 零尺寸字段绑定给出明确 panic 提示, 而非误导性报错
- 移除死导出 `collections.ErrLen`

### 优化

- `collections.Map` 复用 `collectErrors`, 消除重复实现
# Changelog

**English** | [简体中文](CHANGELOG.zh-CN.md)

All notable changes to this project are documented in this file. The format follows [Keep a Changelog](https://keepachangelog.com/), and this project adheres to [Semantic Versioning](https://semver.org/).

## 0.1.0 - 2026-08-26

### Added

- General rules: `RequiredIf` / `Optional` / `Length` / `Len` / `Min` / `Max` / `Between` / `Gt` / `Lt` / `Eq` / `Ne` / `In` / `NotIn` / `Match` / `Date` / `Contains` / `ContainsAny` / `Excludes` / `StartWith` / `EndWith` / `MultipleOf` / `TrimSpace` / `NotNil`
- Cross-field comparison: `WithFieldEq` / `WithFieldNe` / `WithFieldGt` / `WithFieldGte` / `WithFieldLt` / `WithFieldLte`
- Collection validation: `Slice` / `Each` / `Map` / `SliceLen` / `MapLen` / `Unique`
- Format assertions: email / character sets / numeric / color / UUID / ULID / MongoID / text (JSON/Base64/Base64URL/JWT/DataURI) / network (URL/IP/MAC/CIDR/Domain/Host/Port/E164/DialString) / geography (Latitude/Longitude) / cards (CreditCard) / book numbers (ISBN) / business formats (SSN/Semver/Origin/TimeZone/Boolean)
- China localization: `MobileCN` / `IDCardCN` / `PostalCodeCN`
- International standard codes: `CountryCode` / `CurrencyCode` / `LanguageCode`
- Internationalization: six built-in languages en / zh-CN / zh-TW / ja / fr / de, with full per-code coverage of all error codes

### Fixed

- `is.Latitude` / `is.Longitude` accepted `"NaN"` input
- `TestMember` test data inconsistent with its assertion
- `WithLabel` had no effect on bare errors returned by `WithCheckFn`
- zero-size field binding now gives an explicit panic hint instead of a misleading error
- removed the dead export `collections.ErrLen`

### Optimized

- `collections.Map` reuses `collectErrors`, removing duplicated logic
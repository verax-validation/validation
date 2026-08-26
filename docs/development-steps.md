# verax 开发步骤

基于 ozzo-validation 的设计经验, 用 Go 1.26 特性全新实现的泛型数据验证库。

- module 路径: `github.com/verax-validation/validation`
- API 策略: 全泛型 API (`Rule[T]`, 编译期类型安全)
- 流程约束: 每完成一个步骤暂停, 经人工审核通过后才能进入下一步骤

## 核心技术决策备忘

1. 规则即函数: `type Rule[T any] func(T) error`
   - 依赖 Go 1.21+ 泛型函数实参推断, `Validate(name, verax.Required)` 处编译器自动实例化 T=string
2. struct 验证零反射: `Field[T](&u.Name, rules...)` 绑定字段指针, 闭包直接解引用取值,
   替代 ozzo 的 reflect.Value 遍历
3. 数值规则统一: `Min/Max/Between/MultipleOf[T cmp.Ordered]` 一套实现覆盖所有有序类型,
   消除 ozzo 中 int64/float64/string 三套重复代码
4. 集合验证基于 `iter.Seq[T]`(Go 1.23+), 遍历零反射
5. 错误提取对齐 Go 1.26: 配套提供 `errors.AsType[verax.Error]` 使用范式
6. 错误码机制保留(ozzo 特色), 支持 MessageMap 多语言翻译
7. 空值语义采用显式 Optional (2026-08-25 审核确认):
   - 默认严格模式, 规则无条件执行, 行为完全由代码表达
   - rules.Optional(rules...) 包装器声明可选字段, 零值短路放行
   - 替代 ozzo 每条规则内部的隐式 IsEmpty 放行, 消除隐藏魔法与重复判断
8. 国际化双层设计 (2026-08-25 审核确认):
   - 进程级: RegisterLocale(locale, MessageMap) 登记语言包, SetLocale 设定生效语言,
     Error() 默认按生效语言渲染
   - 请求级: (*Error).Localized / (Errors).Localized 按指定语言渲染,
     验证链路零改动, 适配多租户网关场景
   - 并发模型: localeStore 写时复制 + atomic.Pointer, 读路径无锁
   - 回退链: 指定语言未命中错误码时, 回退到错误自带英文默认模板
9. 内置语言包 locales 子包 (2026-08-25 审核确认, 调整原"主库不内置语言包"决策):
   - RegisterZhCN/Ja/Fr/De/En 注入对应翻译数据(不切换语言), 五语言 API 对称
   - NewError 自动收集错误码, Codes() 导出全集, DefaultMessages() 导出英文基准
   - locales 完整性测试强制各语言表覆盖全部错误码, 防止新增规则漏翻
10. 错误文案集中化 (2026-08-25 审核确认):
   - 错误码常量统一定义于 internal/codes, 各公开包转导出(如 rules.CodeRequired),
     翻译表 key 全部引用常量, 拼写错误升级为编译期错误
   - NewError 改为单参数构造, 英文默认文案集中存放于内置模板表,
     Error.Message 由构造函数自动填充, 规则文件零文案
   - 消息渲染回退链: 生效语言表 -> 内置英文 -> 错误码本身
11. 语言数据下沉与注册 API 收敛 (2026-08-25 审核确认):
   - 全部语言数据(含 en)集中于 internal/messages, 为唯一事实来源;
     该层仅依赖 internal/codes, 杜绝循环导入
   - Register 系列函数由根包提供(RegisterZhCN/ZhTW/Ja/Fr/De/En)
   - init 默认注入英文并激活 LocaleEn, 直接构造初始 store,
     不经持锁的 RegisterLocale(非可重入锁, 嵌套调用会死锁)

## 目录结构规划

```
verax/
├── docs/                    # 开发文档
├── internal/
│   └── zero/                # 零值判断内部实现
├── is/                      # 高频格式断言糖 (Email/URL/UUID...)
├── rules/                   # 内置通用值规则
├── collections/             # 集合验证 (Slice/Map/Each)
├── errors.go                # Error/Errors/MessageMap
├── rule.go                  # Rule[T]/组合器
├── validate.go              # Validate/Field/ValidateStruct/Skip/When
├── validatable.go           # Validatable 接口与 Valid 规则
├── go.mod
└── README.md
```

## 步骤清单

### 步骤 1: 项目脚手架

- 目标: 清理模板代码, 固化模块身份与目录骨架
- 产出:
  - go.mod 改为 `github.com/verax-validation/validation`
  - 删除 GoLand 模板 main.go (库项目不含入口)
  - 创建 rules/ collections/ is/ internal/zero/ 目录占位
  - LICENSE (MIT), .gitignore
- 验证: `go build ./...` 通过, 目录树与规划一致

### 步骤 2: 核心错误体系 (errors.go)

- 目标: 结构化错误, 兼顾单值错误与字段聚合错误
- 产出: errors.go 及其单测
  - `Error{Code, Message}` 实现 error 接口
  - `Errors map[string]error` 字段级错误集合
  - 默认消息注册表, 支持运行期覆盖翻译
- 验证: `go test ./...` 通过; 错误消息格式与设计一致

### 步骤 3: 核心验证引擎 (rule.go + validate.go)

- 目标: 泛型验证主链路
- 产出:
  - rule.go: `Rule[T]`, 组合器 `All/When`
  - validate.go: `Validate[T]`/`Field[T]`/`ValidateStruct[S]`/`Skip[T]`
  - validatable.go: `Validatable` 接口与 `Valid` 规则
- 验证: 单测覆盖成功/短路/嵌套 struct/nil 指针场景

### 步骤 4: rules 子包 (内置通用规则)

- 目标: 对齐 ozzo 内置规则的常用集合
- 产出: required/length/numeric/membership/pattern/date 各文件及单测
- 验证: `go test ./...` 通过; 与 ozzo 同名规则行为语义一致

### 步骤 5: is 子包 (格式断言)

- 目标: Email/URL/UUID/IP/JSON 等高频格式快捷规则
- 产出: is 包及单测
- 验证: `go test ./...` 通过

### 步骤 6: collections 子包 (集合验证)

- 目标: 切片/map/迭代器的元素级验证
- 产出: slice.go/map.go/each.go 及单测
- 验证: 嵌套结构 (slice of struct) 验证正确

### 步骤 7: 基准测试对比

- 目标: 量化相对 ozzo 的性能收益
- 产出: bench_test.go, 覆盖简单值/struct/切片三场景
- 验证: 产出对比报告至 docs/benchmark.md, 无性能回归

### 步骤 8: README 与 CI

- 目标: 正式开源形象
- 产出: README.md (含用法示例), .github/workflows/ci.yml
- 验证: 示例代码可直接编译运行

### 步骤 9: 全量质量验证

- 目标: 收尾把关
- 产出: 无新增代码, 仅修复问题
- 验证: `go vet ./...`, `golangci-lint run`, `go test ./...` 全绿;
  codegraph sync 后 impact 抽查关键符号影响面

## 审核记录

| 步骤 | 完成时间 | 审核结论 |
|------|---------|---------|
| 1 | 2026-08-25 | 通过 |
| 2 | 2026-08-25 | 通过 |
| 3 | 2026-08-25 | 通过 |
| 4 | 2026-08-25 | 通过(含 i18n 双层增强) |
| 5 | 2026-08-25 | 通过(55 条规则全量对齐, 清除第三方引用, 补齐注释) |
| 6 | 2026-08-25 | 通过(断言统一改为 errors.AsType) |
| 7 | 2026-08-25 | 通过(bench 独立子模块, 切片场景 10 到 15 倍领先) |
| 8 | | 待审核(README 示例已通过编译运行验证) |
| 9 | | |

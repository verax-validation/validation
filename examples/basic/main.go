// Command basic 演示 verax 的完整用法: 校验简单值 / struct / 跨字段 / 集合 / 国际化 / 错误处理。
package main

import (
	"errors"
	"fmt"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/collections"
	"github.com/verax-validation/validation/is"
	"github.com/verax-validation/validation/rules"
)

// RegisterForm 是注册表单模型。
type RegisterForm struct {
	Name     string
	Email    string
	Age      int
	Password string
	Confirm  string
	Tags     []string
}

func main() {
	// 注入简体中文, 之后校验失败的错误自动为中文
	verax.RegisterZhCN()

	form := &RegisterForm{
		Name:     "alice",
		Email:    "alice@example.com",
		Age:      25,
		Password: "secret#123",
		Confirm:  "secret#123",
		Tags:     []string{"go", "web"},
	}

	// 路径 1: 全部字段合法, 校验通过
	if err := validate(form); err != nil {
		fmt.Println("合法表单不应失败:", err)
	} else {
		fmt.Println("1) 合法表单: 校验通过")
	}

	// 路径 2: 密码不含特殊字符, 且确认密码不一致, 校验失败
	form.Password = "secret123"
	form.Confirm = "secret456"
	if err := validate(form); err != nil {
		fmt.Println("2) 非法表单:")
		fmt.Println("   ", err)
		// 错误处理: 聚合错误按字段提取, 得到结构化错误码
		if errs, ok := errors.AsType[verax.Errors](err); ok {
			if e, found := errs.Get("password"); found {
				if ve, ok2 := errors.AsType[*verax.Error](e); ok2 {
					fmt.Printf("    错误码: %s (字段 %s)\n", ve.Code, ve.Field)
				}
			}
		}
	} else {
		fmt.Println("非法表单不应通过")
	}
}

func validate(form *RegisterForm) error {
	return verax.ValidateStruct(form,
		// 必填 + 长度区间 + 自定义标签
		verax.Field[string]().WithField(&form.Name).
			WithRules(rules.Required[string], rules.Length[string](2, 32)).
			WithLabel("姓名"),
		// 格式断言 + 自定义标签
		verax.Field[string]().WithField(&form.Email).
			WithRules(rules.Required[string], is.EmailFormat).
			WithLabel("邮箱"),
		// 数值区间
		verax.Field[int]().WithField(&form.Age).
			WithRules(rules.Between(18, 120)).
			WithLabel("年龄"),
		// 长度 + 密码复杂度
		verax.Field[string]().WithField(&form.Password).
			WithRules(rules.Length[string](8, 64), rules.ContainsAny("!@#$")).
			WithLabel("密码"),
		// 跨字段比较: 确认密码必须与密码一致
		verax.Field[string]().WithField(&form.Confirm).
			WithFieldEq(&form.Password).
			WithLabel("确认密码"),
		// 集合: 数量区间 + 元素唯一 + 逐元素非空
		verax.Field[[]string]().WithField(&form.Tags).
			WithRules(
				collections.SliceLen[string](1, 5),
				collections.Unique[string](),
				collections.Slice[string](rules.Required),
			).
			WithLabel("标签"),
	)
}

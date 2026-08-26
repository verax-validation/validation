module github.com/verax-validation/validation/bench

go 1.26

require (
	github.com/go-ozzo/ozzo-validation/v4 v4.0.0
	github.com/verax-validation/validation v0.0.0
)

require github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect

replace github.com/verax-validation/validation => ../

replace github.com/go-ozzo/ozzo-validation/v4 => ../../ozzo-validation

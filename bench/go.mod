module github.com/verax-validation/validation/bench

go 1.26

require (
	github.com/go-ozzo/ozzo-validation/v4 v4.0.0
	github.com/verax-validation/validation v0.0.0
)

require gopkg.in/asaskevich/govalidator.v9 v9.0.0-20180315120708-ccb8e960c48f // indirect

replace github.com/verax-validation/validation => ../

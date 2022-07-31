//go:generate go run github.com/mailru/easyjson/easyjson -no_std_marshalers -lower_camel_case $GOFILE
package util

const (
	ResultOK = Result(iota)

	ResultUserAlreadyExists
	ResultWrongLoginOrPassword

	ResultWrongOldPassword

	ResultBondNotFound
)

type Result int

//easyjson:json
type BaseResult struct {
	Status Result
}

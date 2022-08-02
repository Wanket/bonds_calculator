//go:generate go run github.com/mailru/easyjson/easyjson -no_std_marshalers -lower_camel_case $GOFILE
package util

const (
	StatusOK = Status(iota)

	StatusUserAlreadyExists
	StatusWrongLoginOrPassword

	StatusWrongOldPassword

	StatusBondNotFound
)

type Status int

//easyjson:json
type BaseResult struct {
	Status Status
}

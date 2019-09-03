package transaction

import "fmt"

const (
	errBeanFuncType int8 = iota
	errBeanFuncArgsNum
	errBeanFuncArgsMatch
)

var errMessage = map[int8]string{
	errBeanFuncType:      "bean func not a func",
	errBeanFuncArgsNum:   "bean func args num not match",
	errBeanFuncArgsMatch: "bean func args type not match",
}

type transactionError struct {
	code   int8
	detail string
}

func (e *transactionError) Error() string {
	return fmt.Sprintf("[%s][%s]", errMessage[e.code], e.detail)
}

func newError(code int8, detail string) *transactionError {
	return &transactionError{code: code, detail: detail}
}

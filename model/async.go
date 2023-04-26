package model

type AsyncTaskRes struct {
	Idx int
	Res interface{}
	Err error
}

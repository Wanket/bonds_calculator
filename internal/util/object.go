package util

func CopyObjectByRef[T any](obj *T) *T {
	copiedObj := *obj

	return &copiedObj
}

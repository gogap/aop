package aop

type MethodMetadata struct {
	method     interface{}
	MethodName string
	IsStatic   bool
	File       string
	Line       int
}

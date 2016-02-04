package aop

import (
	"reflect"
)

type MethodMetadata struct {
	Method reflect.Method
	File   string
	Line   int
}

func (p *MethodMetadata) IsEqual(t reflect.Type) bool {
	if t.ConvertibleTo(p.Method.Type) {
		return false
	}

	baseIndex := 0
	if p.Method.Index >= 0 {
		baseIndex = 1
	}

	if t.NumIn()+baseIndex != p.Method.Type.NumIn() {
		return false
	}

	for i := 0; i < p.Method.Type.NumIn()-baseIndex; i++ {
		if p.Method.Type.In(baseIndex+i) != t.In(i) {
			return false
		}
	}

	for i := 0; i < p.Method.Type.NumOut(); i++ {
		if p.Method.Type.Out(baseIndex+i) != t.Out(i) {
			return false
		}
	}

	return true
}

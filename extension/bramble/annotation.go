package bramble

import "entgo.io/ent/entc"

const ConnectionAnnotationName = `BrambleConnection`

func ConnectionAnnotation() entc.Annotation {
	return &connectionAnnotation{}
}

type connectionAnnotation struct {
	entc.Annotation
}

func (connectionAnnotation) Name() string {
	return `BrambleConnection`
}

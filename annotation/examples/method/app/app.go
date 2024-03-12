package app

// barStruct Lorem ipsum dolor sit amet, consectetur adipiscing elit
// @MyMethodAnnotation(param=xpto)
type barStruct struct {
}

// barMethod Lorem ipsum dolor sit amet, consectetur adipiscing elit
// @MyMethodAnnotation(param=xpto)
func (s *barStruct) barMethod(r string) (xpto *FooStruct) {
	return &FooStruct{}
}

// FooStruct Lorem ipsum dolor sit amet, consectetur adipiscing elit
// @MyMethodAnnotation(param=xpto)
type FooStruct struct {
}

// FooMethod Lorem ipsum dolor sit amet, consectetur adipiscing elit
// @MyMethodAnnotation(param=xpto)
func (s *FooStruct) FooMethod(r string) (xpto *string) {
	x := ""
	return &x
}

// FooFunc Lorem ipsum dolor sit amet, consectetur adipiscing elit
// @MyMethodAnnotation(param=xpto)
func FooFunc(x string) (string, error) {
	return "", nil
}

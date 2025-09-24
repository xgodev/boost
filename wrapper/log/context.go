package log

import "context"

func ToContext(ctx context.Context) context.Context {
	return l.ToContext(ctx)
}

// FromContext calls concrete Logger.FromContext().
func FromContext(ctx context.Context) Logger {
	return l.FromContext(ctx)
}

func Ctx(ctx context.Context, obj interface{}) Logger {
	return l.FromContext(ctx).WithTypeOf(obj)
}

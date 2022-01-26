// Copyright (c) 2021 Jing-Ying Chen. Subject to the MIT License.

package goutil

import "context"

func ContextWithMap(ctx context.Context, maps ...map[string]interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, m := range maps {
		for k, v := range m {
			ctx = context.WithValue(ctx, k, v)
		}
	}
	return ctx
}

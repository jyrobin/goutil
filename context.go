// Copyright (c) 2021 Jing-Ying Chen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func ContextWithKVs(ctx context.Context, kvs ...KV) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	for _, kv := range kvs {
		ctx = context.WithValue(ctx, kv.Key, kv.Value)
	}
	return ctx
}

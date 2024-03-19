// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package util

import "reflect"

func Invoke(fn interface{}, types []interface{}, params ...interface{}) {
	v := reflect.ValueOf(fn)
	if v.IsValid() && !v.IsNil() {
		in := make([]reflect.Value, len(params))
		for k, param := range params {
			if param != nil {
				in[k] = reflect.ValueOf(param)
			} else {
				in[k] = reflect.Zero(reflect.TypeOf(types[k]).Elem())
			}
		}
		v.Call(in)
	}
}

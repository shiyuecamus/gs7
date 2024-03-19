// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package util

import "time"

func StrOrDefault(str string, defVal string) string {
	if str == "" {
		return defVal
	}
	return str
}

func IntOrDefault(val int, defVal int) int {
	if val == 0 {
		return defVal
	}
	return val
}

func DurationOrDefault(val time.Duration, defVal time.Duration) time.Duration {
	if val == 0 {
		return defVal
	}
	return val
}

func AnyOrDefault(val any, defVal any) any {
	if val == nil {
		return defVal
	}
	return val
}

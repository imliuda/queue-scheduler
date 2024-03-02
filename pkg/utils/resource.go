/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import corev1 "k8s.io/api/core/v1"

func LessEqual(r1, r2 corev1.ResourceList, skipUndefined bool) bool {
	if r1 == nil {
		return true
	}

	for r1Name, r1Value := range r1 {
		if r2Value, ok := r2[r1Name]; ok {
			if r1Value.Cmp(r2Value) <= 0 {
				continue
			} else {
				return false
			}
		} else if skipUndefined {
			continue
		} else {
			return false
		}
	}
	return true
}

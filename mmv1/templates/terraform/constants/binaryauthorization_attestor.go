<%# The license inside this block applies to this file.
	# Copyright 2020 Google Inc.
	# Licensed under the Apache License, Version 2.0 (the "License");
	# you may not use this file except in compliance with the License.
	# You may obtain a copy of the License at
	#
	#     http://www.apache.org/licenses/LICENSE-2.0
	#
	# Unless required by applicable law or agreed to in writing, software
	# distributed under the License is distributed on an "AS IS" BASIS,
	# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	# See the License for the specific language governing permissions and
	# limitations under the License.
-%>
func CompareSignatureAlgorithm(_, old, new string, _ *schema.ResourceData) bool {
	// See https://cloud.google.com/binary-authorization/docs/reference/rest/v1/projects.attestors#signaturealgorithm
	normalizedAlgorithms := map[string]string{
		"ECDSA_P256_SHA256": "ECDSA_P256_SHA256",
		"EC_SIGN_P256_SHA256": "ECDSA_P256_SHA256",
		"ECDSA_P384_SHA384": "ECDSA_P384_SHA384",
		"EC_SIGN_P384_SHA384": "ECDSA_P384_SHA384",
		"ECDSA_P521_SHA512": "ECDSA_P521_SHA512",
		"EC_SIGN_P521_SHA512": "ECDSA_P521_SHA512",
	}

	normalizedOld := old
	normalizedNew := new

	if normalized, ok := normalizedAlgorithms[old]; ok {
		normalizedOld = normalized
	}
	if normalized, ok := normalizedAlgorithms[new]; ok {
		normalizedNew = normalized
	}

	if normalizedNew == normalizedOld {
		return true
	}

	return false
}

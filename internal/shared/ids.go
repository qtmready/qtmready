// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package shared

import (
	"crypto/rand"
	"math"
)

const (
	size = 12 // number of characters in the generated ID.
)

var (
	chars = []rune("abcdefghijklnopqrstvxyz2345689") // character set for encoding.
)

// bitmask calculates the optimal bitmask for encoding a given alphabet size.
//
// For example, for a limit of 30, the optimal bitmask is 31, calculated as 2^5 - 1. This bitmask allows us to encode
// 30 distinct characters (0-29) using 5 bits per character.
//
// It iterates from 1 to 8 bits, calculating the bitmask (2^i - 1) for each iteration. If the calculated mask is greater
// than or equal to the alphabet size minus one, it returns the mask. This ensures that the mask covers the entire range
// of the alphabet, allowing for efficient encoding without overflowing the character set.
func bitmask(limit int) int {
	for i := 1; i <= 8; i++ {
		mask := (2 << uint(i)) - 1 // nolint
		if mask >= limit-1 {
			return mask
		}
	}

	return 0
}

// encode generates a 12-character string from a byte slice using base-33 encoding.  For each character in the generated
// ID, random bytes are generated and used for encoding, increasing entropy and ensuring that each character is a valid
// character from the predefined character set. This prevents errors or inconsistencies.
//
// The function operates by iterating over the provided byte slice, extracting bits using the provided mask, and mapping
// the extracted value to a character from the character set. This process continues for the specified number of steps,
// generating a 12-character string.
func encode(bytes []byte, mask, steps int) string {
	id := make([]rune, size)
	done := false

	for !done {
		idx := 0
		_, _ = rand.Read(bytes)

		for step := 0; step < steps; step++ {
			current := bytes[idx] & byte(mask)

			if current < byte(len(chars)) {
				id[idx] = chars[current]
				idx++
			}

			if idx == size {
				done = true
				break
			}
		}
	}

	return string(id[:size])
}

// Idempotent generates a unique 12-character ID using base-30 encoding, meaning a namespace of 30^12, or approximately
// 531.441 quadrillion unique IDs, providing a high degree of collision resistance.
//
// This approach is sufficient for the current use case, where IDs are ephemeral and only a few hundred will be used
// concurrently. However, for scenarios where IDs are persistent or the number of concurrent IDs is expected to grow
// significantly, a more robust collision-resistant algorithm may be required. This is because the birthday paradox
// dictates that the probability of collisions increases with the number of IDs generated, even if not all are in use
// simultaneously.
//
// Future improvements could involve using a larger character set, encoding more bytes, or implementing a more robust
// collision-resistant algorithm.
func Idempotent() string {
	mask := bitmask(len(chars))                            // 31
	factor := 2 * float64(mask*size) / float64(len(chars)) // 2 * 31 * 12 / 30 = 24.8
	steps := int(math.Ceil(factor))                        // 25
	bytes := make([]byte, steps)

	_, _ = rand.Read(bytes)

	return encode(bytes, mask, steps)
}

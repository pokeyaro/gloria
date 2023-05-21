// Copyright (c) 2023 Pokeya Boa <pokeya.mystic@gmail.com>, All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package gloria

import (
	"fmt"
	"time"
)

// DecoratorTimer a decorator for timing functions to test api performance.
func DecoratorTimer(fn func() error) {
	// Start
	startTime := time.Now()

	// Execute
	if e := fn(); e != nil {
		panic(e.Error())
	}

	// Finish
	duration := time.Since(startTime)

	// Output
	fmt.Printf("api request duration: %v\n", duration)
}

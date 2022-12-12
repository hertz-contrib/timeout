// The MIT License (MIT)
//
// Copyright (c) 2022 Friend Chen
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// This file may have been modified by CloudWeGo authors. All CloudWeGo
// Modifications are Copyright 2022 CloudWeGo Authors.

package timeout

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
)

// New implementation of timeout middleware.
// To use this method, you need to add the following code to the Handler you defined to listen for the context timeout.
// select { case <-ctx.Done(): _ = c.Error(context.DeadlineExceeded) return }
func New(opts ...Option) app.HandlerFunc {
	opt := NewOptions(opts...)
	return func(ctx context.Context, c *app.RequestContext) {
		timeoutContext, cancel := context.WithTimeout(ctx, opt.Timing)
		defer cancel()
		c.Next(timeoutContext)
		if errorChain := c.Errors; errorChain != nil && opt.TimeoutHandler != nil {
			for i := range errorChain {
				if errors.Is(context.DeadlineExceeded, errorChain[i].Err) || errors.Is(opt.TErr, errorChain[i].Err) {
					opt.TimeoutHandler(ctx, c)
					return
				}
			}
			return
		} else {
			return
		}
	}
}

// Copyright (c) 2023 Pokeya Boa <pokeya.mystic@gmail.com>, All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package gloria

/*
	A then-catch-finally statement callback style like javascript axios library
*/

// type CallbackOk[T any] func(data *RESTFulResp[T])
type CallbackOk[T any] func(data T)
type CallbackErr func(e *Exception)
type CallbackExtra[T any] func(c *Client[T])

// Then sets a callback function to be executed when the HTTP request is successful.
// The provided callback function cb is invoked only if no exception occurred during the request.
// The cb function is called with the result of the request as its argument.
// After executing the callback function, the client instance is returned.
func (c *Client[T]) Then(cb CallbackOk[T]) *Client[T] {
	if isEmpty(c.Exception.PanicError) && isEmpty(c.Exception.FailureReason) {
		// The default is 0, which can be changed by the WithModifySuccessCode(code int) function.
		if c.Result.Code == c.Config.DefaultOkCode {
			c.ChalkStr(LogLevelSuccess, "HTTP request successful~ 🎉🎉🎉")
		} else {
			c.ChalkStr(LogLevelFail, "The HTTP request was successful, but the business failed, please check!")
		}
		cb(c.Result.Data)
	}

	return c
}

// Catch sets a callback function to be executed when an exception occurs during the HTTP request.
// The provided callback function cb is invoked only if an exception exists in the client instance.
// The cb function is called with the exception object as its argument.
// After executing the callback function, the client instance is returned.
func (c *Client[T]) Catch(cb CallbackErr) *Client[T] {
	if !isEmpty(c.Exception) {
		if !isEmpty(c.Exception.PanicError) {
			c.ChalkPrintf(LogLevelPanic, "Panic Request! 😭")
		}
		if !isEmpty(c.Exception.FailureReason) {
			c.ChalkPrintf(LogLevelFail, "Business Failed! 🥹")
		}
		cb(c.Exception)
	}

	return c
}

// Finally function is a flexible custom callback function with checkpoint functionality.
func (c *Client[T]) Finally(cb CallbackExtra[T], printLog ...bool) {
	if len(printLog) == 1 && printLog[0] {
		c.
			ChalkStr(LogLevelInfo, c.Meta.Method).
			ChalkStr(LogLevelInfo, c.Meta.Url).
			ChalkObj(LogLevelInfo, c.Config).
			ChalkInt(LogLevelInfo, c.Context.Response.Status).
			ChalkPrintf(LogLevelInfo, "HTTP Status Code: %d, Business Error Code: %d", c.Context.Response.Status, c.Result.Code)
	}
	cb(c)
}

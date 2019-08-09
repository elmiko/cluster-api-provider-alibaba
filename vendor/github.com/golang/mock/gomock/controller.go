// Copyright 2010 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// GoMock - a mock framework for Go.
//
// Standard usage:
//   (1) Define an interface that you wish to mock.
//         type MyInterface interface {
//           SomeMethod(x int64, y string)
//         }
//   (2) Use mockgen to generate a mock from the interface.
//   (3) Use the mock in a test:
//         func TestMyThing(t *testing.T) {
//           mockCtrl := gomock.NewController(t)
//           defer mockCtrl.Finish()
//
//           mockObj := something.NewMockMyInterface(mockCtrl)
//           mockObj.EXPECT().SomeMethod(4, "blah")
//           // pass mockObj to a real object and play with it.
//         }
//
// By default, expected calls are not enforced to run in any particular order.
// Call order dependency can be enforced by use of InOrder and/or Call.After.
// Call.After can create more varied call order dependencies, but InOrder is
// often more convenient.
//
// The following examples create equivalent call order dependencies.
//
// Example of using Call.After to chain expected call order:
//
//     firstCall := mockObj.EXPECT().SomeMethod(1, "first")
//     secondCall := mockObj.EXPECT().SomeMethod(2, "second").After(firstCall)
//     mockObj.EXPECT().SomeMethod(3, "third").After(secondCall)
//
// Example of using InOrder to declare expected call order:
//
//     gomock.InOrder(
//         mockObj.EXPECT().SomeMethod(1, "first"),
//         mockObj.EXPECT().SomeMethod(2, "second"),
//         mockObj.EXPECT().SomeMethod(3, "third"),
//     )
package gomock

import (
	"fmt"
	"golang.org/x/net/context"
	"reflect"
	"runtime"
	"sync"
)

// A TestReporter is something that can be used to report test failures.
// It is satisfied by the standard library's *testing.T.
type TestReporter interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

<<<<<<< HEAD
// TestHelper is a TestReporter that has the Helper method.  It is satisfied
// by the standard library's *testing.T.
type TestHelper interface {
	TestReporter
	Helper()
}

// cleanuper is used to check if TestHelper also has the `Cleanup` method. A
// common pattern is to pass in a `*testing.T` to
// `NewController(t TestReporter)`. In Go 1.14+, `*testing.T` has a cleanup
// method. This can be utilized to call `Finish()` so the caller of this library
// does not have to.
type cleanuper interface {
	Cleanup(func())
}

// A Controller represents the top-level control of a mock ecosystem.  It
// defines the scope and lifetime of mock objects, as well as their
// expectations.  It is safe to call Controller's methods from multiple
// goroutines. Each test should create a new Controller and invoke Finish via
// defer.
//
//   func TestFoo(t *testing.T) {
//     ctrl := gomock.NewController(t)
//     defer ctrl.Finish()
//     // ..
//   }
//
//   func TestBar(t *testing.T) {
//     t.Run("Sub-Test-1", st) {
//       ctrl := gomock.NewController(st)
//       defer ctrl.Finish()
//       // ..
//     })
//     t.Run("Sub-Test-2", st) {
//       ctrl := gomock.NewController(st)
//       defer ctrl.Finish()
//       // ..
//     })
//   })
=======
// A Controller represents the top-level control of a mock ecosystem.
// It defines the scope and lifetime of mock objects, as well as their expectations.
// It is safe to call Controller's methods from multiple goroutines.
>>>>>>> 79bfea2d (update vendor)
type Controller struct {
	mu            sync.Mutex
	t             TestReporter
	expectedCalls *callSet
	finished      bool
}

<<<<<<< HEAD
// NewController returns a new Controller. It is the preferred way to create a
// Controller.
//
// New in go1.14+, if you are passing a *testing.T into this function you no
// longer need to call ctrl.Finish() in your test methods.
func NewController(t TestReporter) *Controller {
	h, ok := t.(TestHelper)
	if !ok {
		h = &nopTestHelper{t}
	}
	ctrl := &Controller{
		T:             h,
=======
func NewController(t TestReporter) *Controller {
	return &Controller{
		t:             t,
>>>>>>> 79bfea2d (update vendor)
		expectedCalls: newCallSet(),
	}
	if c, ok := isCleanuper(ctrl.T); ok {
		c.Cleanup(func() {
			ctrl.T.Helper()
			ctrl.finish(true, nil)
		})
	}

	return ctrl
}

type cancelReporter struct {
<<<<<<< HEAD
	t      TestHelper
	cancel func()
}

func (r *cancelReporter) Errorf(format string, args ...interface{}) {
	r.t.Errorf(format, args...)
}
func (r *cancelReporter) Fatalf(format string, args ...interface{}) {
	defer r.cancel()
	r.t.Fatalf(format, args...)
}

func (r *cancelReporter) Helper() {
	r.t.Helper()
=======
	t      TestReporter
	cancel func()
}

func (r *cancelReporter) Errorf(format string, args ...interface{}) { r.t.Errorf(format, args...) }
func (r *cancelReporter) Fatalf(format string, args ...interface{}) {
	defer r.cancel()
	r.t.Fatalf(format, args...)
>>>>>>> 79bfea2d (update vendor)
}

// WithContext returns a new Controller and a Context, which is cancelled on any
// fatal failure.
func WithContext(ctx context.Context, t TestReporter) (*Controller, context.Context) {
<<<<<<< HEAD
	h, ok := t.(TestHelper)
	if !ok {
		h = &nopTestHelper{t: t}
	}

	ctx, cancel := context.WithCancel(ctx)
	return NewController(&cancelReporter{t: h, cancel: cancel}), ctx
}

type nopTestHelper struct {
	t TestReporter
}

func (h *nopTestHelper) Errorf(format string, args ...interface{}) {
	h.t.Errorf(format, args...)
}
func (h *nopTestHelper) Fatalf(format string, args ...interface{}) {
	h.t.Fatalf(format, args...)
}

func (h nopTestHelper) Helper() {}

// RecordCall is called by a mock. It should not be called by user code.
=======
	ctx, cancel := context.WithCancel(ctx)
	return NewController(&cancelReporter{t, cancel}), ctx
}

>>>>>>> 79bfea2d (update vendor)
func (ctrl *Controller) RecordCall(receiver interface{}, method string, args ...interface{}) *Call {
	if h, ok := ctrl.t.(testHelper); ok {
		h.Helper()
	}

	recv := reflect.ValueOf(receiver)
	for i := 0; i < recv.Type().NumMethod(); i++ {
		if recv.Type().Method(i).Name == method {
			return ctrl.RecordCallWithMethodType(receiver, method, recv.Method(i).Type(), args...)
		}
	}
	ctrl.t.Fatalf("gomock: failed finding method %s on %T", method, receiver)
	panic("unreachable")
}

func (ctrl *Controller) RecordCallWithMethodType(receiver interface{}, method string, methodType reflect.Type, args ...interface{}) *Call {
	if h, ok := ctrl.t.(testHelper); ok {
		h.Helper()
	}

	call := newCall(ctrl.t, receiver, method, methodType, args...)

	ctrl.mu.Lock()
	defer ctrl.mu.Unlock()
	ctrl.expectedCalls.Add(call)

	return call
}

func (ctrl *Controller) Call(receiver interface{}, method string, args ...interface{}) []interface{} {
	if h, ok := ctrl.t.(testHelper); ok {
		h.Helper()
	}

	// Nest this code so we can use defer to make sure the lock is released.
	actions := func() []func([]interface{}) []interface{} {
		ctrl.mu.Lock()
		defer ctrl.mu.Unlock()

		expected, err := ctrl.expectedCalls.FindMatch(receiver, method, args)
		if err != nil {
<<<<<<< HEAD
			// callerInfo's skip should be updated if the number of calls between the user's test
			// and this line changes, i.e. this code is wrapped in another anonymous function.
			// 0 is us, 1 is controller.Call(), 2 is the generated mock, and 3 is the user's test.
			origin := callerInfo(3)
			ctrl.T.Fatalf("Unexpected call to %T.%v(%v) at %s because: %s", receiver, method, args, origin, err)
=======
			origin := callerInfo(2)
			ctrl.t.Fatalf("Unexpected call to %T.%v(%v) at %s because: %s", receiver, method, args, origin, err)
>>>>>>> 79bfea2d (update vendor)
		}

		// Two things happen here:
		// * the matching call no longer needs to check prerequite calls,
		// * and the prerequite calls are no longer expected, so remove them.
		preReqCalls := expected.dropPrereqs()
		for _, preReqCall := range preReqCalls {
			ctrl.expectedCalls.Remove(preReqCall)
		}

		actions := expected.call()
		if expected.exhausted() {
			ctrl.expectedCalls.Remove(expected)
		}
		return actions
	}()

	var rets []interface{}
	for _, action := range actions {
		if r := action(args); r != nil {
			rets = r
		}
	}

	return rets
}

<<<<<<< HEAD
// Finish checks to see if all the methods that were expected to be called
// were called. It should be invoked for each Controller. It is not idempotent
// and therefore can only be invoked once.
//
// New in go1.14+, if you are passing a *testing.T into NewController function you no
// longer need to call ctrl.Finish() in your test methods.
func (ctrl *Controller) Finish() {
	// If we're currently panicking, probably because this is a deferred call.
	// This must be recovered in the deferred function.
	err := recover()
	ctrl.finish(false, err)
}

func (ctrl *Controller) finish(cleanup bool, panicErr interface{}) {
	ctrl.T.Helper()
=======
func (ctrl *Controller) Finish() {
	if h, ok := ctrl.t.(testHelper); ok {
		h.Helper()
	}
>>>>>>> 79bfea2d (update vendor)

	ctrl.mu.Lock()
	defer ctrl.mu.Unlock()

	if ctrl.finished {
<<<<<<< HEAD
		if _, ok := isCleanuper(ctrl.T); !ok {
			ctrl.T.Fatalf("Controller.Finish was called more than once. It has to be called exactly once.")
		}
		return
=======
		ctrl.t.Fatalf("Controller.Finish was called more than once. It has to be called exactly once.")
>>>>>>> 79bfea2d (update vendor)
	}
	ctrl.finished = true

	// Short-circuit, pass through the panic.
	if panicErr != nil {
		panic(panicErr)
	}

	// Check that all remaining expected calls are satisfied.
	failures := ctrl.expectedCalls.Failures()
	for _, call := range failures {
		ctrl.t.Errorf("missing call(s) to %v", call)
	}
	if len(failures) != 0 {
<<<<<<< HEAD
		if !cleanup {
			ctrl.T.Fatalf("aborting test due to missing call(s)")
			return
		}
		ctrl.T.Errorf("aborting test due to missing call(s)")
=======
		ctrl.t.Fatalf("aborting test due to missing call(s)")
>>>>>>> 79bfea2d (update vendor)
	}
}

// callerInfo returns the file:line of the call site. skip is the number
// of stack frames to skip when reporting. 0 is callerInfo's call site.
func callerInfo(skip int) string {
	if _, file, line, ok := runtime.Caller(skip + 1); ok {
		return fmt.Sprintf("%s:%d", file, line)
	}
	return "unknown file"
}

<<<<<<< HEAD
// isCleanuper checks it if t's base TestReporter has a Cleanup method.
func isCleanuper(t TestReporter) (cleanuper, bool) {
	tr := unwrapTestReporter(t)
	c, ok := tr.(cleanuper)
	return c, ok
}

// unwrapTestReporter unwraps TestReporter to the base implementation.
func unwrapTestReporter(t TestReporter) TestReporter {
	tr := t
	switch nt := t.(type) {
	case *cancelReporter:
		tr = nt.t
		if h, check := tr.(*nopTestHelper); check {
			tr = h.t
		}
	case *nopTestHelper:
		tr = nt.t
	default:
		// not wrapped
	}
	return tr
=======
type testHelper interface {
	TestReporter
	Helper()
>>>>>>> 79bfea2d (update vendor)
}
// Copyright © 2020 sqos <sqos4os@yandex.com>
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

// waitroutine包通过WaitRoutine类型提供了创建/等待/取消go routine的方便方式
//
// 在一般使用go routine时,都需要使用Context进行管理和取消.等待退出常用WaitGroup.
// WaitRoutine将两者结合起来,通过简单的接口达到快捷使用效果.
//
// 典型使用场景如下,通过WaitRoutine.Go()运行某些特定功能的go routine,并等待其退出.
//
//  wg := waitroutine.New(nil)
//  wg.Go(func() {
//     // do something
//  }).Go(func() {
//     // do other something
//  })
//  wg.Wait()
//
// 也可以通过接收到某种信号进行Cancel(),通过WaitRoutine.GoRoutine()运行一个持久
// 运行类型为Routine的go routine,在满足特定条件时,退出.
// 这个特定条件可以是ctx传递进去,也可以是特定功能运行结束.如下通过ctx.Done()退出:
//  routine := func(ctx context.Context) {
//    tick := time.NewTicker(time.Second)
//    for {
//      select {
//      case <-tick.C:
//      case <-ctx.Done():
//    	  return
//      }
//    }
//  }
//  wg := New(context.Background())
//  time.AfterFunc(waitSecond, func() {
//     fmt.Println("it's time to cancel")
//     wg.Cancel()
//  })
//
//  wg.GoRoutine(routine)
//  wg.GoRoutine(routine)
//
//  wg.Wait()
//
package waitroutine

import (
	"context"
	"sync"
)

// Routine 可以通过Go()函数运行的routine原型
type Routine func(ctx context.Context)

// WaitRoutine 管理go routine
type WaitRoutine struct {
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// DefaultWaitRoutine 默认WaitRoutine
var DefaultWaitRoutine = New(context.Background())

// New 新建一个WaitRoutine
//
// 在ctx为nil值时,默认使用context.Background()作为父context
func New(ctx context.Context) *WaitRoutine {
	wgc := &WaitRoutine{}
	if ctx == nil {
		ctx = context.Background()
	}
	wgc.ctx, wgc.cancelFunc = context.WithCancel(ctx)
	return wgc
}

func (c *WaitRoutine) goFn(fn func()) {
	fn()
	c.wg.Done()
}

// Go 运行参数传递的routines,类型为func()
//
// 接收不定个数func(),所有都会运行
// 该接口一般用于不需要context的go routine调用
func (c *WaitRoutine) Go(fns ...func()) *WaitRoutine {
	for _, fn := range fns {
		c.wg.Add(1)
		go c.goFn(fn)
	}
	return c
}

func (c *WaitRoutine) goRoutine(routine Routine) {
	routine(c.ctx)
	c.wg.Done()
}

// GoRoutine 运行参数传递的routines,类型Routine
//
// 接收不定个数Routine,所有都会运行
// 该接口会传递context.Context,go routine可以根据context决定是否结束,或者从中获取相关参数
func (c *WaitRoutine) GoRoutine(routines ...Routine) *WaitRoutine {
	for _, routine := range routines {
		c.wg.Add(1)
		go c.goRoutine(routine)
	}
	return c
}

// Cancel 取消所有Routine运行,如果已经运行,则ctx参数会接收到ctx.Done()信号
func (c *WaitRoutine) Cancel() {
	c.cancelFunc()
}

// Wait 等待所有Routine运行结束或者被取消
func (c *WaitRoutine) Wait() {
	c.wg.Wait()
}

// WaitGroup 返回内部WaitGroup结构
func (c *WaitRoutine) WaitGroup() *sync.WaitGroup {
	return &c.wg
}

// Context 返回内部Context结构
func (c *WaitRoutine) Context() context.Context {
	return c.ctx
}

// Go 通过DefaultWaitRoutine运行参数传递的routines,类型为func()
//
// 接收不定个数func(),所有都会运行
func Go(fns ...func()) *WaitRoutine {
	return DefaultWaitRoutine.Go(fns...)
}


// Go 通过DefaultWaitRoutine运行参数传递的routines,类型为Routine
//
// 接收不定个数Routine,所有都会运行
func GoRoutine(routines ...Routine) *WaitRoutine {
	return DefaultWaitRoutine.GoRoutine(routines...)
}

// Cancel 通过DefaultWaitRoutine取消所有Routine运行,
// 如果已经运行,则ctx参数会接收到ctx.Done()信号
func Cancel() {
	DefaultWaitRoutine.cancelFunc()
}

// Wait 通过DefaultWaitRoutine等待所有Routine运行结束或者被取消
func Wait() {
	DefaultWaitRoutine.Wait()
}

// WaitGroup 通过DefaultWaitRoutine返回内部WaitGroup结构
func WaitGroup() *sync.WaitGroup {
	return DefaultWaitRoutine.WaitGroup()
}

// Context 通过DefaultWaitRoutine返回内部Context结构
func Context() context.Context {
	return DefaultWaitRoutine.Context()
}

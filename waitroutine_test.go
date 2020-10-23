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

package waitroutine

import (
	"context"
	"testing"
	"time"
)

func timestamp() string {
	return time.Now().Format(time.RFC3339)
}

func routine(ctx context.Context) {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-tick.C:
		case <-ctx.Done():
			return
		}
	}
}

func TestWaitRoutine_Wait(t *testing.T) {
	waitSecond := time.Second * 5
	ctx, _ := context.WithTimeout(context.Background(), waitSecond)

	wg := New(ctx)
	wg.GoRoutine(routine)
	wg.GoRoutine(routine)
	wg.Go(func() {
		after := 10 * time.Second
		<-time.After(after)
		t.Logf("%s time after %s exit", timestamp(), after)
	})

	t.Logf("%s wait for timeout after %v", timestamp(), waitSecond)
	wg.Wait()
	t.Logf("%s now exit", timestamp())
}

func TestWaitRoutine_Cancel(t *testing.T) {
	waitSecond := time.Second * 5

	wg := New(context.Background())
	time.AfterFunc(waitSecond, func() {
		t.Logf("%s time to cancel", timestamp())
		wg.Cancel()
	})

	wg.GoRoutine(routine)
	wg.GoRoutine(routine)
	wg.Go(func() {
		after := 10 * time.Second
		<-time.After(after)
		t.Logf("%s time after %s exit", timestamp(), after)
	})

	t.Logf("%s wait for cancel after %v", timestamp(), waitSecond)
	wg.Wait()
	t.Logf("%s now exit", timestamp())
}

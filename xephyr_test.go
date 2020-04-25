// Copyright 2020 Aleksandar Radivojevic
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xvfb

import (
	"testing"
	"time"
)

func TestXephyr(t *testing.T) {
	time.Sleep(500 * time.Millisecond)

	xephyr := NewXephyr(Options{
		Display: 99,
		Width:   1280,
		Height:  720,
		SetEnv:  false,
	})

	err := xephyr.Start()
	if err != nil {
		t.Errorf("cannot start xephyr\n%v", err)
		t.FailNow()
	}

	defer func() {
		if err := xephyr.Stop(); err != nil {
			t.Errorf("stopping xephyr returned an error\n%v", err)
		}

		if err := xephyr.Wait(); err != nil {
			t.Errorf("xephyr did not quit peacefully\n%v", err)
		}

		if xephyr.IsRunning() {
			t.Errorf("xephyr is still detected as if it was running even though should have been stopped")
		}
	}()

	if err := xephyr.WaitUntilReady(15); err != nil {
		t.Errorf("xephyr did not get ready\n%v", err)
	}
}

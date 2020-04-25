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

func TestXvfb(t *testing.T) {
	time.Sleep(500 * time.Millisecond)

	xvfb := NewXvfb(Options{
		Display:    99,
		Width:      1280,
		Height:     720,
		ColorDepth: 24,
	})

	err := xvfb.Start()
	if err != nil {
		t.Errorf("cannot start xvfb\n%v", err)
		t.FailNow()
	}

	defer func() {
		if err := xvfb.Stop(); err != nil {
			t.Errorf("stopping xvfb returned an error\n%v", err)
		}

		if err := xvfb.Wait(); err != nil {
			t.Errorf("xvfb did not quit peacefully\n%v", err)
		}

		if xvfb.IsRunning() {
			t.Errorf("xvfb is still detected as if it was running even though should have been stopped")
		}
	}()

	if err := xvfb.WaitUntilReady(15); err != nil {
		t.Errorf("xvfb did not get ready\n%v", err)
	}
}

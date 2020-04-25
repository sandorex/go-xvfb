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
)

// really basic test of running and closing without errors
func testDisplay(display Display, testName string, t *testing.T) {
	t.Run(testName, func(t *testing.T) {
		err := display.Start()
		if err != nil {
			t.Fatalf("cannot start display %q:\n%v", display.GetBackend(), err)
		}

		defer func() {
			if err := display.Stop(); err != nil {
				t.Errorf("stopping display %q returned an error:\n%v", display.GetBackend(), err)
			}

			if err := display.Wait(); err != nil {
				t.Errorf("display %q did not quit peacefully:\n%v", display.GetBackend(), err)
			}

			if display.IsRunning() {
				t.Errorf("display %q is still detected as if it was running even though should have been stopped", display.GetBackend())
			}
		}()

		if err := display.WaitUntilReady(15); err != nil {
			t.Errorf("display %q did not get ready\n%v", display.GetBackend(), err)
		}
	})
}

func TestBasic(t *testing.T) {
	xvfb := NewXvfb(99, 1280, 720)
	xephyr := NewXephyr(88, 1280, 720, 99)

	if ok, errs := xvfb.HasDependencies(); !ok {
		t.Logf("xvfb dependencies are missing:")
		for _, err := range errs {
			t.Log(err)
		}
		t.SkipNow()
	}

	if ok, errs := xephyr.HasDependencies(); !ok {
		t.Logf("xephyr dependencies are missing:")
		for _, err := range errs {
			t.Log(err)
		}
		t.SkipNow()
	}

	// test xvfb first
	testDisplay(xvfb, "xvfb", t)

	t.Log("running Xephyr under Xvfb, it will fail if Xvfb test failed")

	err := xvfb.Start()
	if err != nil {
		t.Fatalf("cannot start xvfb:\n%v", err)
	}

	defer func() {
		if err := xvfb.Stop(); err != nil {
			t.Errorf("stopping xvfb returned an error:\n%v", err)
		}

		if err := xvfb.Wait(); err != nil {
			t.Errorf("xvfb did not quit peacefully:\n%v", err)
		}
	}()

	if err := xvfb.WaitUntilReady(15); err != nil {
		t.Fatalf("xvfb did not get ready\n%v", err)
	}

	testDisplay(xephyr, "xephyr", t)
}

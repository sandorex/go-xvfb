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

// Display is generic interface for all displays
type Display interface {
	Start() error
	Stop() error
	Kill() error
	Wait() error
	IsRunning() bool
	IsReady() (bool, error)
	WaitUntilReady(timeout int) error
	IsVisible() bool
	GetDependencies() []string
	HasDependencies() (bool, []error)
}

// StartDisplay starts a display
func StartDisplay(visible bool, options Options) (Display, error) {
	var display Display
	if visible {
		display = NewXvfb(options)
	} else {
		display = NewXephyr(options)
	}

	return display, display.Start()
}

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
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Xephyr is a display backend that opens a virtual x server as a window in
// existing x server, as such it requires a running x server
type Xephyr struct {
	baseDisplay
	HostDisplay int
}

// check if it implements Display
var _ Display = (*Xephyr)(nil)

// NewXephyr constructor for Xephyr
func NewXephyr(display, width, height, hostDisplay int) *Xephyr {
	xephyr := Xephyr{}
	xephyr.Display = display
	xephyr.Width = width
	xephyr.Height = height
	xephyr.HostDisplay = hostDisplay

	if xephyr.Width == 0 || xephyr.Height == 0 {
		xephyr.Width = 1280
		xephyr.Height = 720
	}

	return &xephyr
}

// Start starts Xephyr
func (x *Xephyr) Start() error {
	if x.IsRunning() {
		return ErrAlreadyRunning
	}

	if !isDisplayInUse(x.HostDisplay) {
		return ErrNoDisplay
	}

	display := fmt.Sprintf(":%d", x.Display)

	if isDisplayInUse(x.Display) {
		return ErrDisplayInUse(x.Display)
	}

	x.cmd = exec.Command(
		"Xephyr",
		append([]string{
			fmt.Sprintf(":%d", x.Display),
			"-screen",
			fmt.Sprintf("%dx%d", x.Width, x.Height),
		}, x.Args...)...,
	)

	x.cmd.Env = append(os.Environ(), fmt.Sprintf("DISPLAY=:%d", x.HostDisplay))

	// start and return if any error rises
	if err := x.cmd.Start(); err != nil {
		return err
	}

	// wait for it to quit if it's going to quit
	time.Sleep(500 * time.Millisecond)

	// check if it's still running
	if !x.IsRunning() {
		return ErrCrashed
	}

	if err := os.Setenv("DISPLAY", display); err != nil {
		return err
	}

	return nil
}

// IsVisible is the backend visible
func (Xephyr) IsVisible() bool {
	return true
}

// GetBackend returns name of the backend
func (Xephyr) GetBackend() string {
	return "xephyr"
}

// GetDependencies returns dependencies that should be in path for it to run
func (Xephyr) GetDependencies() []string {
	return []string{
		"Xephyr",
		"xdpyinfo",
	}
}

// HasDependencies checks if all required dependencies are in path
func (x Xephyr) HasDependencies() (bool, []error) {
	return hasDependencies(x.GetDependencies())
}

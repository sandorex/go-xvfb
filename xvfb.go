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

// Xvfb is a display backend that opens a virtual x server in the background
type Xvfb struct {
	baseDisplay
}

// check if it implements Display
var _ Display = (*Xvfb)(nil)

// NewXvfb constructor for Xvfb
func NewXvfb(display, width, height int) *Xvfb {
	xvfb := Xvfb{}
	xvfb.Display = display
	xvfb.Width = width
	xvfb.Height = height

	if xvfb.Width == 0 || xvfb.Height == 0 {
		xvfb.Width = 1280
		xvfb.Height = 720
	}

	return &xvfb
}

// Start starts Xvfb
func (x *Xvfb) Start() error {
	if x.IsRunning() {
		return ErrAlreadyRunning
	}

	display := fmt.Sprintf(":%d", x.Display)

	if isDisplayInUse(x.Display) {
		return ErrDisplayInUse(x.Display)
	}

	x.cmd = exec.Command(
		"Xvfb",
		append([]string{
			fmt.Sprintf(":%d", x.Display),
			"-screen",
			"0",
			fmt.Sprintf("%dx%dx24", x.Width, x.Height),
		}, x.Args...)...,
	)

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
func (Xvfb) IsVisible() bool {
	return false
}

// GetBackend returns name of the backend
func (Xvfb) GetBackend() string {
	return "xvfb"
}

// GetDependencies returns dependencies that should be in path for it to run
func (Xvfb) GetDependencies() []string {
	return []string{
		"Xvfb",
		"xdpyinfo",
	}
}

// HasDependencies checks if all required dependencies are in path
func (x Xvfb) HasDependencies() (bool, []error) {
	return hasDependencies(x.GetDependencies())
}

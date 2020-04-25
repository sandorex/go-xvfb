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
	"syscall"
	"time"
)

// Xvfb is a display backend that opens a virtual x server in the background
type Xvfb struct {
	*Options
	cmd *exec.Cmd
}

// check if it implements display
var _ Display = (*Xvfb)(nil)

// NewXvfb Xvfb constructor
func NewXvfb(options Options) *Xvfb {
	// set default values
	if options.ColorDepth == 0 {
		options.ColorDepth = 24
	}

	if options.Width == 0 || options.Height == 0 {
		options.Width = 1280
		options.Height = 720
	}

	return &Xvfb{Options: &options, cmd: nil}
}

// Start starts Xvfb
func (x *Xvfb) Start() error {
	if x.IsRunning() {
		return ErrAlreadyRunning
	}

	display := fmt.Sprintf(":%d", x.Display)

	_, err := os.Stat(fmt.Sprintf("/tmp/.X11-unix/X%d", x.Display))
	if !os.IsNotExist(err) {
		return ErrDisplayInUse(x.Display)
	}

	// TODO capture process output and return it with an error
	x.cmd = exec.Command(
		"Xvfb",
		append([]string{
			display,
			"-screen",
			"0",
			fmt.Sprintf("%dx%dx%d", x.Width, x.Height, x.ColorDepth),
		}, x.Args...)...,
	)

	if err := os.Setenv("DISPLAY", display); err != nil {
		return err
	}

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

	return nil
}

// Stop stops Xvfb gracefully
func (x Xvfb) Stop() error {
	if !x.IsRunning() {
		return ErrNotRunning
	}

	return x.cmd.Process.Signal(syscall.SIGTERM)
}

// Kill stops Xvfb forcefully
func (x Xvfb) Kill() error {
	if !x.IsRunning() {
		return ErrNotRunning
	}

	return x.cmd.Process.Kill()
}

// Wait waits for Xvfb to quit
func (x Xvfb) Wait() error {
	if !x.IsRunning() {
		return ErrNotRunning
	}

	return x.cmd.Wait()
}

// IsRunning checks if Xvfb process is running
func (x Xvfb) IsRunning() bool {
	if x.cmd == nil || x.cmd.Process == nil {
		return false
	}

	return x.cmd.Process.Signal(syscall.Signal(0)) == nil
}

// IsReady checks if Xvfb is ready to display windows
func (x Xvfb) IsReady() (bool, error) {
	return isDisplayReady(x.Display)
}

// WaitUntilReady waits until the Xvfb is ready or timeout has been exceeded
func (x Xvfb) WaitUntilReady(timeout int) error {
	if !x.IsRunning() {
		return ErrNotRunning
	}

	i := 0
	for timeout == 0 || i < timeout {
		ready, err := x.IsReady()

		// NOTE ExitError is caught by IsReady(), so this is for other errors
		if err != nil {
			return err
		}

		if ready {
			return nil
		}

		time.Sleep(1 * time.Second)

		if timeout != 0 {
			i++
		}
	}

	return ErrTimeout(timeout)
}

// IsVisible is the backend visible
func (Xvfb) IsVisible() bool {
	return false
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

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

// Xephyr is a display backend that opens a virtual x server as a window in
// existing x server, as such it requires a running x server
type Xephyr struct {
	*Options
	cmd *exec.Cmd
}

// check if it implements display
var _ Display = (*Xephyr)(nil)

// NewXephyr Xephyr constructor
func NewXephyr(options Options) *Xephyr {
	// set default values
	if options.Width == 0 || options.Height == 0 {
		options.Width = 1280
		options.Height = 720
	}

	return &Xephyr{Options: &options, cmd: nil}
}

// Start starts Xephyr
func (x *Xephyr) Start() error {
	if x.IsRunning() {
		return ErrAlreadyRunning
	}

	// TODO check if xserver is running

	display := fmt.Sprintf(":%d", x.Display)

	_, err := os.Stat(fmt.Sprintf("/tmp/.X11-unix/X%d", x.Display))
	if !os.IsNotExist(err) {
		return ErrDisplayInUse(x.Display)
	}

	// TODO capture process output and return it with an error
	x.cmd = exec.Command(
		"Xephyr",
		append([]string{
			display,
			"-screen",
			fmt.Sprintf("%dx%d", x.Width, x.Height),
		}, x.Args...)...,
	)

	if x.SetEnv {
		if err := os.Setenv("DISPLAY", display); err != nil {
			return err
		}
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

// Stop stops Xephyr gracefully
func (x Xephyr) Stop() error {
	if !x.IsRunning() {
		return ErrNotRunning
	}

	return x.cmd.Process.Signal(syscall.SIGTERM)
}

// Kill stops Xephyr forcefully
func (x Xephyr) Kill() error {
	if !x.IsRunning() {
		return ErrNotRunning
	}

	return x.cmd.Process.Kill()
}

// Wait waits for Xephyr to quit
func (x Xephyr) Wait() error {
	if !x.IsRunning() {
		return ErrNotRunning
	}

	return x.cmd.Wait()
}

// IsRunning checks if Xephyr process is running
func (x Xephyr) IsRunning() bool {
	if x.cmd == nil || x.cmd.Process == nil {
		return false
	}

	return x.cmd.Process.Signal(syscall.Signal(0)) == nil
}

// IsReady checks if Xephyr is ready to display windows
func (x Xephyr) IsReady() (bool, error) {
	return isDisplayReady(x.Display)
}

// WaitUntilReady waits until the Xephyr is ready or timeout has been exceeded
func (x Xephyr) WaitUntilReady(timeout int) error {
	if !x.IsRunning() {
		return ErrNotRunning
	}

	i := 0
	for timeout == 0 || i < timeout {
		ready, err := x.IsReady()

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
func (Xephyr) IsVisible() bool {
	return true
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

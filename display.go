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
	"os/exec"
	"syscall"
	"time"
)

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
	GetBackend() string
	GetDependencies() []string
	HasDependencies() (bool, []error)
}

// TODO capture process output and return it with an error

type baseDisplay struct {
	Display       int
	Width, Height int
	Args          []string
	cmd           *exec.Cmd
}

// Stop stops the display gracefully
func (x baseDisplay) Stop() error {
	if !x.IsRunning() {
		return ErrNotRunning
	}

	return x.cmd.Process.Signal(syscall.SIGTERM)
}

// Kill stops the display forcefully
func (x baseDisplay) Kill() error {
	if !x.IsRunning() {
		return ErrNotRunning
	}

	return x.cmd.Process.Kill()
}

// Wait waits for the display to quit
func (x baseDisplay) Wait() error {
	if !x.IsRunning() {
		return ErrNotRunning
	}

	return x.cmd.Wait()
}

// IsRunning checks if display process is running
func (x baseDisplay) IsRunning() bool {
	if x.cmd == nil || x.cmd.Process == nil {
		return false
	}

	return x.cmd.Process.Signal(syscall.Signal(0)) == nil
}

// IsReady checks if display is ready to be used
func (x baseDisplay) IsReady() (bool, error) {
	return isDisplayReady(x.Display)
}

// WaitUntilReady waits until the display is ready or timeout has been exceeded
//
// Timeout of 0 is infinite
func (x baseDisplay) WaitUntilReady(timeout int) error {
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

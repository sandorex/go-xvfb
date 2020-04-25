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
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// ErrCrashed means the display has quit unexpectedly
var ErrCrashed = errors.New("display has quit unexpectedly")

// ErrAlreadyRunning means the display is already running
var ErrAlreadyRunning = errors.New("display is already running")

// ErrNotRunning means the display is not running
var ErrNotRunning = errors.New("display is not running")

// ErrNoDisplay means the host display is not running
var ErrNoDisplay = errors.New("xephyr requires a running x server")

// ErrTimeout means timeout has been exceeded
type ErrTimeout int

func (e ErrTimeout) Error() string {
	return fmt.Sprintf("timeout of %ds has been exceeded", int(e))
}

// ErrDisplayInUse means the DISPLAY set is already in use by a x server
type ErrDisplayInUse int

func (e ErrDisplayInUse) Error() string {
	return fmt.Sprintf("display %d is in use, please remove the lockfile at /tmp/.X11-unix/X%d", int(e), int(e))
}

// isDisplayInUse checks if display is in use
func isDisplayInUse(display int) bool {
	_, err := os.Stat(fmt.Sprintf("/tmp/.X11-unix/X%d", display))
	return !os.IsNotExist(err)
}

// isDisplayReady checks if display is ready
func isDisplayReady(display int) (bool, error) {
	cmd := exec.Command("xdpyinfo")
	cmd.Env = append(os.Environ(), fmt.Sprintf("DISPLAY=:%d", display))

	err := cmd.Run()
	if errors.Is(err, &exec.ExitError{}) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// hasDependencies checks if all dependencies can be found in path
func hasDependencies(dependencies []string) (bool, []error) {
	errors := []error{}
	for _, i := range dependencies {
		if _, err := exec.LookPath(i); err != nil {
			errors = append(errors, err)
		}
	}

	return len(errors) == 0, errors
}

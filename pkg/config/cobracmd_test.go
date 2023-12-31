// Copyright © 2023 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShowConfigCobraCommandOk(t *testing.T) {

	err := ShowConfigCommand(func() error {
		AddRootKey("things")
		return nil
	}).Execute()
	assert.NoError(t, err)

}

func TestShowConfigCobraCommandErrInit(t *testing.T) {

	err := ShowConfigCommand(func() error {
		return fmt.Errorf("pop")
	}).Execute()
	assert.Regexp(t, "pop", err)

}

// Copyright 2020-2021 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestData(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// Call OnVMStart -> set initial value.
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())
	// Initialize http context.
	contextID := host.InitializeHttpContext()
	// Call OnHttpRequestHeaders.
	action := host.CallOnRequestHeaders(contextID, nil, false)
	require.Equal(t, types.ActionContinue, action)

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	require.Contains(t, logs, "shared value: 1")

	// Call OnHttpRequestHeaders again.
	action = host.CallOnRequestHeaders(contextID, nil, false)
	require.Equal(t, types.ActionContinue, action)
	action = host.CallOnRequestHeaders(contextID, nil, false)
	require.Equal(t, types.ActionContinue, action)

	// Check Envoy logs.
	logs = host.GetLogs(types.LogLevelInfo)
	require.Contains(t, logs, "shared value: 3")
}

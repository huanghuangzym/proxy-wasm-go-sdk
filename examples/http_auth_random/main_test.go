package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestHttpAuthRandom_OnHttpRequestHeaders(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// Initialize context.
	contextID := host.InitializeHttpContext()

	// Call OnHttpRequestHeaders.
	action := host.CallOnRequestHeaders(contextID,
		types.Headers{{"key", "value"}}, false)
	require.Equal(t, types.ActionPause, action)

	// Verify DispatchHttpCall is called.
	attrs := host.GetCalloutAttributesFromContext(contextID)
	require.Equal(t, len(attrs), 1)
	require.Equal(t, "httpbin", attrs[0].Upstream)
	// Check if the current action is pause.
	require.Equal(t, types.ActionPause,
		host.GetCurrentHttpStreamAction(contextID))

	// Check Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	require.Contains(t, logs, "http call dispatched to "+clusterName)
	require.Contains(t, logs, "request header: key: value")
}

func TestHttpAuthRandom_OnHttpCallResponse(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithNewRootContext(newRootContext)
	host := proxytest.NewHostEmulator(opt)
	// Release the host emulation lock so that other test cases can insert their own host emulation.
	defer host.Done()

	// http://httpbin.org/uuid
	headers := [][2]string{
		{"HTTP/1.1", "200 OK"}, {"Date:", "Thu, 17 Sep 2020 02:47:07 GMT"},
		{"Content-Type", "application/json"}, {"Content-Length", "53"},
		{"Connection", "keep-alive"}, {"Server", "gunicorn/19.9.0"},
		{"Access-Control-Allow-Origin", "*"}, {"Access-Control-Allow-Credentials", "true"},
	}

	// Access granted case -> Local response must not be sent.
	contextID := host.InitializeHttpContext()
	// Call OnHttpRequestHeaders.
	action := host.CallOnRequestHeaders(contextID, nil,
		false)
	require.Equal(t, types.ActionPause, action)
	// Verify DispatchHttpCall is called.
	attrs := host.GetCalloutAttributesFromContext(contextID)
	require.Equal(t, len(attrs), 1)
	// Call OnHttpCallResponse.
	body := []byte(`{"uuid": "7b10a67a-1c67-4199-835b-cbefcd4a63d4"}`)
	host.CallOnHttpCallResponse(attrs[0].CalloutID, headers, nil, body)
	// Check local response.
	assert.Nil(t, host.GetSentLocalResponse(contextID))
	// CHeck Envoy logs.
	logs := host.GetLogs(types.LogLevelInfo)
	require.Contains(t, logs, "access granted")

	// Access denied case -> Local response must be sent.
	contextID = host.InitializeHttpContext()
	// Call OnHttpRequestHeaders.
	action = host.CallOnRequestHeaders(contextID, nil, false)
	require.Equal(t, types.ActionPause, action)
	// Verify DispatchHttpCall is called.
	attrs = host.GetCalloutAttributesFromContext(contextID)
	require.Equal(t, len(attrs), 1)
	// Call OnHttpCallResponse.
	body = []byte(`{"uuid": "aaaaaaaa-1c67-4199-835b-cbefcd4a63d4"}`)
	host.CallOnHttpCallResponse(attrs[0].CalloutID, headers, nil, body)
	// Check local response.
	localResponse := host.GetSentLocalResponse(contextID)
	assert.NotNil(t, localResponse)
	require.Equal(t, uint32(403), localResponse.StatusCode)
	require.Equal(t, []byte("access forbidden"), localResponse.Data)
	require.Len(t, localResponse.Headers, 1)
	require.Equal(t, "powered-by", localResponse.Headers[0][0])
	require.Equal(t, "proxy-wasm-go-sdk!!", localResponse.Headers[0][1])
	// Check Envoy logs.
	logs = host.GetLogs(types.LogLevelInfo)
	require.Contains(t, logs, "access forbidden")
}

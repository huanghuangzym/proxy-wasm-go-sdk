package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetNewRootContext(newRootContext)
}

type rootContext struct {
	// You'd better embed the default root context
	// so that you don't need to reimplement all the methods by yourself.
	proxywasm.DefaultRootContext
}

func newRootContext(uint32) proxywasm.RootContext { return &rootContext{} }

// Override DefaultRootContext.
func (*rootContext) NewHttpContext(contextID uint32) proxywasm.HttpContext {
	return &httpHeaders{contextID: contextID}
}

type httpHeaders struct {
	// You'd better embed the default root context
	// so that you don't need to reimplement all the methods by yourself.
	proxywasm.DefaultHttpContext
	contextID uint32
}

// Override DefaultHttpContext.
func (ctx *httpHeaders) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("wasm  OnHttpRequestHeaders")
	err := proxywasm.SetHttpRequestHeader("test", "best")
	if err != nil {
		proxywasm.LogCritical("failed to set request header: test")
	}

	_, err = proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
	}

	/*
	for _, h := range hs {
		proxywasm.LogInfof("request header --> %s: %s", h[0], h[1])
	}*/

	err = proxywasm.AddHttpRequestHeader("hl-test","OnHttpRequestHeaders")
	if err != nil {
		proxywasm.LogWarnf("failed to AddHttpRequestHeader: %v", err)
	}

	return types.ActionContinue
}



// Override DefaultHttpContext.
func (ctx *httpHeaders) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("wasm  OnHttpResponseHeaders")

	_, err := proxywasm.GetHttpResponseHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get response headers: %v", err)
	}

	/*for _, h := range hs {
		proxywasm.LogInfof("response header <-- %s: %s", h[0], h[1])
	}*/
	err = proxywasm.AddHttpResponseHeader("hltest","OnHttpResponseHeaders")
	if err != nil {
		proxywasm.LogWarnf("failed to AddHttpRequestHeader: %v", err)
	}

	err = proxywasm.RemoveHttpResponseHeader("content-length")
	if err != nil {
		proxywasm.LogWarnf("failed to RemoveHttpResponseHeader: %v", err)
	}


	return types.ActionContinue
}


func (ctx *httpHeaders) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {
	proxywasm.LogInfof("body size: %d", bodySize)
	if bodySize != 0 {
		initialBody, err := proxywasm.GetHttpResponseBody(0, bodySize)
		if err != nil {
			proxywasm.LogErrorf("failed to get request body: %v", err)
			return types.ActionContinue
		}
		proxywasm.LogInfof("initial request body: %s", string(initialBody))

		b := []byte("testhlbodyrsp zheli bixu zu gou chang de\n")

		err = proxywasm.SetHttpResponseBody(b)
		if err != nil {
			proxywasm.LogErrorf("failed to SetHttpResponseBody: %v", err)
			return types.ActionContinue
		}

		proxywasm.LogInfof("on OnHttpResponseBody finished")
	}

	return types.ActionContinue
}


// Override DefaultHttpContext.
func (ctx *httpHeaders) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}

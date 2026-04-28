// Package tesote is the Go SDK for the equipo.tesote.com API.
//
// Versioned clients ship side by side under v1, v2, v3 sub-packages. The
// top-level package exposes the shared transport, error types, and
// configuration so consumers can construct any version against the same
// underlying *Client.
//
//	import (
//	    "github.com/tesote/sdk/go"
//	    "github.com/tesote/sdk/go/v3"
//	)
//
//	tc, err := tesote.NewClient(tesote.Options{APIKey: "..."} )
//	if err != nil { return err }
//	v3c := v3.New(tc)
//	accounts, err := v3c.Accounts.List(ctx, v3.ListOptions{})
//
// All public I/O methods accept context.Context as their first argument; pass
// a context with a deadline to bound the request.
package tesote

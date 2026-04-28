// Package tesote is the Go SDK for the equipo.tesote.com API.
//
// Versioned clients ship side by side under v1, v2 sub-packages. The
// top-level package exposes the shared transport, error types, and
// configuration so consumers can construct any version against the same
// underlying *Client.
//
//	import (
//	    "github.com/tesote/sdk/go"
//	    "github.com/tesote/sdk/go/v2"
//	)
//
//	tc, err := tesote.NewClient(tesote.Options{APIKey: "..."} )
//	if err != nil { return err }
//	v2c := v2.New(tc)
//	accounts, err := v2c.Accounts.List(ctx, v2.ListOptions{})
//
// All public I/O methods accept context.Context as their first argument; pass
// a context with a deadline to bound the request.
package tesote

package v3

import "errors"

// VerifyWebhookSignature verifies a v3 webhook signature. Stub until the
// platform's signature scheme is confirmed.
func VerifyWebhookSignature(_ []byte, _ string, _ string) error {
	return errors.New("tesote/v3: webhook signature verification not implemented")
}

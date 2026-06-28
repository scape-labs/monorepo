package effectors

import "context"

// WebhookNotify notifies a downstream system that an event happened.
//
// TODO: implement — POST to the registered webhook URL with retry + signing.
func WebhookNotify(ctx context.Context, url string, payload []byte) error {
	_ = ctx
	_ = url
	_ = payload
	// TODO: implement.
	return nil
}

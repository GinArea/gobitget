package bitgetv3

// RateLimit - Bitget does not document per-request rate limit headers, so no
// http-tagged fields are defined; the empty type keeps the Response[T] shape
// shared with the sibling exchange SDKs
type RateLimit struct {
}

package api

// CacheControlImmutable is for static/immutable assets (images, fonts, etc.)
//   - public:       Allows CDNs and shared proxies to cache.
//   - immutable:    Prevents browsers from sending "Is this modified?" (304) checks on refresh.
//   - no-transform: Prevents intermediary compression.
//   - max-age:      1 year (31536000s).
const CacheControlImmutable = "public, immutable, no-transform, max-age=31536000"

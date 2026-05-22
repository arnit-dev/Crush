// Package crush provides the public API for Crush, a terminal-based AI
// coding assistant. It exposes the core types and interfaces needed to
// programmatically configure, initialize, and interact with Crush
// from external Go code.
//
// This package is the stable, public surface of Crush. Types that are
// re-exported here are guaranteed to remain available and stable across
// minor releases. Internal implementation details in internal/ packages
// may change without notice.
//
// For command-line usage, see the crush binary. For HTTP API access, see
// the server/client packages.
package crush

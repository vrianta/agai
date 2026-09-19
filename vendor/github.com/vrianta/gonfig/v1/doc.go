// Package structparser populates Go structs from environment variables,
// command-line arguments, and default values declared in struct tags.
//
// Import this version with:
//
//	import gonfig "github.com/vrianta/gonfig/v1"
//
// The package supports env, arg, default, required, and description tags.
// Use Parse when the application needs to handle errors explicitly, or use
// New for concise initialization during application startup.
package structparser

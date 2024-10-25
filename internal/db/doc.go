// Package db provides a database abstraction layer, encapsulating connection and entity interactions.
//
// External packages must use the provided wrapper functions to interact with database configuration.
// Direct imports of the config package are strictly prohibited.
//
// This package is the second-lowest level in the internal hierarchy and must not import from any other internal
// package besides the erratic package for error handling.
//
// Error codes in this package are in the 200xxx address space. Refer to the status sub-package for details.
package db

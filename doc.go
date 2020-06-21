// Package plug contains a small runtime for communication between the
// host and plugin applications. The only function in this package that
// should be directly called is the `plug.Run` function which will start
// the main loop in each plugin implementation. All other variables,
// constants, types, and functions in this package are only exported
// so that they can be used with the code that is generated using
// `protoc-gen-plug`.
package plug

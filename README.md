# plug

plug is a minimal plugin system for go based on the amazing 
[Hashicorp go-plugin](https://github.com/hashicorp/go-plugin) system.
It uses the same basic concept of launching the plugin as a sub-process
using the `os/exec` package from the standard library and communicating
with the process with the plugin through IPC. Here are a few differences:

- The interface will be generated automatically using a new protobuf code
  generator called `protoc-gen-plug`
- Instead of attaching the stdin/stdout streams to the host process it uses
  them for communication between the plugin and host.
- It uses a custom protocol for communication between the host and plugin,
  this is much faster but also has fewer bells and whistles. There are many
  limitations that will probably not be addressed.
  
## Usage

This package uses the new go protobuf implementation `google.golang.org/protobuf`,

1. Install `protoc-gen-go` (from google.golang.org/protobuf) and `protoc-gen-plug`
   to your path
2. Create a protobuf with a _single_ service definition (See `example/kv.proto`)
3. Generate both the go and plug code (See `example/gen.go`)
4. Implement the generated plugin interface (See `example/kv/main.go`)
5. Use the plugin with the generated `Load` command (See `example/kv/kv_test.go`)
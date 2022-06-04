## Chapter 2 - Structure Data with Protocol Buffers

- Use Protocol Buffers to encode data to be sent over the network
  - consistent schemas for structuring and serializing data
  - guarantees type-safety
  - fast serialization and deserialization - less boilerplate code
  - backward compatibility with versioning
- gRPC: a high-performance remote procedure call (RPC) framework

### Setup

- Install the Protocol Buffer compiler
```bash
wget https://github.com/protocolbuffers/protobuf/\
  releases/download/v21.1/protoc-21.1-osx-aarch_64.zip

unzip protoc-21.1-osx-aarch_64.zip -d /usr/local/protobuf

echo 'export PATH="$PATH:/usr/local/protobuf/bin"' >> ~/.bash_profile
source ~/.bash_profile

protoc --version
# libprotoc 3.19.4

rm protoc-21.1-osx-aarch_64.zip 
```

- Install the protobuf runtime for Go
```bash
go get google.golang.org/protobuf/...@v1.28.0
```

- Update and source `~/.bash_profile`
```bash
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOPATH/bin
```

---
### Define Domain Types as Protocol Buffers

- A convention for Go projects is to put protobuf definitions in an `api` directory
- Use the `repeated` keyword to define a slice of some type
- Compile the protobuf
```bash
protoc api/v1/*.proto \
  --go_out=. \
  --go_opt=paths=source_relative \
  --proto_path=.
```
- Write a Makefile to automate commands

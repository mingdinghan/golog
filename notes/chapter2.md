## Chapter 2 - Structure Data with Protocol Buffers

- Use Protocol Buffers to encode data to be sent over the network
  - consistent schemas for structuring and serializing data
  - guarantees type-safety
  - fast serialization and deserialization - less boilerplate code
  - backward compatibility with versioning
- gRPC: a high-performance remote procedure call (RPC) framework

### Install the Protocol Buffer compiler
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

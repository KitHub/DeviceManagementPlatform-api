# DeviceManagementPlatform-api
DeviceManagementPlatform backend api service

## 环境准备

### protobuf

```
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go get buf.build/go/protovalidate
```

* proto文件的validate插件需要将`validat.proto`文件放到`./proto/buf/validate`文件夹下
* `validate.proto`文件在 https://github.com/bufbuild/protovalidate/blob/main/proto/protovalidate/buf/validate/validate.proto
* validation规则参见 https://protovalidate.com/about/ 
// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: main.proto

package pb

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/gogo/protobuf/gogoproto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	_ "github.com/mwitkow/go-proto-validators"
	_ "github.com/nats-rpc/nrpc"
	_ "github.com/TenderPro/rpckit/app/ticker"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Empty) Validate() error {
	return nil
}
func (this *PingRequest) Validate() error {
	if !(this.SleepTimeMs > 0) {
		return github_com_mwitkow_go_proto_validators.FieldError("SleepTimeMs", fmt.Errorf(`value '%v' must be greater than '0'`, this.SleepTimeMs))
	}
	if !(this.SleepTimeMs < 10) {
		return github_com_mwitkow_go_proto_validators.FieldError("SleepTimeMs", fmt.Errorf(`value '%v' must be less than '10'`, this.SleepTimeMs))
	}
	return nil
}
func (this *PingResponse) Validate() error {
	return nil
}

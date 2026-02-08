//go:build reference

package grpctelemetryservice

/*
Reference Solution - gRPC Client-Side Streaming RPC
==================================================

This file implements a gRPC server for a telemetry service. The RecordTelemetry
RPC is client-side streaming: the client sends a stream of TelemetryData messages,
and the server responds with a single TelemetrySummary (count and sum) when the
client closes the stream.

This connects to gRPC streaming:
- Client streams: client sends many, server sends one (or none)
- stream.Recv(): blocks until next message or EOF
- stream.SendAndClose(): sends final response and closes send direction
- io.EOF: signals client has finished sending (CloseAndRecv on client side)

The exercise teaches:
- Implementing gRPC service interfaces (embed Unimplemented* for forward compat)
- Receiving from a stream in a loop
- Aggregating streamed data (count, sum)
- Sending the final response with SendAndClose
*/

import (
	"io"

	pb "github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/proto"
)

// Server implements the gRPC TelemetryService. Embedding UnimplementedTelemetryServiceServer
// satisfies the interface and allows adding new RPCs to the proto without breaking existing code.
type Server struct {
	pb.UnimplementedTelemetryServiceServer
}

/*
RecordTelemetry - Client-Side Streaming Handler

Receives TelemetryData messages from the client, aggregates count and sum,
then sends TelemetrySummary when client signals done (io.EOF).
*/
func (s *Server) RecordTelemetry(stream pb.TelemetryService_RecordTelemetryServer) error {
	var count int32
	var sum float64

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// Client closed send; send summary and return
			return stream.SendAndClose(&pb.TelemetrySummary{Count: count, Sum: sum})
		}
		if err != nil {
			return err
		}
		count++
		sum += msg.Value
	}
}

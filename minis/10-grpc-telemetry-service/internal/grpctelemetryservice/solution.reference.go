//go:build reference

package grpctelemetryservice

import (
	pb "github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/proto"
)

/*
Reference Solution
==================

This file is the canonical reference for this exercise. It keeps failure paths
explicit when an operation can fail, so callers can decide how to handle
errors at API boundaries.

Read this alongside exercise.go and the tests to understand the intended data
flow, ownership boundaries, and invariants that keep behavior deterministic.

Teaching notes:
- Memory/ownership: make copies when returning mutable data that should not
  alias internal state; share references only when aliasing is intentional.
- Invariants: establish assumptions close to construction, and rely on them in
  smaller helper functions to keep logic easy to audit.
- Error surfaces: prefer explicit returns over hidden panics so learners can
  reason about control flow in production-style code.
*/

// Server is the struct that will implement the gRPC server interface.
type Server struct {
	// You must embed the UnimplementedTelemetryServiceServer on your struct to satisfy
	// the interface for forward compatibility.
	pb.UnimplementedTelemetryServiceServer
}

// RecordTelemetry is the method that handles the client-side streaming RPC.
func (s *Server) RecordTelemetry(stream pb.TelemetryService_RecordTelemetryServer) error {
	// TODO:
	// 1. Initialize variables to hold the aggregated results (count and sum).
	// 2. Start a `for` loop to receive messages from the stream.
	// 3. Inside the loop, call `stream.Recv()`.
	//    a. If the error is `io.EOF`, the client has finished sending.
	//       - Use `stream.SendAndClose()` to send a `TelemetrySummary` message
	//         back to the client with the final count and sum.
	//       - Return `nil` to indicate success.
	//    b. If there is any other error, return it.
	//    c. If there is no error, process the received `TelemetryData` message:
	//       - Increment your count.
	//       - Add the message's value to your sum.
	return nil
}

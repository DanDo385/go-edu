package grpctelemetryservice

import (
	"context"
	"log"
	"net"
	"testing"

	pb "github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// startTestServer starts an in-memory gRPC server for testing.
func startTestServer(t *testing.T) (*grpc.ClientConn, *bufconn.Listener) {
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	pb.RegisterTelemetryServiceServer(s, &Server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	return conn, lis
}

func TestRecordTelemetry(t *testing.T) {
	conn, lis := startTestServer(t)
	defer lis.Close()
	client := pb.NewTelemetryServiceClient(conn)

	stream, err := client.RecordTelemetry(context.Background())
	if err != nil {
		t.Fatalf("Failed to call RecordTelemetry: %v", err)
	}

	points := []float64{10.5, 20.0, 5.5}
	for _, p := range points {
		if err := stream.Send(&pb.TelemetryData{Value: p}); err != nil {
			t.Fatalf("Failed to send point to stream: %v", err)
		}
	}

	summary, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("Failed to receive summary: %v", err)
	}

	if summary.Count != 3 {
		t.Errorf("Expected count 3, got %d", summary.Count)
	}
	if summary.Sum != 36.0 {
		t.Errorf("Expected sum 36.0, got %f", summary.Sum)
	}
}

func TestRecordTelemetry_EmptyStream(t *testing.T) {
	conn, lis := startTestServer(t)
	defer lis.Close()
	client := pb.NewTelemetryServiceClient(conn)

	stream, err := client.RecordTelemetry(context.Background())
	if err != nil {
		t.Fatalf("Failed to call RecordTelemetry: %v", err)
	}

	summary, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("Failed to receive summary: %v", err)
	}

	if summary.Count != 0 {
		t.Errorf("Expected count 0, got %d", summary.Count)
	}
	if summary.Sum != 0.0 {
		t.Errorf("Expected sum 0.0, got %f", summary.Sum)
	}
}
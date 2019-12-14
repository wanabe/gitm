package main
import(
    "log"
    "net"
    "context"

    "google.golang.org/grpc"

    pb "github.com/wanabe/gitm/api/gitm"
)

type server struct {
    pb.UnimplementedLogServer
}

func (s *server) Get(ctx context.Context, in *pb.LogInput) (*pb.Commit, error) {
    return &pb.Commit{}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterLogServer(s, &server{})
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

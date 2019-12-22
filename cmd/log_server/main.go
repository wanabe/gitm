package main
import(
	consoleLog "log"
	"net"
	"context"

	"google.golang.org/grpc"

	pb "github.com/wanabe/gitm/api/gitm"
	"github.com/wanabe/gitm/internal/gitm"
)

type server struct {
	pb.UnimplementedLogServer
}

func fatalIfError(err error, format string) {
	if err == nil {
		return
	}
	consoleLog.Fatalf(format, err)
}

func (s *server) Init(ctx context.Context, iter *pb.LogIterator) (*pb.LogIterator, error) {
	log, err := gitm.NewLog(iter)
	fatalIfError(err, "%v")

	err = log.InitPointers()
	fatalIfError(err, "%v")

	return log.Iter, nil
}

func (s *server) Get(ctx context.Context, iter *pb.LogIterator) (*pb.LogIterator, error) {
	log, err := gitm.NewLog(iter)
	fatalIfError(err, "%v")

	err = log.InitWalker()
	fatalIfError(err, "%v")

	err = log.Get()
	fatalIfError(err, "%v")

	return log.Iter, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	fatalIfError(err, "failed to listen: %v")
	s := grpc.NewServer()
	pb.RegisterLogServer(s, &server{})
	fatalIfError(s.Serve(lis), "failed to serve: %v")
}

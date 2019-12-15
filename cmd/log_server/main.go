package main
import(
	"log"
	"net"
	"context"

	"google.golang.org/grpc"
	"gopkg.in/libgit2/git2go.v27"

	pb "github.com/wanabe/gitm/api/gitm"
)

type server struct {
	pb.UnimplementedLogServer
}

func fatalIfError(err error, format string) {
	if err == nil {
  	return
	}
	log.Fatalf(format, err)
}

func (s *server) Get(ctx context.Context, iter *pb.LogIterator) (*pb.LogIterator, error) {
	path := "./"
	repo := iter.Repository
	if (repo != nil) {
		if (repo.Path != "") {
			path = repo.Path
		}
	}

	r, err := git.OpenRepository(path)
	fatalIfError(err, "%v")

	walker, err := r.Walk()
	fatalIfError(err, "%v")
	walker.Sorting(git.SortTime)
	walker.PushGlob("refs/*")

	oid := new(git.Oid)
	fatalIfError(walker.Next(oid), "%v")
	commit, err := r.LookupCommit(oid)
	fatalIfError(err, "%v")
	iter.Commits = append(iter.Commits, &pb.Commit{Object: &pb.Object{Hash: commit.Id()[:]}})
	return iter, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	fatalIfError(err, "failed to listen: %v")
	s := grpc.NewServer()
	pb.RegisterLogServer(s, &server{})
	fatalIfError(s.Serve(lis), "failed to serve: %v")
}

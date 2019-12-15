package main
import(
	"log"
	"net"
	"context"
	"bytes"

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

func repository(iter *pb.LogIterator) (*git.Repository, error) {
	path := "./"
	repo := iter.Repository
	if (repo != nil) {
		if (repo.Path != "") {
			path = repo.Path
		}
	}

	r, err := git.OpenRepository(path)
	fatalIfError(err, "%v")
	return r, nil
}

func (s *server) Init(ctx context.Context, iter *pb.LogIterator) (*pb.LogIterator, error) {
	r, err := repository(iter)
	fatalIfError(err, "%v")

	if len(iter.Pointers) == 0 {
		refIter, err := r.NewReferenceIterator()
		fatalIfError(err, "%v")
		for {
			ref, err := refIter.Next()
			if ref == nil {
				break
			}
			fatalIfError(err, "%v")
			if ref.Type() != git.ReferenceOid {
				continue
			}
			iter.Pointers = append(iter.Pointers, &pb.Object{Hash: ref.Target()[:]})
		}
	}
	return iter, nil
}

func (s *server) Get(ctx context.Context, iter *pb.LogIterator) (*pb.LogIterator, error) {
	r, err := repository(iter)
	fatalIfError(err, "%v")

	walker, err := r.Walk()
	fatalIfError(err, "%v")
	walker.Sorting(git.SortTime)

	for i := range iter.Pointers {
		oid := git.NewOidFromBytes(iter.Pointers[i].Hash)
		walker.Push(oid)
	}

	oid := new(git.Oid)
	fatalIfError(walker.Next(oid), "%v")
	commit, err := r.LookupCommit(oid)
	fatalIfError(err, "%v")

	hash := oid[:]
	iter.Commits = append(iter.Commits, &pb.Commit{Object: &pb.Object{Hash: hash}})
	for i := 0; i < len(iter.Pointers); i++ {
		if (bytes.Equal(hash, iter.Pointers[i].Hash)) {
			iter.Pointers = append(iter.Pointers[:i], iter.Pointers[i+1:]...)
			i--
		}
	}
	count := commit.ParentCount()
	for j := uint(0); j < count; j++ {
		iter.Pointers = append(iter.Pointers, &pb.Object{Hash: commit.ParentId(j)[:]})
	}
	return iter, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	fatalIfError(err, "failed to listen: %v")
	s := grpc.NewServer()
	pb.RegisterLogServer(s, &server{})
	fatalIfError(s.Serve(lis), "failed to serve: %v")
}

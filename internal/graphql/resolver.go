package graphql

import (
	"context"
	"encoding/hex"

	"google.golang.org/grpc"

	pb "github.com/wanabe/gitm/api/gitm"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func input2pb(input *LogIteratorInput) (*pb.LogIterator) {
	pbLogIterator := &pb.LogIterator { Repository: &pb.Repository{}}
	if (input != nil) {
		if (input.Num != nil) {
			pbLogIterator.Num = int32(*input.Num)
		}
		if (input.Repository != nil && input.Repository.Path != nil) {
			pbLogIterator.Repository.Path = *input.Repository.Path
		}
		pbLogIterator.Pointers =  make([]*pb.Object, len(input.Pointers))
		for i, pointer := range input.Pointers {
			b, _ := hex.DecodeString(*pointer.Hash)
			pbLogIterator.Pointers[i] = &pb.Object{Hash: b}
		}
	}
	return pbLogIterator
}

func pb2log(pbLogIterator *pb.LogIterator) (*LogIterator) {
	logIterator := &LogIterator{Num: int(pbLogIterator.Num), Repository: &Repository{Path: pbLogIterator.Repository.Path}}
	logIterator.Pointers =  make([]*Object, len(pbLogIterator.Pointers))
	for i, pointer := range pbLogIterator.Pointers {
		logIterator.Pointers[i] = &Object{Hash: hex.EncodeToString(pointer.Hash)}
	}
	logIterator.Commits =  make([]*Commit, len(pbLogIterator.Commits))
	for i, commit := range pbLogIterator.Commits {
		logIterator.Commits[i] = &Commit{Object: &Object{Hash: hex.EncodeToString(commit.Object.Hash)}}
	}
	return logIterator
}

func (r *queryResolver) Get(ctx context.Context, input *LogIteratorInput) (*LogIterator, error) {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewLogClient(conn)
	pbLogIterator := input2pb(input)
	ret, err := client.Get(context.TODO(), pbLogIterator)
	if err != nil {
		return nil, err
	}
	return pb2log(ret), nil
}

func (r *queryResolver) Init(ctx context.Context, input *LogIteratorInput) (*LogIterator, error) {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewLogClient(conn)
	pbLogIterator := input2pb(input)
	ret, err := client.Init(context.TODO(), pbLogIterator)
	if err != nil {
		return nil, err
	}
	return pb2log(ret), nil
}

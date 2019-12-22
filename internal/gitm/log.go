package gitm

import(
	"bytes"

	"gopkg.in/libgit2/git2go.v27"

	pb "github.com/wanabe/gitm/api/gitm"
)

type Log struct {
	Repository *git.Repository
	Walker *git.RevWalk

	Iter *pb.LogIterator
}

func NewLog(iter *pb.LogIterator) (*Log, error) {
	path := "./"
	repo := iter.Repository
	if (repo != nil) {
		if (repo.Path != "") {
			path = repo.Path
		}
	}

	r, err := git.OpenRepository(path)
	if (err != nil) {
		return nil, err
	}
	return &Log{Repository: r, Iter: iter}, err
}

func (log *Log) InitPointers() (error) {
	if len(log.Iter.Pointers) == 0 {
		refIter, err := log.Repository.NewReferenceIterator()
		if err != nil {
			return err
		}
		for {
			ref, err := refIter.Next()
			if ref == nil {
				break
			}
			if err != nil {
				return err
			}
			if ref.Type() != git.ReferenceOid {
				continue
			}
			log.Iter.Pointers = append(log.Iter.Pointers, &pb.Object{Hash: ref.Target()[:]})
		}
	}
	return nil
}

func (log *Log) InitWalker() (error) {
	walker, err := log.Repository.Walk()
	if err != nil {
		return err
	}
	walker.Sorting(git.SortTime)

	for i := range log.Iter.Pointers {
		oid := git.NewOidFromBytes(log.Iter.Pointers[i].Hash)
		walker.Push(oid)
	}

	log.Walker = walker
	return nil
}

func (log *Log) Get() (error) {
	oid := new(git.Oid)
	err := log.Walker.Next(oid)
	if err != nil {
		return err
	}
	commit, err := log.Repository.LookupCommit(oid)
	if err != nil {
		return err
	}

	hash := oid[:]
	log.Iter.Commits = append(log.Iter.Commits, &pb.Commit{Object: &pb.Object{Hash: hash}})
	for i := 0; i < len(log.Iter.Pointers); i++ {
		if (bytes.Equal(hash, log.Iter.Pointers[i].Hash)) {
			log.Iter.Pointers = append(log.Iter.Pointers[:i], log.Iter.Pointers[i+1:]...)
			i--
		}
	}
	count := commit.ParentCount()
	for j := uint(0); j < count; j++ {
		log.Iter.Pointers = append(log.Iter.Pointers, &pb.Object{Hash: commit.ParentId(j)[:]})
	}
	return nil
}

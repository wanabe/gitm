package main
import (
    "fmt"
    "encoding/hex"

    "github.com/golang/protobuf/proto"

    pb "github.com/wanabe/gitm/api/gitm"
)

func main() {
    commit := pb.Commit {
        Object: &pb.Object {
            Hash: "4df003e28a16e91e9667c7e6ea5852202820ac67",
        },
    }
    out, err := proto.Marshal(&commit)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%s\ndone.\n", hex.Dump(out))
}

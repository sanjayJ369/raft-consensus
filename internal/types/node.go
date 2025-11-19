package types

type Term int64   // represents election term
type Index int64  // represetns log index
type NodeId int64 // represents node Id

type Vote struct {
	Term        Term   // current election term
	VoteGranted bool   // is the vote granted or not
	From        NodeId // follower node id
	To          NodeId // candidate node id
}

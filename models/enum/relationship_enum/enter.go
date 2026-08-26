package relationship_enum

type Relation int8

const (
	RelationStranger Relation = 1
	RelationFocus    Relation = 2
	RelationFans     Relation = 3
	RelationFriend   Relation = 4
)

func (r Relation) String() string {
	switch r {
	case RelationStranger:
		return "Stranger"
	case RelationFocus:
		return "Focus"
	case RelationFans:
		return "Fans"
	case RelationFriend:
		return "Friend"
	}
	return "Unknown"
}

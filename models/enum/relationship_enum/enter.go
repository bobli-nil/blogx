package relationship_enum

type Relation int8

const (
	RelationStranger Relation = 1
	RelationFocus    Relation = 2
	RelationFans     Relation = 3
	RelationFriend   Relation = 4
)

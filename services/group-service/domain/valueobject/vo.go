package valueobject

type GroupID struct {
	value string
}

func NewGroupID(value string) GroupID {
	return GroupID{value: value}
}

func (id GroupID) String() string {
	return id.value
}

func (id GroupID) IsValid() bool {
	return id.value != ""
}

type GroupName struct {
	value string
}

func NewGroupName(value string) GroupName {
	return GroupName{value: value}
}

func (n GroupName) String() string {
	return n.value
}

func (n GroupName) IsValid() bool {
	return len(n.value) >= 2 && len(n.value) <= 100
}
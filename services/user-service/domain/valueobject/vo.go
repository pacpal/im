package valueobject

type UserID struct {
	value string
}

func NewUserID(value string) UserID {
	return UserID{value: value}
}

func (id UserID) String() string {
	return id.value
}

func (id UserID) IsValid() bool {
	return id.value != ""
}

type PhoneNumber struct {
	value string
}

func NewPhoneNumber(value string) PhoneNumber {
	return PhoneNumber{value: value}
}

func (pn PhoneNumber) String() string {
	return pn.value
}

func (pn PhoneNumber) IsValid() bool {
	return len(pn.value) >= 11
}

type Password struct {
	hash string
}

func NewPassword(hash string) Password {
	return Password{hash: hash}
}

func (p Password) Hash() string {
	return p.hash
}

func (p Password) IsEmpty() bool {
	return p.hash == ""
}
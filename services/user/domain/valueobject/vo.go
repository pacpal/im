package valueobject

import (
	"regexp"
	"strings"
)

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
	return id.value != "" && strings.HasPrefix(id.value, "user_")
}

type PhoneNumber struct {
	value string
}

var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

func NewPhoneNumber(value string) PhoneNumber {
	return PhoneNumber{value: value}
}

func (pn PhoneNumber) String() string {
	return pn.value
}

func (pn PhoneNumber) IsValid() bool {
	return phoneRegex.MatchString(pn.value)
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

type UserName struct {
	value string
}

func NewUserName(value string) UserName {
	return UserName{value: value}
}

func (n UserName) String() string {
	return n.value
}

func (n UserName) IsValid() bool {
	return len(n.value) >= 2 && len(n.value) <= 50
}

type AvatarURL struct {
	value string
}

func NewAvatarURL(value string) AvatarURL {
	return AvatarURL{value: value}
}

func (a AvatarURL) String() string {
	return a.value
}

func (a AvatarURL) IsEmpty() bool {
	return a.value == ""
}

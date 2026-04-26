package users

import "time"

// Shared pagination types

type Meta struct {
	Paging *Paging `json:"paging,omitempty"`
}

type Paging struct {
	Total      int    `json:"total,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`
}

type Links struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Last  string `json:"last,omitempty"`
}

type ResourceLinks struct {
	Self string `json:"self,omitempty"`
}

// UsersResponse is the response for a list of users.
type UsersResponse struct {
	Data  []User `json:"data"`
	Links *Links `json:"links,omitempty"`
	Meta  *Meta  `json:"meta,omitempty"`
}

// UserResponse is the response for a single user.
type UserResponse struct {
	Data  User   `json:"data"`
	Links *Links `json:"links,omitempty"`
}

// User represents a user resource.
type User struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes *UserAttributes `json:"attributes,omitempty"`
	Links      *ResourceLinks  `json:"links,omitempty"`
}

// UserAttributes contains the attributes of a user.
type UserAttributes struct {
	FirstName           string        `json:"firstName,omitempty"`
	LastName            string        `json:"lastName,omitempty"`
	MiddleName          string        `json:"middleName,omitempty"`
	Status              string        `json:"status,omitempty"`
	ManagedAppleAccount string        `json:"managedAppleAccount,omitempty"`
	IsExternalUser      bool          `json:"isExternalUser,omitempty"`
	RoleOuList          []RoleOu      `json:"roleOuList,omitempty"`
	Email               string        `json:"email,omitempty"`
	EmployeeNumber      string        `json:"employeeNumber,omitempty"`
	CostCenter          string        `json:"costCenter,omitempty"`
	Division            string        `json:"division,omitempty"`
	Department          string        `json:"department,omitempty"`
	JobTitle            string        `json:"jobTitle,omitempty"`
	StartDateTime       *time.Time    `json:"startDateTime,omitempty"`
	CreatedDateTime     *time.Time    `json:"createdDateTime,omitempty"`
	UpdatedDateTime     *time.Time    `json:"updatedDateTime,omitempty"`
	PhoneNumbers        []PhoneNumber `json:"phoneNumbers,omitempty"`
}

// RoleOu represents a role and organizational unit assignment.
type RoleOu struct {
	RoleName string `json:"roleName,omitempty"`
	OuId     string `json:"ouId,omitempty"`
}

// PhoneNumber represents a phone number with its type.
type PhoneNumber struct {
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Type        string `json:"type,omitempty"`
}

// RequestQueryOptions represents query parameters for user endpoints.
type RequestQueryOptions struct {
	// Fields specifies which fields to return. Use Field* constants.
	Fields []string
	// Limit is the number of resources to return (max 1000).
	Limit int
}

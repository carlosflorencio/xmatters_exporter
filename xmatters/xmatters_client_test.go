package xmatters

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	validGroupMembersPayload = []byte(`{
  "count": 1,
  "total": 1,
  "data": [
    {
      "group": {
        "id": "a6f7a219-8ee2-4462-9eb4-92b014e1c091",
        "targetName": "Dummy Group Name",
        "recipientType": "GROUP",
        "links": {
          "self": "/api/xm/1/groups/a6f7a219-8ee2-4462-9eb4-92b014e1c091"
        }
      },
      "member": {
        "id": "25a7c82f-d862-41d3-98ef-4c7a8611f511",
        "targetName": "userid",
        "firstName": "John",
        "lastName": "Doe",
        "recipientType": "PERSON",
        "links": {
          "self": "/api/xm/1/people/25a7c82f-d862-41d3-98ef-4c7a8611f511"
        }
      }
    }
  ],
  "links": {
    "self": "/api/xm/1/groups/a6f7a219-8ee2-4462-9eb4-92b014e1c091/members?offset=0&limit=100"
  }
}`)
	wrongTokenPayload = []byte(`{
  "code": 401,
  "message": "Invalid Credentials",
  "reason": "Unauthorized"
}`)
)

func TestGroupMembers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(validGroupMembersPayload)
	}))
	defer server.Close()

	xmatters := NewXmattersClient(server.URL, "test")
	group, err := xmatters.GetGroupRoster("group")

	assert.Nil(t, err)
	assert.Equal(t, "Dummy Group Name", group.Data[0].Group.TargetName)
	assert.Equal(t, "John", group.Data[0].Member.FirstName)
	assert.Equal(t, 1, len(group.Data))
}

func TestWrongToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write(wrongTokenPayload)
	}))
	defer server.Close()

	xmatters := NewXmattersClient(server.URL, "test")
	_, err := xmatters.GetGroupRoster("group")

	assert.NotNil(t, err)
}
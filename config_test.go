package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfig(t *testing.T) {
	servers, err := readHostConfig("")
	fmt.Println(servers, err)
}

func TestParseConfig(t *testing.T) {
	type test struct {
		mock   func() (string, string)
		expect []string
		err    error
	}

	cases := []test{
		{
			mock: func() (string, string) {
				raw := `
				root@123.456
				jeki@109.109
				`

				return raw, ""
			},
			expect: []string{
				"root@123.456",
				"jeki@109.109",
			},
			err: nil,
		},
		{
			mock: func() (string, string) {
				raw := `
					[my-server1]
					root@123.124
					#something
					root@123.123
				`
				return raw, "my-server1"
			},
			expect: []string{
				"root@123.124",
				"root@123.123",
			},
		},
		{
			mock: func() (string, string) {
				raw := `
					[server1]
					root@123.124
					root@123.123
				`
				return raw, "server"
			},
			err: fmt.Errorf("Group 'server' not found"),
		}, {
			mock: func() (string, string) {
				raw := `
					[server1]
					root@123.124
					root@123.123
				`
				return raw, ""
			},
			err: fmt.Errorf("Config use group mode but no group specified"),
		}, {
			mock: func() (string, string) {
				raw := `
					[server1]
					root@123.124
					root@123.123
				`
				return raw, "server2"
			},
			err: fmt.Errorf("Group 'server2' not found"),
		},
	}

	for _, tc := range cases {
		hosts, err := parseConfig(tc.mock())
		assert.Equal(t, tc.err, err)
		assert.Equal(t, tc.expect, hosts)
	}
}

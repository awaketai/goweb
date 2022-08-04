package framework

import (
	"strings"
	"testing"
)

// 1./user/login
// 2./user/logout
// 3./subject/name
// 4./subject/name/age
// 5/subject/:id/name
func TestFilterChildNodes(t *testing.T) {
	root := &node{
		isLast:  false,
		segment: "/",
		handler: func(c *Context) error { return nil },
		childs: []*node{
			{
				isLast:  false,
				segment: "USER",
				handler: func(c *Context) error { return nil },
				childs: []*node{
					{
						isLast:  true,
						segment: "LOGIN",
						handler: func(c *Context) error { return nil },
					},
					{
						isLast:  true,
						segment: "LOGOUT",
						handler: func(c *Context) error { return nil },
					},
				},
			},
			{
				isLast:  false,
				segment: "SUBJECT",
				handler: func(c *Context) error { return nil },
				childs: []*node{
					{
						isLast:  false,
						segment: "NAME",
						handler: func(c *Context) error { return nil },
						childs: []*node{
							{
								isLast:  true,
								segment: "AGE",
								handler: func(c *Context) error { return nil },
							},
						},
					},
					{
						isLast:  false,
						segment: "SUBJECT",
						handler: func(c *Context) error { return nil },
						childs: []*node{
							{
								isLast:  true,
								segment: "NAME",
								handler: func(c *Context) error { return nil },
							},
						},
					},
				},
			},
		},
	}

	{
		nodes := root.filterChildNodes("/")
		for _, v := range nodes {
			userNodes := v.filterChildNodes("user")
			if len(userNodes) != 2 {
				t.Errorf("expected node length 2,actual:%d", len(nodes))
			}

			for _, v := range userNodes {
				if !strings.Contains(v.segment, "login") || !strings.Contains(v.segment, "logout") {
					t.Errorf("expected segment is login or logout,actual not contain")
				}
			}
		}

	}
}

func TestMatchNode(t *testing.T) {
	root := &node{
		isLast:  false,
		segment: "",
		handler: func(c *Context) error { return nil },
		childs: []*node{
			{
				isLast:  false,
				segment: "USER",
				handler: func(c *Context) error { return nil },
				childs: []*node{
					{
						isLast:  true,
						segment: "LOGIN",
						handler: func(c *Context) error { return nil },
					},
					{
						isLast:  true,
						segment: "LOGOUT",
						handler: func(c *Context) error { return nil },
					},
				},
			},
			{
				isLast:  false,
				segment: "SUBJECT",
				handler: func(c *Context) error { return nil },
				childs: []*node{
					{
						isLast:  false,
						segment: "NAME",
						handler: func(c *Context) error { return nil },
						childs: []*node{
							{
								isLast:  true,
								segment: "AGE",
								handler: func(c *Context) error { return nil },
							},
						},
					},
					{
						isLast:  false,
						segment: "NAME",
						handler: func(c *Context) error { return nil },
					},
				},
			},
		},
	}

	{
		nodes := root.matchNode("user/login")
		if nodes == nil {
			t.Errorf("expected not nil,actual nil")
		}
	}

	{
		nodes := root.matchNode("subject/name/age")
		if nodes == nil {
			t.Errorf("expected not nil,actual nil")
		}
	}
}

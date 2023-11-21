package goweb

import (
	"context"
	"testing"
)

func Test_filterChildNodes(t *testing.T) {
	root := &node{
		isLast:  false,
		segment: "",
		handler: func(c context.Context) error { return nil },
		childs: []*node{
			{
				isLast:  true,
				segment: "FOO",
				handler: func(c context.Context) error { return nil },
				childs:  nil,
			},
			{
				isLast:  true,
				segment: ":id",
				handler: nil,
				childs:  nil,
			},
		},
	}
	{
		nodes := root.filterChildNodes("FOO")
		if len(nodes) != 2 {
			t.Error("foo error")
		}
	}

	{

		nodes := root.filterChildNodes(":foo")
		if len(nodes) != 2 {
			t.Error(":foo error")
		}
	}
}

func Text_matchNode(t *testing.T) {
	root := &node{
		isLast:  false,
		segment: "",
		handler: func(c context.Context) error { return nil },
		childs: []*node{
			{
				isLast:  true,
				segment: "FOO",
				handler: func(c context.Context) error { return nil },
				childs: []*node{
					&node{
						isLast:  true,
						segment: "BAR",
						handler: func(c context.Context) error { panic("not implemented") },
						childs:  []*node{},
					},
				},
			},
			{
				isLast:  true,
				segment: ":id",
				handler: nil,
				childs:  nil,
			},
		},
	}
	{
		node := root.matchNode("foo/bar")
		if node == nil {
			t.Error("match mornal node error")
		}
	}

	{
		node := root.matchNode("test")
		if node == nil {
			t.Error("match test")
		}
	}
}

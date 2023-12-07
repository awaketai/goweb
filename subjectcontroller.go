package goweb

import "github.com/awaketai/goweb/framework2"

func SubjectAddController(c *framework2.Context) error {
	c.JSON(200,"ok,SubjectAddController")
	return nil
}

func SubjectListController(c *framework2.Context) error {
	c.JSON(200,"ok,SubjectListController")
	return nil
}

func SubjectUpdateController(c *framework2.Context) error {
	c.JSON(200,"ok,SubjectUpdateController")
	return nil
}

func SubjectDeleteController(c *framework2.Context) error {
	c.JSON(200,"ok,SubjectDeleteController")
	return nil
}

func SubjectGetController(c *framework2.Context) error {
	c.JSON(200,"ok,SubjectGetController")
	return nil
}
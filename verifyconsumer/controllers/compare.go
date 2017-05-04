package controllers

import (
	_ "fmt"
	"github.com/astaxie/beego"
	"verifyconsumer/models"
)

type ComparedController struct {
	beego.Controller
}

type Result struct {
	Detail string `json:"detail"`
	Status string `json:"status"`
}

type ResultT struct {
	Confidence float64 `json:"confidence"`
	Status     string  `json:"status"`
}

// URLMapping ...
func (c *ComparedController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// Post ...
// @Title Post
// @Description upload files and compare
// @Success 200 {int} confidence
// @Failure 403 body is empty
// @router / [post]
func (c *ComparedController) Post() {
	f1, _, err := c.GetFile("file1")
	if err != nil {
		c.Data["json"] = Result{err.Error(), "fail"}
		c.ServeJSON()
	}
	f2, _, err := c.GetFile("file2")
	if err != nil {
		c.Data["json"] = Result{err.Error(), "fail"}
		c.ServeJSON()
	}

	conf, _ := models.Comparing(f1, f2)
	f1.Close()
	f2.Close()

	c.Data["json"] = ResultT{conf, "success"}
	c.ServeJSON()
}

// Get ...
// @Title Get
// @Description get restapi
// @Success 200
// @Failure 403
// @router / [get]
func (c *ComparedController) Get() {

	c.Data["json"] = 222222222
	c.ServeJSON()

}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmoule/go2hal/hal"
)

func RetrieveRootResources(ctx *gin.Context) {
	root := hal.NewResourceObject()
	link := &hal.LinkObject{Href: "http://localhost:8080/v1/"}

	self, _ := hal.NewLinkRelation("self") //skipped error handling
	self.SetLink(link)

	root.AddLink(self)

	encoder := hal.NewEncoder()
	bytes, error := encoder.ToJSON(root)
	if error != nil {
		//TODO Implements in future
		return
	}

	ctx.Data(http.StatusOK, "application/json", bytes)
}

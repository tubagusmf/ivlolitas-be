// @title IVLOLITAS API
// @version 1.0
// @description REST API for IVLOLITAS Ecommerce
// @BasePath /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	_ "github.com/tubagusmf/ivlolitas-be/docs"
	"github.com/tubagusmf/ivlolitas-be/internal/console"
)

func main() {
	console.Execute()
}

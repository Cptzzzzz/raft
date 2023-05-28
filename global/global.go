package global

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var Me int
var Router *gin.Engine
var HttpClient *http.Client
var Peers []string
var JudgeHost string

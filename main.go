package main

import (
	"github.com/bmizerany/pat"
	"github.com/macococo/go-webbase/conf"
	"github.com/macococo/go-zipcode/controllers"
	"github.com/macococo/go-zipcode/tasks"
	"github.com/yvasiyarov/go-metrics"
	"github.com/yvasiyarov/gorelic"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

var (
	agent *gorelic.Agent
)

func initGorelic() {
	if !conf.IsDev() {
		return
	}

	agent = gorelic.NewAgent()
	agent.NewrelicName = "go-zipcode"
	agent.Verbose = true
	agent.NewrelicLicense = conf.Config.NewrelicLicense
	agent.CollectHTTPStat = true
	agent.HTTPTimer = metrics.NewTimer()
	agent.Run()
}

func wrapController(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	if agent == nil {
		return f
	}
	return agent.WrapHTTPHandlerFunc(f)
}

func initRouter() {
	port := conf.Config.Port

	m := pat.New()
	m.Get("/api/search", http.HandlerFunc(wrapController(controllers.SearchController)))
	http.Handle("/", m)

	log.Println("HTTP listen port:", port)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func main() {
	log.Println("Runmode:", conf.Config.Runmode)
	log.Println("NumCPU:", runtime.NumCPU())
	log.Println("GOMAXPROCS:", runtime.GOMAXPROCS(runtime.NumCPU()))

	if len(os.Args) >= 2 && os.Args[1] == "ImportAll" {
		tasks.ImportAll()
	} else {
		initGorelic()
		initRouter()
	}
}

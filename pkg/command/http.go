package command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/gregbiv/sandbox/pkg/api/docs"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// HTTPCommand is responsible for running the http server
type HTTPCommand struct {
	Meta
}

// Run is responsible for starting the http server
func (c *HTTPCommand) Run(args []string) int {
	flags := c.FlagSet("http")
	flags.Usage = func() { c.UI.Output(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	router := chi.NewRouter()

	// A good base middleware stack
	router.Use(
		chiMiddleware.WithValue("app.config", c.Config),
		chiMiddleware.Recoverer,
		lg.RequestLogger(logrus.StandardLogger()),
	)

	// HelloWorld
	router.Route("/", func(r chi.Router) {
		r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("Sandbox API"))
		})
	})

	// Documentation
	docs.DocServer(router, docs.DocsPATH, http.Dir(docs.DocsDIR))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Config.Port), router))

	return 0
}

// Help outputs a helper text for the command
func (*HTTPCommand) Help() string {
	helpText := `
Usage: sandbox http [options]

  Start the Http Rest API server
`

	return strings.TrimSpace(helpText)
}

// Synopsis outputs to the console the synopsis of the command
func (c *HTTPCommand) Synopsis() string {
	return "Start the Http Rest API server"
}

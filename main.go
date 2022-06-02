package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Charl88/hayvnapi/routes"
	"github.com/Charl88/hayvnapi/schedulers"
	"github.com/Charl88/hayvnapi/shared"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/swaggest/rest"
	"github.com/swaggest/rest/chirouter"
	"github.com/swaggest/rest/jsonschema"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/openapi"
	"github.com/swaggest/rest/request"
	"github.com/swaggest/rest/response"
	"github.com/swaggest/rest/response/gzip"
	"github.com/swaggest/swgui/v3cdn"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func main() {

	log.Print("Starting API")
	// Init API documentation schema.
	apiSchema := &openapi.Collector{}
	apiSchema.Reflector().SpecEns().Info.Title = "HAYVN API"
	apiSchema.Reflector().SpecEns().Info.WithDescription("API for receiving and sending batched messages")
	apiSchema.Reflector().SpecEns().Info.Version = "v1.0.0"

	r := chirouter.NewWrapper(chi.NewRouter())

	// Setup request decoder and validator
	validatorFactory := jsonschema.NewFactory(apiSchema, apiSchema)
	decoderFactory := request.NewDecoderFactory()
	decoderFactory.ApplyDefaults = true
	decoderFactory.SetDecoderFunc(rest.ParamInPath, chirouter.PathToURLValues)

	// Build CORS options
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}

	// Setup middlewares.
	r.Use(
		middleware.Logger,                             // Logger
		middleware.Recoverer,                          // Panic recovery.
		nethttp.OpenAPIMiddleware(apiSchema),          // Documentation collector.
		request.DecoderMiddleware(decoderFactory),     // Request decoder setup.
		request.ValidatorMiddleware(validatorFactory), // Request validator setup.
		response.EncoderMiddleware,                    // Response encoder setup.
		gzip.Middleware,                               // Response compression with support for direct gzip pass through.
		cors.Handler(corsOptions),
	)

	r.Method(http.MethodGet, "/", nethttp.NewHandler(HealthCheck()))
	r.Method(http.MethodGet, "/docs/openapi.json", apiSchema)
	r.Mount("/docs", v3cdn.NewHandler(apiSchema.Reflector().Spec.Info.Title, "/docs/openapi.json", "/docs"))

	// Initialize the empty array to store the incoming messages in memory
	messageQueue := shared.MessageArray{
		Messages: []shared.Message{},
	}

	r.Method(http.MethodPost, "/message", nethttp.NewHandler(routes.MessageReceiver(&messageQueue)))
	r.Method(http.MethodPost, "/aggregated-messages", nethttp.NewHandler(routes.AggregatedMessages()))

	schedulers.CreateMessageScheduler(&messageQueue)

	http.ListenAndServe(":3000", r)
}

// A health check endpoint for the API to report it's current status
func HealthCheck() usecase.IOInteractor {
	type healthCheckOutput struct {
		Status string `json:"status"`
		Data   string `json:"data"`
	}

	u := usecase.NewIOI(nil, new(healthCheckOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			out = output.(*healthCheckOutput)
		)
		out.Data = "Server is up and running"
		out.Status = status.OK.String()
		return nil
	})
	u.SetTitle("Health Check")

	return u
}

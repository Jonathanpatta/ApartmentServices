package Router

import (
	"github.com/gorilla/mux"
	"github.com/jonathanpatta/apartmentservices/Consumers"
	"github.com/jonathanpatta/apartmentservices/Items"
	"github.com/jonathanpatta/apartmentservices/Middleware"
	"github.com/jonathanpatta/apartmentservices/Orders"
	"github.com/jonathanpatta/apartmentservices/Producers"
	"github.com/jonathanpatta/apartmentservices/Services"
	"github.com/jonathanpatta/apartmentservices/Settings"
	"github.com/jonathanpatta/apartmentservices/Subscriptions"
	"log"
)

func GetMainRouter() *mux.Router {
	settings, err := Settings.NewSettings()
	if err != nil {
		log.Fatalf("unable to load settings, %v", err)
	}

	router := mux.NewRouter()
	Consumers.AddSubrouter(router, settings)
	Producers.AddSubrouter(router, settings)
	Services.AddSubrouter(router, settings)
	Items.AddSubrouter(router, settings)
	Orders.AddSubrouter(router, settings)
	Subscriptions.AddSubrouter(router, settings)
	router.Use(Middleware.CorsMiddleware)

	return router
}

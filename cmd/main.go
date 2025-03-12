package main

import (
	"fmt"
	"linkshorter/configs"
	"linkshorter/internal/auth"
	"linkshorter/internal/link"
	"linkshorter/internal/stat"
	"linkshorter/internal/user"
	"linkshorter/pkg/db"
	"linkshorter/pkg/event"
	"linkshorter/pkg/middleware"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	db := db.NewDB(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	linkRepository := link.NewLinkRepository(db)
	userRpository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	authService := auth.NewAuthService(userRpository)
	statsServces := stat.NewStatService(&stat.StatServiceDeps{
		EventBus: eventBus,
		StatRepository: statRepository,
	})

	auth.NewAuthHundler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		EventBus:       eventBus,
		Config:         conf,
	})
	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepository,
		Config: conf,
	})

	go statsServces.AddClick()

	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	return stack(router)
}

func main() {
	app := App()

	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	fmt.Println("Server is listeing on port 8081")
	server.ListenAndServe()
}

package main

import (
	"harmonica/config"
	"harmonica/internal/handler"
	"harmonica/internal/handler/middleware"
	auth "harmonica/internal/microservices/auth/proto"
	image "harmonica/internal/microservices/image/proto"
	like "harmonica/internal/microservices/like/proto"

	"harmonica/internal/repository"
	"harmonica/internal/service"
	"log"
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/joho/godotenv"
	v3 "github.com/swaggest/swgui/v3"
)

func runServer(addr string) {
	logger := config.ConfigureZapLogger("monolit")
	defer logger.Sync()

	authCli, imageCli, likeCli := makeMicroservicesClients()

	conf := config.New()
	connector, err := repository.NewConnector(conf, imageCli)
	if err != nil {
		logger.Info(err.Error())
		return
	}
	defer connector.Disconnect()

	r := repository.NewRepository(connector, logger)
	s := service.NewService(r, likeCli)

	hub := handler.NewHub() // ws-server
	h := handler.NewAPIHandler(s, logger, hub, authCli, imageCli, likeCli)
	mux := http.NewServeMux()

	configureUserRoutes(logger, h, mux)
	configurePinRoutes(logger, h, mux)
	configureBoardRoutes(logger, h, mux)
	configureChatRoutes(logger, h, mux)
	configureDraftRoutes(logger, h, mux)
	configureSearchRoutes(logger, h, mux)
	configureSubscriptionRoutes(logger, h, mux)
	configureNotificationRoutes(logger, h, mux)
	configureCommentsRoutes(logger, h, mux)

	mux.Handle("GET /docs/swagger.json", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))
	mux.Handle("GET /swagger/", v3.NewHandler("My API", "/docs/swagger.json", "/swagger"))
	mux.HandleFunc("GET /img/{image_name}", h.GetImage)

	go hub.Run()
	mux.HandleFunc("GET /ws", middleware.AuthRequired(logger, h.AuthService, h.ServeWs))

	loggedMux := middleware.Logging(logger, mux)

	server := http.Server{
		Addr: addr,
	}

	if config.GetEnvAsBool("DEBUG", false) {
		server.Handler = middleware.CORS(loggedMux)
		server.ListenAndServe()
		return
	}
	server.Handler = middleware.CSRF(middleware.CORS(loggedMux))
	server.ListenAndServeTLS("/etc/letsencrypt/live/harmoniums.ru/fullchain.pem", "/etc/letsencrypt/live/harmoniums.ru/privkey.pem")
}

func makeMicroservicesClients() (auth.AuthorizationClient, image.ImageClient, like.LikeClient) {
	authConn, err := grpc.Dial(config.GetEnv("AUTH_MICROSERVICE_PORT", ":8002"), grpc.WithInsecure())
	if err != nil {
		log.Print(err)
		return nil, nil, nil
	}

	imageConn, err := grpc.Dial(config.GetEnv("IMAGE_MICROSERVICE_PORT", ":8003"), grpc.WithInsecure())
	if err != nil {
		log.Print(err)
		return nil, nil, nil
	}

	likeConn, err := grpc.Dial(config.GetEnv("LIKE_MICROSERVICE_PORT", ":8004"), grpc.WithInsecure())
	if err != nil {
		log.Print(err)
		return nil, nil, nil
	}
	authCli := auth.NewAuthorizationClient(authConn)
	imageCli := image.NewImageClient(imageConn)
	likeCli := like.NewLikeClient(likeConn)
	return authCli, imageCli, likeCli
}

func configureUserRoutes(logger *zap.Logger, h *handler.APIHandler, mux *http.ServeMux) {
	authRoutes := map[string]http.HandlerFunc{
		"POST /api/v1/users/{user_id}": h.UpdateUser,
	}
	notAuthRoutes := map[string]http.HandlerFunc{
		"POST /api/v1/login": h.Login,
		"POST /api/v1/users": h.Register,
	}
	checkAuthRoutes := map[string]http.HandlerFunc{
		"GET /api/v1/is_auth":          h.IsAuth, // check it
		"GET /api/v1/logout":           h.Logout,
		"GET /api/v1/users/{nickname}": h.GetUser,
		"GET /api/v1/all/users":        h.GetAllUsers,
	}
	for pattern, f := range authRoutes {
		mux.HandleFunc(pattern, middleware.AuthRequired(logger, h.AuthService, f))
	}
	for pattern, f := range notAuthRoutes {
		mux.HandleFunc(pattern, middleware.NoAuthRequired(logger, h.AuthService, f))
	}
	for pattern, f := range checkAuthRoutes {
		mux.HandleFunc(pattern, middleware.CheckAuth(logger, h.AuthService, f))
	}
}

func configurePinRoutes(logger *zap.Logger, h *handler.APIHandler, mux *http.ServeMux) {
	authRoutes := map[string]http.HandlerFunc{
		"POST /api/v1/pins":                 h.CreatePin,
		"POST /api/v1/pins/{pin_id}":        h.UpdatePin,
		"DELETE /api/v1/pins/{pin_id}":      h.DeletePin,
		"POST /api/v1/pins/{pin_id}/like":   h.CreateLike,
		"DELETE /api/v1/pins/{pin_id}/like": h.DeleteLike,
		"GET /api/v1/favorites":             h.GetFavorites,
	}
	checkAuthRoutes := map[string]http.HandlerFunc{
		"GET /api/v1/pins/{pin_id}": h.GetPin,
		"GET /api/v1/pins":          h.Feed,
	}
	publicRoutes := map[string]http.HandlerFunc{
		//"GET /api/v1/pins":                    h.Feed,
		"GET /api/v1/pins/created/{nickname}": h.UserPins,
		"GET /api/v1/likes/{pin_id}/users":    h.UsersLiked,
	}
	for pattern, f := range authRoutes {
		mux.HandleFunc(pattern, middleware.AuthRequired(logger, h.AuthService, f))
	}
	for pattern, f := range checkAuthRoutes {
		mux.HandleFunc(pattern, middleware.CheckAuth(logger, h.AuthService, f))
	}
	for pattern, f := range publicRoutes {
		mux.HandleFunc(pattern, f)
	}
}

func configureBoardRoutes(logger *zap.Logger, h *handler.APIHandler, mux *http.ServeMux) {
	authRoutes := map[string]http.HandlerFunc{
		"POST /api/v1/boards":                            h.CreateBoard,
		"POST /api/v1/boards/{board_id}":                 h.UpdateBoard,
		"DELETE /api/v1/boards/{board_id}":               h.DeleteBoard,
		"POST /api/v1/boards/{board_id}/pins/{pin_id}":   h.AddPinToBoard,
		"DELETE /api/v1/boards/{board_id}/pins/{pin_id}": h.DeletePinFromBoard,
		"GET /api/v1/boards/excluding/{pin_id}":          h.GetUserBoardsWithoutPin,
	}
	checkAuthRoutes := map[string]http.HandlerFunc{
		"GET /api/v1/boards/{board_id}":         h.GetBoard,
		"GET /api/v1/boards/created/{nickname}": h.GetUserBoards,
	}
	for pattern, f := range authRoutes {
		mux.HandleFunc(pattern, middleware.AuthRequired(logger, h.AuthService, f))
	}
	for pattern, f := range checkAuthRoutes {
		mux.HandleFunc(pattern, middleware.CheckAuth(logger, h.AuthService, f))
	}
}

func configureChatRoutes(logger *zap.Logger, h *handler.APIHandler, mux *http.ServeMux) {
	authRoutes := map[string]http.HandlerFunc{
		"POST /api/v1/messages/{receiver_id}": h.SendMessage,
		"GET /api/v1/messages/{user_id}":      h.ReadMessages,
		"GET /api/v1/chats":                   h.GetUserChats,
	}
	for pattern, f := range authRoutes {
		mux.HandleFunc(pattern, middleware.AuthRequired(logger, h.AuthService, f))
	}
}

func configureDraftRoutes(logger *zap.Logger, h *handler.APIHandler, mux *http.ServeMux) {
	authRoutes := map[string]http.HandlerFunc{
		"POST /api/v1/drafts/{receiver_id}": h.UpdateDraft,
	}
	for pattern, f := range authRoutes {
		mux.HandleFunc(pattern, middleware.AuthRequired(logger, h.AuthService, f))
	}
}

func configureSearchRoutes(logger *zap.Logger, h *handler.APIHandler, mux *http.ServeMux) {
	checkAuthRoutes := map[string]http.HandlerFunc{
		"GET /api/v1/search/{search_query}": h.Search,
	}
	for pattern, f := range checkAuthRoutes {
		mux.HandleFunc(pattern, middleware.CheckAuth(logger, h.AuthService, f))
	}
}

func configureSubscriptionRoutes(logger *zap.Logger, h *handler.APIHandler, mux *http.ServeMux) {
	authRoutes := map[string]http.HandlerFunc{
		"POST /api/v1/users/subscribe/{user_id}":   h.SubscribeToUser,
		"DELETE /api/v1/users/subscribe/{user_id}": h.UnsubscribeFromUser,
	}
	publicRoutes := map[string]http.HandlerFunc{
		"GET /api/v1/users/subscribers/{user_id}":   h.GetUserSubscribers,
		"GET /api/v1/users/subscriptions/{user_id}": h.GetUserSubscriptions,
	}
	for pattern, f := range authRoutes {
		mux.HandleFunc(pattern, middleware.AuthRequired(logger, h.AuthService, f))
	}
	for pattern, f := range publicRoutes {
		mux.HandleFunc(pattern, f)
	}
}

func configureNotificationRoutes(logger *zap.Logger, h *handler.APIHandler, mux *http.ServeMux) {
	authRoutes := map[string]http.HandlerFunc{
		"GET /api/v1/notifications":                         h.GetUnreadNotifications,
		"POST /api/v1/notifications/read/{notification_id}": h.ReadNotification,
		"POST /api/v1/notifications/read/all":               h.ReadAllNotifications,
	}
	for pattern, f := range authRoutes {
		mux.HandleFunc(pattern, middleware.AuthRequired(logger, h.AuthService, f))
	}
}

func configureCommentsRoutes(logger *zap.Logger, h *handler.APIHandler, mux *http.ServeMux) {
	authRoutes := map[string]http.HandlerFunc{
		"POST /api/v1/pin/comments/{pin_id}": h.AddComment,
	}
	checkAuthRoutes := map[string]http.HandlerFunc{
		"GET /api/v1/pin/comments/{pin_id}": h.GetComments,
	}
	for pattern, f := range authRoutes {
		mux.HandleFunc(pattern, middleware.AuthRequired(logger, h.AuthService, f))
	}
	for pattern, f := range checkAuthRoutes {
		mux.HandleFunc(pattern, middleware.CheckAuth(logger, h.AuthService, f))
	}
}

func init() {
	if err := godotenv.Load("conf.env"); err != nil {
		log.Print("No conf.env file found")
	}
}

// @title			Harmonium backend API
// @version		1.0
// @description	This is API-docs of backend server of Harmonica team.
// @host			https://harmoniums.ru
// @BasePath		api/v1
func main() {
	runServer(":8080")
}

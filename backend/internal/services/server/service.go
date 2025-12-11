package server

import (
	"context"
	"database/sql"
	"fmt"

	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/config"
	handlersaccounts "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/handlers/accounts"
	handlersgames "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/handlers/games"
	handlersparticipants "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/handlers/participants"
	handlersrandom "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/handlers/random"
	handlersrooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/handlers/rooms"
	handlersvotes "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/handlers/votes"
	middlewares "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/middlewares"
	repositorygames "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/games"
	repositoryparticipants "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/participants"
	repositoryrefreshtokens "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/refresh_tokens"
	repositoryresults "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/results"
	repositoryrooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/rooms"
	repositoryusers "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/users"
	repositoryvotes "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/votes"
	servicegames "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/games"
	serviceparticipants "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/participants"
	serviceresults "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/results"
	servicerooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/rooms"
	servicetokens "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/tokens"
	serviceusers "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/users"
	servicevotes "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/votes"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/hub"
)

type Server interface {
	Start() error
}

type Service struct {
	config config.AppConfig

	// repositories
	userRepo         repositoryusers.UserRepository
	refreshTokenRepo repositoryrefreshtokens.RefreshTokenRepository
	gamesRepo        repositorygames.GameRepository
	participantsRepo repositoryparticipants.ParticipantRepository
	resultsRepo      repositoryresults.ResultRepository
	roomsRepo        repositoryrooms.RoomRepository
	votesRepo        repositoryvotes.VoteRepository

	// servicess
	userService        serviceusers.UserService
	tokenService       servicetokens.TokenService
	gameService        servicegames.GameService
	participantService serviceparticipants.ParticipantService
	resultService      serviceresults.ResultService
	roomService        servicerooms.RoomService
	voteService        servicevotes.VoteService

	// handlers
	// accounts handlers
	signUpHandler     handlersaccounts.SignUpHandler
	signInHandler     handlersaccounts.SignInHandler
	refreshHandler    handlersaccounts.RefreshHandler
	logoutHandler     handlersaccounts.LogoutHandler
	userHandler       handlersaccounts.UserHandler
	getByNameHandler  handlersaccounts.GetByNameHandler
	updateUserHandler handlersaccounts.UpdateUserHandler

	// games handlers
	addGameHandler    handlersgames.AddGameHandler
	getGamesHandler   handlersgames.GetGamesHandler
	deleteGameHandler handlersgames.DeleteGameHandler

	// participants handlers
	inviteHandler            handlersparticipants.InviteHandler
	getParticipantsHandler   handlersparticipants.GetParticipantsHandler
	deleteParticipantHandler handlersparticipants.DeleteParticipantHandler

	// random handlers
	getRandomHandler  handlersrandom.GetRandomHandler
	getLastHandler    handlersrandom.GetLastHandler
	getHistoryHandler handlersrandom.GetHistoryHandler

	// rooms handlers
	createRoomHandler  handlersrooms.CreateRoomHandler
	getAllRoomsHandler handlersrooms.GetAllRoomsHandler
	getRoomInfoHandler handlersrooms.GetRoomInfoHandler
	updateRoomHandler  handlersrooms.UpdateRoomHandler
	deleteRoomHandler  handlersrooms.DeleteRoomHandler

	// votes handlers
	addVoteHandler    handlersvotes.AddVoteHandler
	getVotesHandler   handlersvotes.GetVotesHandler
	deleteVoteHandler handlersvotes.DeleteVoteHandler

	// realtime
	wsRoomHandler handlersrooms.WSRoomHandler
	hub           hub.Hub

	// middleware
	authMiddleware      middlewares.AuthMiddleware
	checkRoomMiddleware middlewares.CheckRoomMiddleware

	ctx context.Context
	app *fiber.App
	db  *sql.DB
}

func NewServer(app *fiber.App, db *sql.DB, cfg config.AppConfig, ctx context.Context) (*Service, error) {
	userRepo := repositoryusers.NewRepository(db)
	refreshTokenRepo := repositoryrefreshtokens.NewRepository(db)
	gamesRepo := repositorygames.NewRepository(db)
	participantsRepo := repositoryparticipants.NewRepository(db)
	resultsRepo := repositoryresults.NewRepository(db)
	roomsRepo := repositoryrooms.NewRepository(db)
	votesRepo := repositoryvotes.NewRepository(db)

	userService := serviceusers.NewService(userRepo)
	tokenService := servicetokens.NewService(cfg, refreshTokenRepo)
	gameService := servicegames.NewService(gamesRepo)
	participantService := serviceparticipants.NewService(participantsRepo)
	resultService := serviceresults.NewService(resultsRepo)
	roomService := servicerooms.NewService(roomsRepo)
	voteService := servicevotes.NewService(votesRepo)

	// accounts handlers
	signUpHandler := handlersaccounts.NewSignupHandler(tokenService, userService)
	signInHandler := handlersaccounts.NewSignInHandler(tokenService, userService)
	refreshHandler := handlersaccounts.NewRefreshHandler(tokenService, userService)
	logoutHandler := handlersaccounts.NewLogoutHandler(tokenService, userService)
	userHandler := handlersaccounts.NewUserHandler(userService)
	getByNameHandler := handlersaccounts.NewGetByNameHandler(userService)
	updateUserHandler := handlersaccounts.NewUpdateUserHandler(userService)

	// games handlers
	addGameHandler := handlersgames.NewAddGameHandler(gameService, participantService)
	getGamesHandler := handlersgames.NewGetGamesHandler(gameService)
	deleteGameHandler := handlersgames.NewDeleteGameHandler(gameService, participantService)

	// participants handlers
	inviteHandler := handlersparticipants.NewInviteHandler(participantService, userService)
	getParticipantsHandler := handlersparticipants.NewGetParticipantsHandler(participantService)
	deleteParticipantHandler := handlersparticipants.NewDeleteParticipantHandler(participantService)

	// random handlers
	getRandomHandler := handlersrandom.NewGetRandomHandler(resultService)
	getLastHandler := handlersrandom.NewGetLastHandler(resultService)
	getHistoryHandler := handlersrandom.NewGetHistoryHandler(resultService)

	// rooms handlers
	createRoomHandler := handlersrooms.NewCreateRoomHandler(roomService, participantService)
	getAllRoomsHandler := handlersrooms.NewGetAllRoomsHandler(roomService)
	getRoomInfoHandler := handlersrooms.NewGetRoomInfoHandler(roomService)
	updateRoomHandler := handlersrooms.NewUpdateRoomHandler(roomService)
	deleteRoomHandler := handlersrooms.NewDeleteRoomHandler(roomService)

	// votes handlers
	addVoteHandler := handlersvotes.NewAddVoteHandler(voteService)
	getVotesHandler := handlersvotes.NewGetVotesHandler(voteService)
	deleteVoteHandler := handlersvotes.NewDeleteVoteHandler(voteService)

	// realtime hub & handler
	h := hub.NewHubWS()
	wsRoomHandler := handlersrooms.NewWSRoomHandler(h, tokenService)

	// pass hub to services that emit events
	gameService.SetHub(h)
	participantService.SetHub(h)
	voteService.SetHub(h)
	roomService.SetHub(h)

	authMiddleware := middlewares.NewAuthMiddleware(tokenService)
	checkRoomMiddleware := middlewares.NewCheckRoomMiddleware(roomService, participantService)

	return &Service{
		ctx:    ctx,
		config: cfg,
		app:    app,
		db:     db,

		// repositories
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		gamesRepo:        gamesRepo,
		participantsRepo: participantsRepo,
		resultsRepo:      resultsRepo,
		roomsRepo:        roomsRepo,
		votesRepo:        votesRepo,

		// services
		userService:        userService,
		tokenService:       tokenService,
		gameService:        gameService,
		participantService: participantService,
		resultService:      resultService,
		roomService:        roomService,
		voteService:        voteService,

		// handlers
		// accounts handlers
		signUpHandler:     *signUpHandler,
		signInHandler:     *signInHandler,
		refreshHandler:    *refreshHandler,
		logoutHandler:     *logoutHandler,
		userHandler:       *userHandler,
		getByNameHandler:  *getByNameHandler,
		updateUserHandler: *updateUserHandler,

		// games handlers
		addGameHandler:    *addGameHandler,
		getGamesHandler:   *getGamesHandler,
		deleteGameHandler: *deleteGameHandler,

		// participants handlers
		inviteHandler:            *inviteHandler,
		getParticipantsHandler:   *getParticipantsHandler,
		deleteParticipantHandler: *deleteParticipantHandler,

		// random handlers
		getRandomHandler:  *getRandomHandler,
		getLastHandler:    *getLastHandler,
		getHistoryHandler: *getHistoryHandler,

		// rooms handlers
		createRoomHandler:  *createRoomHandler,
		getAllRoomsHandler: *getAllRoomsHandler,
		getRoomInfoHandler: *getRoomInfoHandler,
		updateRoomHandler:  *updateRoomHandler,
		deleteRoomHandler:  *deleteRoomHandler,

		// votes handlers
		addVoteHandler:    *addVoteHandler,
		getVotesHandler:   *getVotesHandler,
		deleteVoteHandler: *deleteVoteHandler,

		// realtime
		wsRoomHandler: *wsRoomHandler,
		hub:           h,

		// middleware
		authMiddleware:      *authMiddleware,
		checkRoomMiddleware: *checkRoomMiddleware,
	}, nil
}

func (s *Service) Start(addr string) error {

	err := s.configure()
	if err != nil {
		return fmt.Errorf("configure error: %v", err)
	}

	logger.Infof(s.ctx, "Starting server on %s", addr)
	return s.app.Listen(addr)
}

func (s *Service) configure() error {
	s.app.Use(middlewares.LogFieldsMiddleware)
	s.app.Get("/api/v1/rooms/:room_id/ws", s.wsRoomHandler.Handle, websocket.New(s.wsRoomHandler.Conn))

	// Public routes
	s.app.Post("/api/auth/signup", s.signUpHandler.HandleSignup)
	s.app.Post("/api/auth/signin", s.signInHandler.HandleSignIn)
	s.app.Post("/api/auth/refresh", s.refreshHandler.HandleRefresh)
	s.app.Post("/api/auth/logout", s.logoutHandler.HandleLogout)

	// Authenticated routes
	authApi := s.app.Group("/api/v1")
	authApi.Use(s.authMiddleware.AuthRequired)

	// Account routes
	authApi.Get("/user", s.userHandler.HandleGetUser)
	authApi.Get("/user/by-name", s.getByNameHandler.HandleGetByName)
	authApi.Put("/user", s.updateUserHandler.HandleUpdateUser)

	// Room routes
	authApi.Post("/rooms", s.createRoomHandler.Handle)
	authApi.Get("/rooms", s.getAllRoomsHandler.HandleGetAllRooms)

	// Room-specific routes (with room middleware)
	roomApi := authApi.Group("/rooms/:room_id")
	roomApi.Use(s.checkRoomMiddleware.Handle)

	// Room info
	roomApi.Get("", s.getRoomInfoHandler.HandleGetRoomInfo)
	roomApi.Put("", s.updateRoomHandler.HandleUpdateRoom)
	roomApi.Delete("", s.deleteRoomHandler.HandleDeleteRoom)

	// Games routes
	roomApi.Post("/games", s.addGameHandler.Handle)
	roomApi.Get("/games", s.getGamesHandler.Handle)
	roomApi.Delete("/games/:game_id", s.deleteGameHandler.Handle)

	// Participants routes
	roomApi.Post("/participants", s.inviteHandler.Handle)
	roomApi.Get("/participants", s.getParticipantsHandler.Handle)
	roomApi.Delete("/participants", s.deleteParticipantHandler.Handle)

	// Votes routes
	roomApi.Post("/votes/", s.addVoteHandler.Handle)
	roomApi.Get("/votes", s.getVotesHandler.Handle)
	roomApi.Delete("/votes/:vote_id", s.deleteVoteHandler.Handle)

	// Random routes
	roomApi.Get("/random", s.getRandomHandler.Handle)
	roomApi.Get("/random/last", s.getLastHandler.Handle)
	roomApi.Get("/random/history", s.getHistoryHandler.Handle)

	// WebSocket route for realtime room updates
	// roomApi.Get("/ws", s.wsRoomHandler.Handle, websocket.New(s.wsRoomHandler.Conn))

	return nil
}

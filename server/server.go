package server

import (
	"errors"
	"github.com/Gewinum/go-df-discord/utils"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"net/http"
)

type Payload struct {
	Error *ErrorInformation `json:"error"`
	Data  interface{}       `json:"data"`
}

func SuccessPayload(data interface{}) Payload {
	return Payload{
		Error: nil,
		Data:  data,
	}
}

func FailurePayload(err ApplicationError) Payload {
	return Payload{
		Error: &ErrorInformation{
			Code:    err.ErrorCode,
			Message: err.Error(),
		},
		Data: nil,
	}
}

type ErrorInformation struct {
	Code    int
	Message string
}

type Server struct {
	accessToken     string
	discordBotToken string
	opts            *Opts
	service         *Service
	bot             *Bot
}

func NewServer(accessToken, discordBotToken string, opts *Opts) *Server {
	FillEmptyOpts(opts)
	service := NewService(opts.Repo, opts.CodeStr)
	bot, err := NewBot(discordBotToken, service)
	if err != nil {
		panic(err)
	}
	return &Server{
		accessToken:     accessToken,
		discordBotToken: discordBotToken,
		opts:            opts,
		service:         service,
		bot:             bot,
	}
}

func (s *Server) Bot() *Bot {
	return s.bot
}

// ServeWeb listens to the specific address
func (s *Server) ServeWeb(addr string) error {
	handler, err := s.GetHttpHandler(false)
	if err != nil {
		return err
	}
	return http.ListenAndServe(addr, handler)
}

// GetHttpHandler returns http.Handler
func (s *Server) GetHttpHandler(isDevelopment bool) (http.Handler, error) {
	if !isDevelopment {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()

	e.Use(sloggin.New(s.opts.Logger))
	e.Use(s.recoveryMiddleware)
	e.Use(s.authMiddleware)

	e.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world!")
	})

	e.POST("/codes/issue", func(c *gin.Context) {
		rawData, err := c.GetRawData()
		utils.ErrorPanic(err)
		xuid := string(rawData)
		info, err := s.service.IssueCode(xuid)
		utils.ErrorPanic(err)
		c.JSON(http.StatusOK, SuccessPayload(info))
	})

	e.POST("/codes/check", func(c *gin.Context) {
		rawData, err := c.GetRawData()
		utils.ErrorPanic(err)
		code := string(rawData)
		info, err := s.service.CheckCode(code)
		utils.ErrorPanic(err)
		c.JSON(http.StatusOK, SuccessPayload(info))
	})

	e.POST("/codes/revoke", func(c *gin.Context) {
		rawData, err := c.GetRawData()
		utils.ErrorPanic(err)
		code := string(rawData)
		utils.ErrorPanic(s.service.RevokeCode(code))
		c.JSON(http.StatusOK, SuccessPayload(nil))
	})

	e.GET("/users/discord/:id", func(c *gin.Context) {
		discordId := c.Param("id")
		if discordId == "" {
			panic(NewApplicationError(40000, "Discord ID is not specified"))
		}
		user, err := s.service.GetUserByDiscord(discordId)
		utils.ErrorPanic(err)
		c.JSON(http.StatusOK, SuccessPayload(user))
	})

	e.GET("/users/xuid/:xuid", func(c *gin.Context) {
		xuid := c.Param("xuid")
		if xuid == "" {
			panic(NewApplicationError(40000, "XUID is not specified"))
		}
		user, err := s.service.GetUserByXUID(xuid)
		utils.ErrorPanic(err)
		c.JSON(http.StatusOK, SuccessPayload(user))
	})

	return e, nil
}

func (s *Server) authMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != s.accessToken {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
}

func (s *Server) recoveryMiddleware(c *gin.Context) {
	defer func() {
		rawErr := recover()
		if rawErr == nil {
			return
		}
		err, ok := rawErr.(error)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var appError ApplicationError
		isAppError := errors.As(err, &appError)
		if !isAppError {
			s.opts.Logger.Error(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		accordingErrorCode := utils.GetNumberFirstDigits(appError.ErrorCode, 3)
		c.AbortWithStatusJSON(accordingErrorCode, FailurePayload(appError))
	}()
	c.Next()
}

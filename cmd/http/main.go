package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	_ "go-tonify-backend/docs"
	"go-tonify-backend/internal/api/middleware"
	"go-tonify-backend/internal/api/route"
	"go-tonify-backend/internal/bootstrap"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/repository"
	"go-tonify-backend/internal/service"
	"go-tonify-backend/pkg/logger"
	"log"
)

//	@title			Swagger Tonify API
//	@version		1.0
//	@description	Server help to interact with telegram mini app Tonify.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	app := bootstrap.App()
	r := gin.Default()

	conn, err := openConnectionToDB(app.DBConfig)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = conn.Close()
	}()
	l := logger.NewLogger(logger.DEV, logger.LevelDebug)
	box := container.NewContainer(l, app, conn)
	accountRepository := repository.NewAccountRepository(conn)
	companyRepository := repository.NewCompanyRepository(conn)
	countryRepository := repository.NewCountryRepository()

	authService := service.NewAuthService(box, accountRepository, companyRepository)
	authMiddleware := middleware.NewAuth(box, authService)
	loggerMiddleware := middleware.NewLogger(box)

	route.Setup(r, box, authService, authMiddleware, loggerMiddleware, accountRepository, countryRepository)

	if err := r.Run(app.Address()); err != nil {
		log.Fatalln(err)
	}
}

func openConnectionToDB(dbConfig bootstrap.DBConfig) (*sql.DB, error) {
	psqlConn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
		dbConfig.Mode,
	)
	conn, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return conn, err
	}
	err = conn.Ping()
	return conn, err
}

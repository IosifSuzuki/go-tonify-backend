package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	_ "go-tonify-backend/docs"
	v1 "go-tonify-backend/internal/api/interface/http"
	"go-tonify-backend/internal/container"
	accountRepository "go-tonify-backend/internal/domain/account/repository"
	accountUsecase "go-tonify-backend/internal/domain/account/usecase"
	countryRepository "go-tonify-backend/internal/domain/country/repository"
	countryUsecase "go-tonify-backend/internal/domain/country/usecase"
	"go-tonify-backend/internal/domain/provider/transaction"
	"go-tonify-backend/internal/infrastructure/config"
	"go-tonify-backend/internal/infrastructure/filestorage/s3"
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
	cont, err := composeContainer()
	if err != nil {
		log.Fatalln("fail to compose container", err)
	}
	defer func() {
		_ = cont.GetDBConnection().Close()
	}()

	fileStorage := s3.NewS3FileStorage(cont)
	transactionProvider := transaction.NewProvider(cont.GetDBConnection())

	accountRep := accountRepository.NewAccount(cont.GetDBConnection())
	attachmentRep := accountRepository.NewAttachment(cont.GetDBConnection())
	_ = accountRepository.NewCompany(cont.GetDBConnection())
	countryRep := countryRepository.NewCountry()

	accountUc := accountUsecase.NewAccount(cont, fileStorage, accountRep, attachmentRep, transactionProvider)
	countryUc := countryUsecase.NewCountry(cont, countryRep)

	handler := v1.NewHandler(cont, accountUc, countryUc)

	if err := handler.Run(); err != nil {
		log.Fatalln("fail to run handler", err)
	}
}

func composeContainer() (container.Container, error) {
	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	l := logger.NewLogger(logger.DEV, logger.Level(conf.Server.LoggerLevel))
	conn, err := openConnectionToDB(conf.PostgreSQL)
	if err != nil {
		return nil, err
	}
	return container.NewContainer(l, conf, conn), nil
}

func openConnectionToDB(dbConfig *config.PostgreSQL) (*sql.DB, error) {
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

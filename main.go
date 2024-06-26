package main

import (
	"context"
	"database/sql"
	"github.com/hibiken/asynq"
	rcors "github.com/rs/cors"
	"github.com/vexsx/Simple-Bank/mail"
	"github.com/vexsx/Simple-Bank/worker"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vexsx/Simple-Bank/api"
	db "github.com/vexsx/Simple-Bank/db/sqlc"
	_ "github.com/vexsx/Simple-Bank/doc/statik"
	"github.com/vexsx/Simple-Bank/gapi"
	"github.com/vexsx/Simple-Bank/pb"
	"github.com/vexsx/Simple-Bank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("error with opening db")
	}

	//for docker base database
	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, store, config)
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)

}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(
	redisOpt asynq.RedisClientOpt,
	store db.Store,
	config util.Config,
) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)

	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	grpcLgger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLgger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create grpc listener")
	}
	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot creat grpc server")
	}
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik FS")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot creat grpc listener")
	}
	log.Info().Msgf("start HTTP gateway server at %s", config.HTTPServerAddress)

	// Configure the CORS middleware
	corsMiddleware := rcors.New(rcors.Options{
		AllowedOrigins:      []string{"*", "http://localhost:4200"},
		AllowedMethods:      []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:      []string{"*", "Origin", "Authorization", "Content-Type", "Access-Control-Allow-Headers"},
		ExposedHeaders:      []string{"*", "Content-Type", "Origin", "Authorization", "Content-Type", "Access-Control-Allow-Headers"},
		OptionsPassthrough:  true,
		AllowCredentials:    true,
		AllowPrivateNetwork: true,
		Debug:               true,
		MaxAge:              int(12 * time.Hour),
	})

	handler := corsMiddleware.Handler(gapi.HttpLogger(mux))
	err = http.Serve(listener, handler)
	if err != nil {
		log.Error().Err(err).Msg("HTTP gateway server failed to serve")
	}
}

func runGinServer(config util.Config, store db.Store) {

	//mewo

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}

}

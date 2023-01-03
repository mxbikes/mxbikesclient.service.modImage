package main

import (
	"net"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mxbikes/mxbikesclient.service.modImage/handler"
	"github.com/mxbikes/mxbikesclient.service.modImage/repository"
	protobuffer "github.com/mxbikes/protobuf/modImage"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	logLevel       = getEnv("LOG_LEVEL", "info")
	port           = getEnv("PORT", "localhost:4092")
	minioHost      = getEnv("MINIO_HOST", "host.docker.internal:9001")
	minioAccessKey = getEnv("MINIO_ACCESKEY", "fcR98JVhqBHP2laQ")
	minioSecret    = getEnv("MINIO_SECRET", "7xMpT01cPq1O1QpWktOjBDjsgZdccPoL")
)

func main() {
	logger := &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &prefixed.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
	}

	/* Database */
	minioConn, err := minio.New(minioHost, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecret, ""),
		Secure: false,
	})
	if err != nil {
		logger.WithFields(logrus.Fields{"prefix": "MINIO"}).Fatal("unable to open a connection to database")
	}
	minioRepository := repository.NewMinioRepository(minioConn)

	/* Server */
	// Create a tcp listner
	listener, err := net.Listen("tcp", port)
	if err != nil {
		logger.WithFields(logrus.Fields{"prefix": "SERVICE.MODIMAGE"}).Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	protobuffer.RegisterModImageServiceServer(grpcServer, handler.New(minioRepository, *logger))
	reflection.Register(grpcServer)

	// Start grpc server on listener
	logger.WithFields(logrus.Fields{"prefix": "SERVICE.MODIMAGE"}).Infof("is listening on Grpc PORT: {%v}", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		logger.WithFields(logrus.Fields{"prefix": "SERVICE.MODIMAGE"}).Fatalf("failed to serve: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

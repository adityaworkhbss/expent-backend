package prisma

import (
	"log"
	"os"

	"expent-backend/prisma/db"
)

type PrismaClient struct {
	Prisma *db.PrismaClient
}

func NewClient(databaseURL string) (*PrismaClient, error) {
	if err := os.Setenv("DATABASE_URL", databaseURL); err != nil {
		log.Printf("Failed to set DATABASE_URL env: %v", err)
	}

	client := db.NewClient(db.WithDatasourceURL(databaseURL))
	if err := client.Prisma.Connect(); err != nil {
		return nil, err
	}
	return &PrismaClient{Prisma: client}, nil
}

package prisma

import (
	"log"
	"os"

	"expent-backend/prisma/db"
)

type PrismaClient struct {
	Prisma *db.PrismaClient
}

// NewClient initializes a Prisma client with the given DSN.
func NewClient(databaseURL string) (*PrismaClient, error) {
	// The Prisma client expects a DSN like "postgresql://..." which is passed via env var.
	// Ensure the env variable is set for the generated client.
	if err := os.Setenv("DATABASE_URL", databaseURL); err != nil {
		log.Printf("Failed to set DATABASE_URL env: %v", err)
	}
	// Initialize the generated Prisma client.
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return nil, err
	}
	return &PrismaClient{Prisma: client}, nil
}

package prisma

import (
    "log"
    "context"

    "github.com/steebchen/prisma-client-go/runtime/transaction"
    "github.com/steebchen/prisma-client-go/runtime"
    "expent-backend/prisma/client"
)

type PrismaClient struct {
    Prisma *client.PrismaClient
    Tx     *transaction.Manager
}

// NewClient initializes a Prisma client with the given DSN.
func NewClient(databaseURL string) (*PrismaClient, error) {
    // The Prisma client expects a DSN like "postgresql://..." which is passed via env var.
    // Ensure the env variable is set for the generated client.
    if err := runtime.SetEnv("DATABASE_URL", databaseURL); err != nil {
        log.Printf("Failed to set DATABASE_URL env: %v", err)
    }
    // Initialize the generated Prisma client.
    client := client.NewClient()
    if err := client.Prisma.Connect(); err != nil {
        return nil, err
    }
    return &PrismaClient{Prisma: client, Tx: transaction.NewManager()}, nil
}

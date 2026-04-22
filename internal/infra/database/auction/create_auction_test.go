package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	t.Helper()

	mongoURI := os.Getenv("MONGODB_URL")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:admin@localhost:27017/auctions?authSource=admin"
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}

	dbName := "auction_test_" + time.Now().Format("20060102150405")
	db := client.Database(dbName)

	cleanup := func() {
		_ = db.Drop(context.Background())
		_ = client.Disconnect(context.Background())
	}

	return db, cleanup
}

func TestAuctionAutoClose(t *testing.T) {
	// Set a short auction duration for testing
	os.Setenv("AUCTION_DURATION", "2s")
	defer os.Unsetenv("AUCTION_DURATION")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewAuctionRepository(db)

	auctionEntity, internalErr := auction_entity.CreateAuction(
		"Test Product",
		"Electronics",
		"A test product description that is long enough",
		auction_entity.New,
	)
	if internalErr != nil {
		t.Fatalf("failed to create auction entity: %v", internalErr)
	}

	ctx := context.Background()
	if err := repo.CreateAuction(ctx, auctionEntity); err != nil {
		t.Fatalf("failed to insert auction: %v", err)
	}

	// Verify it starts as Active
	var result AuctionEntityMongo
	if err := repo.Collection.FindOne(ctx, bson.M{"_id": auctionEntity.Id}).Decode(&result); err != nil {
		t.Fatalf("failed to find auction: %v", err)
	}
	if result.Status != auction_entity.Active {
		t.Errorf("expected status Active (%d), got %d", auction_entity.Active, result.Status)
	}

	// Wait for auto-close goroutine to fire (duration + buffer)
	time.Sleep(3 * time.Second)

	if err := repo.Collection.FindOne(ctx, bson.M{"_id": auctionEntity.Id}).Decode(&result); err != nil {
		t.Fatalf("failed to find auction after wait: %v", err)
	}
	if result.Status != auction_entity.Completed {
		t.Errorf("expected status Completed (%d) after duration, got %d", auction_entity.Completed, result.Status)
	}
}

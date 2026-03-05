package main

import (
	"context"
	"log"
	"time"

	"github.com/casbin/casbin/v2"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/google/uuid"
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/entity"
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/valueobject"
	"github.com/mikiasgoitom/Secure-Asset/internal/infrastructure/config"
	"github.com/mikiasgoitom/Secure-Asset/internal/infrastructure/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Println("--- Seeder Started ---")

	// --- Load Config and Connect to DB ---
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.DatabaseURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)
	db := client.Database(cfg.DatabaseName)

	// --- Initialize Casbin Enforcer ---
	adapter, err := mongodbadapter.NewAdapterWithCollectionName(clientOptions, cfg.DatabaseName, cfg.CasbinCollection)
	if err != nil {
		log.Fatalf("Failed to create casbin adapter: %v", err)
	}
	enforcer, err := casbin.NewEnforcer(cfg.CasbinModelPath, adapter)
	if err != nil {
		log.Fatalf("Failed t o create casbin enforcer: %v", err)
	}

	// --- Run Seeders ---
	seedCasbinPolicies(enforcer)
	admin, err := seedAdminUser(ctx, db, enforcer)
	viewer, err := seedViewerUser(ctx, db, enforcer)
	seedAssetsAndDACPolicies(ctx, db, enforcer, admin, viewer)

	log.Println("--- Seeder Finished Successfully ---")
}

// seedCasbinPolicies adds the initial permission rules.
func seedCasbinPolicies(enforcer *casbin.Enforcer) {
	log.Println("Seeding Casbin policies...")

	policies := [][]string{
		// RBAC: Employee can create assets
		{"Employee", "/api/v1/asset/create", "POST", "allow"},
		// RBAC: Viewer can ONLY read assets
		{"Viewer", "/api/v1/asset/*", "GET", "allow"},
		// RBAC: Admin can do anything to assets
		{"Admin", "/api/v1/asset/*", "(GET)|(POST)|(PUT)|(DELETE)", "allow"},
	}

	for _, policy := range policies {
		if has, _ := enforcer.HasPolicy(policy); !has {
			enforcer.AddPolicy(policy)
			log.Printf("Policy added: %v", policy)
		} else {
			log.Printf("Policy already exists: %v", policy)
		}
	}

	if err := enforcer.SavePolicy(); err != nil {
		log.Fatalf("Failed to save policies: %v", err)
	}
	log.Println("Casbin policies seeded.")
}

// seedAdminUser creates a default admin if one doesn't exist.
func seedAdminUser(ctx context.Context, db *mongo.Database, enforcer *casbin.Enforcer) (*entity.User, error) {
	log.Println("Seeding admin user...")

	// We need a UserRepository to interact with the users collection
	userRepo := repository.NewUserRepository(db, "users") // Assuming "users" is your collection name

	// Check if the admin user already exists
	adminEmail := "admin@example.com"
	existingAdmin, err := userRepo.FindByEmail(ctx, adminEmail)
	if err != nil {
		log.Fatalf("Error checking for existing admin: %v", err)
	}
	if existingAdmin != nil {
		log.Println("Admin user already exists.")
		return existingAdmin, nil
	}

	// Hash the admin password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("adminpassword"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	admin := &entity.User{
		ID:             uuid.New().String(),
		Username:       "admin",
		Email:          adminEmail,
		Password:       string(hashedPassword),
		Role:           "Admin",
		Department:     "IT",
		ClearanceLevel: valueobject.GetClassificationLevel("Confidential"),
	}

	// Create the user in the database
	createdAdmin, err := userRepo.Create(ctx, admin)
	if err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}
	log.Printf("Admin user created with ID: %s", createdAdmin.ID)

	// Assign the 'Admin' role to the new user in Casbin
	if has, _ := enforcer.HasRoleForUser(createdAdmin.ID, "Admin"); !has {
		_, err = enforcer.AddRoleForUser(createdAdmin.ID, "Admin")
		if err != nil {
			log.Fatalf("Failed to assign admin role in Casbin: %v", err)
		}
		log.Println("Admin role assigned in Casbin.")
	}
	return createdAdmin, nil
}

// seedViewerUser creates a default Viewer if one doesn't exist.
func seedViewerUser(ctx context.Context, db *mongo.Database, enforcer *casbin.Enforcer) (*entity.User, error) {
	log.Println("Seeding Viewer user...")

	// We need a UserRepository to interact with the users collection
	userRepo := repository.NewUserRepository(db, "users") // Assuming "users" is your collection name

	// Check if the Viewer user already exists
	ViewerEmail := "viewer1@example.com"
	existingViewer, err := userRepo.FindByEmail(ctx, ViewerEmail)
	if err != nil {
		log.Fatalf("Error checking for existing viewer: %v", err)
	}
	if existingViewer != nil {
		log.Println("Viewer user already exists.")
		return existingViewer, nil
	}

	// Hash the Viewer password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Viewerpassword"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash Viewer password: %v", err)
	}

	Viewer := &entity.User{
		ID:             uuid.New().String(),
		Username:       "viewer1",
		Email:          ViewerEmail,
		Password:       string(hashedPassword),
		Role:           "Viewer",
		Department:     "IT",
		ClearanceLevel: valueobject.GetClassificationLevel("Internal"),
	}

	// Create the user in the database
	createdViewer, err := userRepo.Create(ctx, Viewer)
	if err != nil {
		log.Fatalf("Failed to create Viewer user: %v", err)
	}
	log.Printf("Viewer user created with ID: %s", createdViewer.ID)

	// Assign the 'Viewer' role to the new user in Casbin
	if has, _ := enforcer.HasRoleForUser(createdViewer.ID, "Viewer"); !has {
		_, err = enforcer.AddRoleForUser(createdViewer.ID, "Viewer")
		if err != nil {
			log.Fatalf("Failed to assign Viewer role in Casbin: %v", err)
		}
		log.Println("Viewer role assigned in Casbin.")
	}
	return createdViewer, nil
}

func seedAssetsAndDACPolicies(ctx context.Context, db *mongo.Database, enforcer *casbin.Enforcer, admin, viewer *entity.User) {
	log.Println("Seeding assets and DAC policies...")
	assetRepo := repository.NewAssetRepository(db, "assets")

	// Asset 1: A server (for ABAC test)
	serverAsset := &entity.Asset{
		ID: "server-001", Name: "Main Web Server", AssetType: "Server", OwnerID: admin.ID,
	}
	assetRepo.Create(ctx, serverAsset)

	// Asset 2: A top-secret document (for RuBAC/MAC test)
	secretDoc := &entity.Asset{
		ID: "doc-ts-001", Name: "Project Chimera Blueprint", AssetType: "Document",
		Classification: valueobject.Confidential, OwnerID: admin.ID,
	}
	assetRepo.Create(ctx, secretDoc)

	// Asset 3: A normal report (for DAC test)
	normalReport := &entity.Asset{
		ID: "report-2025-q4", Name: "Q4 Financial Report", AssetType: "Document",
		Classification: valueobject.Internal, OwnerID: admin.ID,
	}
	assetRepo.Create(ctx, normalReport)

	// DAC Policy: Grant the 'viewer' user direct GET access to the normal report.
	// This is a 'p' rule with a user ID instead of a role.
	dacPolicy := []string{viewer.ID, "/api/v1/asset/report-2025-q4", "GET", "allow"}
	if has, _ := enforcer.HasPolicy(dacPolicy); !has {
		enforcer.AddPolicy(dacPolicy)
		log.Printf("DAC Policy added: %v", dacPolicy)
	}

	enforcer.SavePolicy()
	log.Println("Assets and DAC policies seeded.")
}

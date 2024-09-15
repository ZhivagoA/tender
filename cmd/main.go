package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"tender/internal/handlers"
	"tender/internal/repositories"
	"tender/internal/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DATABASE")

	fmt.Println("Host:", host)
	fmt.Println("Port:", port)
	fmt.Println("User:", user)
	fmt.Println("DB Name:", dbname)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	fmt.Println("Успешное подключение к базе данных:", db)

	userRepo := &repositories.UserRepository{DB: db}
	orgRepo := &repositories.OrganizationRepository{DB: db}
	responsibleRepo := &repositories.OrganizationResponsibleRepository{DB: db}
	bidRepo := &repositories.BidRepository{DB: db}
	tenderRepo := &repositories.TenderRepository{DB: db}

	bidVersionRepo := &repositories.BidVersionRepository{DB: db}
	tenderVersionRepo := &repositories.TenderVersionRepository{DB: db}

	userService := &services.UserService{Repo: userRepo}
	orgService := &services.OrganizationService{Repo: orgRepo}
	responsibleService := &services.OrganizationResponsibleService{
		Repo:             responsibleRepo,
		OrganizationRepo: orgRepo,
		UserRepo:         userRepo,
	}
	bidService := &services.BidService{
		Repo:                        bidRepo,
		TenderRepo:                  tenderRepo,
		VersionRepo:                 bidVersionRepo,
		OrganizationResponsibleRepo: responsibleRepo,
	}

	tenderService := &services.TenderService{
		Repo:                        tenderRepo,
		VersionRepo:                 tenderVersionRepo,
		OrganizationRepo:            orgRepo,
		OrganizationResponsibleRepo: responsibleRepo,
	}

	pingHandler := &handlers.PingHandler{}

	userHandler := &handlers.UserHandler{UserService: userService}
	orgHandler := &handlers.OrganizationHandler{OrganizationService: orgService}
	responsibleHandler := &handlers.OrganizationResponsibleHandler{ResponsibleService: responsibleService}
	bidHandler := &handlers.BidHandler{BidService: bidService}
	tenderHandler := &handlers.TenderHandler{TenderService: tenderService}

	r := mux.NewRouter()

	r.HandleFunc("/ping", pingHandler.Ping).Methods("GET")

	r.HandleFunc("/users/new", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{userID}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/users/{userID}", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{userID}", userHandler.DeleteUser).Methods("DELETE")

	r.HandleFunc("/organizations/new", orgHandler.CreateOrganization).Methods("POST")
	r.HandleFunc("/organizations/{orgID}", orgHandler.GetOrganization).Methods("GET")
	r.HandleFunc("/organizations/{orgID}", orgHandler.UpdateOrganization).Methods("PUT")
	r.HandleFunc("/organizations/{orgID}", orgHandler.DeleteOrganization).Methods("DELETE")

	r.HandleFunc("/organizations/responsibles/new", responsibleHandler.AssignResponsible).Methods("POST")
	r.HandleFunc("/organizations/responsibles/{responsibleID}", responsibleHandler.GetResponsible).Methods("GET")
	r.HandleFunc("/organizations/responsibles/{responsibleID}", responsibleHandler.RemoveResponsible).Methods("DELETE")

	r.HandleFunc("/bids/new", bidHandler.CreateBid).Methods("POST")
	r.HandleFunc("/bids/submit_decision", bidHandler.SubmitDecision).Methods("POST")
	r.HandleFunc("/bids/list", bidHandler.ListBids).Methods("GET")
	r.HandleFunc("/bids/my", bidHandler.MyBids).Methods("GET")
	r.HandleFunc("/bids/status/{bidID}", bidHandler.GetBidStatus).Methods("GET")
	r.HandleFunc("/bids/edit/{bidID}", bidHandler.EditBid).Methods("PUT")
	r.HandleFunc("/bids/rollback", bidHandler.RollbackBidVersion).Methods("POST")
	r.HandleFunc("/bids/reviews/{bidID}", bidHandler.GetBidReviews).Methods("GET")
	r.HandleFunc("/bids/feedback", bidHandler.LeaveBidFeedback).Methods("POST")

	r.HandleFunc("/tenders/new", tenderHandler.CreateTender).Methods("POST")
	r.HandleFunc("/tenders", tenderHandler.ListTenders).Methods("GET")
	r.HandleFunc("/tenders/my", tenderHandler.ListTendersByUser).Methods("GET")
	r.HandleFunc("/tenders/status/{tenderID}", tenderHandler.GetTenderStatus).Methods("GET")
	r.HandleFunc("/tenders/edit", tenderHandler.EditTender).Methods("PUT")
	r.HandleFunc("/tenders/rollback", tenderHandler.RollbackTenderVersion).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}

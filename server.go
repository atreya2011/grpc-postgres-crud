package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/atreya2011/grpc-postgres-crud/postgrescrud"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"google.golang.org/grpc"
)

type service struct {
	DB *gorm.DB
}

var grpcAddrFlag = flag.String("addr", ":7000", "Address host:port")

func main() {
	log.Printf("grpc server start on port %v", *grpcAddrFlag)
	// Step 1. listen for connections on tcp
	lis, err := net.Listen("tcp", *grpcAddrFlag)
	// Always handle errors
	if err != nil {
		log.Fatalf("frak")
	}
	// Step 2. Create a new grpc server instance
	srv := grpc.NewServer()
	// Step 3. Initialize the db and store it in service{}
	serv := initDB()
	// Step 4. Register the Service by passing the new server instance
	// and the server struct created above
	// The register function is present in pb.go
	postgrescrud.RegisterPostgresCrudServer(srv, serv)
	// Step 5. Serve the listener created above
	srv.Serve(lis)
}

// Create The following is the implementation of the Create service
// as defined in the proto file. It can be any implemention.
func (s *service) Create(ctx context.Context, req *postgrescrud.CreatePersonRequest) (*postgrescrud.CreatePersonResponse, error) {
	// create new record by passing req.Person into gorm.Create()
	log.Printf("%T\n", req.Person.MiddleName)
	// Note the use of tx as the database handle once you are within a transaction
	// transactions are used to rollback changes in case of an error
	// this is to prevent creation of new record in case of an error
	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	person := &postgrescrud.Person{
		MiddleName: "",
	}
	foo := tx.Model(&postgrescrud.Person{}).Create(person)
	bar := foo.Error
	log.Println(bar)
	if bar != nil {
		log.Println("Error: ", bar)
		log.Println("Rollback: ", tx.Rollback().Error)
		return &postgrescrud.CreatePersonResponse{Id: 0}, bar
	}

	// return response
	return &postgrescrud.CreatePersonResponse{Id: 1}, tx.Commit().Error
}

// Read The following is the implementation of the Read service
// as defined in the proto file. It can be any implemention.
func (s *service) Read(ctx context.Context, req *postgrescrud.ReadPersonRequest) (*postgrescrud.ReadPersonResponse, error) {
	// find person with id
	// initialize p with a point to person type
	p := new(postgrescrud.Person)
	// store the result in p using gorm.Where
	err := s.DB.Where("id = ?", req.GetId()).First(p).Error
	if err != nil {
		log.Println(err)
	}
	// return response
	return &postgrescrud.ReadPersonResponse{Person: p}, nil
}

// List lists all the fullnames. Request is empty.
// Below is an example implementation
func (s *service) List(ctx context.Context, e *empty.Empty) (*postgrescrud.ListPeopleResponse, error) {
	// initialize slice of pointers to person
	var people []*postgrescrud.Person
	// pass this to gorm.Find() to fill it with results
	err := s.DB.Find(&people).Error
	if err != nil {
		log.Println(err)
	}
	return &postgrescrud.ListPeopleResponse{Peoples: people}, nil
}

// Delete deletes record based on id
// below is an example implementation
// func (*service) Delete() {
// 	// Delete - delete product
// 	db.Delete(&product)
// }

func initDB() *service {
	const (
		host   = "localhost"
		port   = 5432
		user   = "atreya"
		dbname = "atreya"
	)
	// Step 1. Build the parameter string based on const declared above
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	// Step 2. Connect to database with above parameters
	log.Println(psqlInfo)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic("failed to connect database")
	}

	// Step 3. Migrate the schema
	db.DropTableIfExists(&postgrescrud.Person{})
	db.AutoMigrate(&postgrescrud.Person{})

	// Step 4. Create a new service instance and store the db
	s := &service{db}

	// return the db
	return s
}

package main

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite
	store *dbStore
	db    *sql.DB
}

func (s *StoreSuite) SetupSuite() {
	connString := "dbname=bird_encyclopedia sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db
	s.store = &dbStore{db: db}
}

func (s *StoreSuite) SetupTest() {
	_, err := s.db.Query("DELETE FROM birds")
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *StoreSuite) TearDownSuite() {
	s.db.Close()
}

func TestStoreSuite(t *testing.T) {
	s := new(StoreSuite)
	suite.Run(t, s)
}

func (s *StoreSuite) TestCreateBird() {
	s.store.CreateBird(&Bird{
		Description: "test description",
		Species:     "test species",
	})

	res, err := s.db.Query(`SELECT COUNT(*) FROM birds WHERE description='test description' AND SPECIES='test species'`)
	if err != nil {
		s.T().Fatal(err)
	}

	var count int
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}

	if count != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", count)
	}
}

func (s *StoreSuite) TestGetBird() {
	// Insert a sample bird into the `birds` table
	_, err := s.db.Query(`INSERT INTO birds (species, description) VALUES('bird','description')`)
	if err != nil {
		s.T().Fatal(err)
	}

	// Get the list of birds through the stores `GetBirds` method
	birds, err := s.store.GetBirds()
	if err != nil {
		s.T().Fatal(err)
	}

	// Assert that the count of birds received must be 1
	nBirds := len(birds)
	if nBirds != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", nBirds)
	}

	// Assert that the details of the bird is the same as the one we inserted
	expectedBird := Bird{"bird", "description"}
	if *birds[0] != expectedBird {
		s.T().Errorf("incorrect details, expected %v, got %v", expectedBird, *birds[0])
	}
}
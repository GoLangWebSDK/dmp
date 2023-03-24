package database

type DBSeeder interface {
	SeedModel() error
}

type Seeder struct {
	DB      *Database
	seeders []DBSeeder
}

func NewSeeder(db *Database) *Seeder {
	return &Seeder{
		DB: db,
	}
}

func (s *Seeder) Run() {
	for _, seeder := range s.seeders {
		err := seeder.SeedModel()
		if err != nil {
			// tmp error handling...
			panic(err)
		}
	}
}

func (s *Seeder) AddSeeder(seeders ...DBSeeder) *Seeder {
	s.seeders = append(s.seeders, seeders...)
	return s
}

package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"sync"
)

type Repository struct {
	Config config
	DB     *gorm.DB
}
type config struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Dbname   string `mapstructure:"DB_NAME"`
	Schema   string `mapstructure:"DB_SCHEMA"`
}

var once sync.Once
var repository *Repository

func GetInstance() *Repository {
	once.Do(func() {
		repository = &Repository{}
	})
	return repository
}

// ConnectGenerator handle for multiple database type in test just use postgres
func (repo *Repository) ConnectGenerator(typeDataBase string) string {
	switch typeDataBase {
	case "PG":
		return fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%d sslmode=disable`,
			repo.Config.Host, repo.Config.User, repo.Config.Password, repo.Config.Dbname, repo.Config.Port)
	default:
		return ""
	}
}

// ConnectToDB connect to db with name. Example: "PG"
func (repo *Repository) ConnectToDB(typeDatabase string, envPath string) (*gorm.DB, error) {
	err := repo.setConfig(envPath)
	if err != nil {
		return nil, err
	}
	//if err = repo.migrateDB(); err != nil {
	//	return nil, err
	//}
	log.Printf("Applied migrations!")
	db, err := gorm.Open(postgres.Open(repo.ConnectGenerator(typeDatabase)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: repo.Config.Schema + ".",
		},
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// SetConfig CreateConfig load config in env with pg
func (repo *Repository) setConfig(envPath string) error {
	config, err := loadConfig(envPath) // for simple make just one file .env for postgres
	if err != nil {
		return err
	}
	repo.Config = *config
	return nil
}

// load config from .env to app
func loadConfig(path string) (*config, error) {
	config := &config{}
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// init data base in docker
func (repo *Repository) migrateDB() error {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	db, err := sql.Open("postgres", repo.ConnectGenerator("PG"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	_, err = migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

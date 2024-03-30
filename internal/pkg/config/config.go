package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	Secret         string
	AccessExpTime  time.Duration
	RefreshExpTime time.Duration
	DBName         string
	CollectionName string
}

func New() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	read := os.Getenv("READTMT")
	readtimeout, err := strconv.Atoi(read)
	if err != nil {
		log.Fatal(err)
	}

	write := os.Getenv("WRITETMT")
	writetimeout, err := strconv.Atoi(write)
	if err != nil {
		log.Fatal(err)
	}

	maxbytes := os.Getenv("MAXB")
	maxheaderbytes, err := strconv.Atoi(maxbytes)
	if err != nil {
		log.Fatal(err)
	}

	actime := os.Getenv("ACCESSTIME")
	access, err := strconv.Atoi(actime)
	if err != nil {
		log.Fatal(err)
	}

	reftime := os.Getenv("ACCESSTIME")
	refresh, err := strconv.Atoi(reftime)
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		Addr:           os.Getenv("ADDR"),
		Secret:         os.Getenv("SECRET"),
		DBName:         os.Getenv("DBNAME"),
		CollectionName: os.Getenv("COLLNAME"),
		ReadTimeout:    time.Duration(readtimeout),
		WriteTimeout:   time.Duration(writetimeout),
		MaxHeaderBytes: maxheaderbytes,
		AccessExpTime:  time.Duration(access),
		RefreshExpTime: time.Duration(refresh),
	}
}

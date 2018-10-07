package env

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
)

func Load() error {
    if os.Getenv("CHITOI_ENV") == "" {
        os.Setenv("CHITOI_ENV", "development")
    }

    return godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("CHITOI_ENV")))
}

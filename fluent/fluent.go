package fluent

import (
    "os"
    "strconv"

    "github.com/fluent/fluent-logger-golang/fluent"
    "github.com/pkg/errors"
)

func Connect() (*fluent.Fluent, error) {
    port, err := strconv.Atoi(os.Getenv("CHITOI_FLUENT_PORT"))
    if err != nil {
        return nil, errors.Wrap(err, "error convert port string to integer")
    }

    return fluent.New(fluent.Config{
        FluentHost: os.Getenv("CHITOI_FULENT_HOST"),
        FluentPort: port,
    })
}

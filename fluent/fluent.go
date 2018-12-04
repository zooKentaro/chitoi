package fluent

import (
    "os"
    "strconv"

    "github.com/fluent/fluent-logger-golang/fluent"
    "github.com/pkg/errors"
)

type Logger struct {
    logger *fluent.Fluent
}

func Connect() (*Logger, error) {
    port, err := strconv.Atoi(os.Getenv("CHITOI_FLUENT_PORT"))
    if err != nil {
        return nil, errors.Wrap(err, "error convert port string to integer")
    }

    cnf := fluent.Config{
        FluentHost: os.Getenv("CHITOI_FULENT_HOST"),
        FluentPort: port,
    }

    fluent, err := fluent.New(cnf)
    if err != nil {
        return nil, errors.Wrap(err, "error new fluent")
    }

    return &Logger{logger: fluent}, nil
}

func (lgr *Logger) PostError() error {
    return nil
}

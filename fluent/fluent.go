package fluent

import (
    "fmt"
    "os"
    "strconv"
    "time"

    "github.com/fluent/fluent-logger-golang/fluent"
    "github.com/pkg/errors"
)

type Logger struct {
    logger *fluent.Fluent
}

const (
    ErrorTagSuffix = "error"
)

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

func (lgr *Logger) PostError(tag, msg string) error {
    errTag := fmt.Sprintf("%s.%s", tag, ErrorTagSuffix)
    data := map[string]string{
        "time":    time.Now().Format("2006-01-02 15:04:05"),
        "message": msg,
    }

    if err := lgr.logger.Post(errTag, data); err != nil {
        return errors.Wrapf(err, "error post log, tag: %s, data: %+v", errTag, data)
    }
    return nil
}

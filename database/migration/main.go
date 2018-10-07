package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"

    gs "github.com/schemalex/git-schemalex"
    "github.com/uenoryo/chitoi/config"
    "github.com/uenoryo/chitoi/env"
)

const (
    revisionTable = "gslex_revision"
)

var (
    deploy = flag.Bool("deploy", false, "migrationを実際に実行するかどうか (デフォルトは dry-run)")
)

func main() {
    flag.Parse()

    if err := env.Load(); err != nil {
        log.Fatal("error load env, error: ", err.Error())
        return
    }

    ctx := context.Background()
    if err := run(ctx, *deploy); err != nil {
        log.Fatal("error run migration, error: ", err.Error())
    }
}

func run(ctx context.Context, isDeploy bool) error {
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s",
        os.Getenv("CHITOI_DB_USER"),
        os.Getenv("CHITOI_DB_PASS"),
        os.Getenv("CHITOI_DB_HOST"),
        os.Getenv("CHITOI_DB_PORT"),
        os.Getenv("CHITOI_DB_NAME"),
    )

    runner := &gs.Runner{
        Workspace: config.Home(),
        Deploy:    isDeploy,
        Table:     revisionTable,
        DSN:       dsn,
        Schema:    "database/sql/main.sql",
    }

    if isDeploy {
        return runner.Run(ctx)
    }

    return nil
}

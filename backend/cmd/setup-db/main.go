package main

import (
	"context"
	"flag"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"log"
	"os"
	"strings"
	"time"
)

var (
	dsn     *string
	scripts []string
	timeout *time.Duration
)

func init() {
	dsn = flag.String("dsn", "", "YDB connection string")
	timeout = flag.Duration("timeout", time.Minute, "timeout")
}

func main() {
	flag.Parse()
	scripts = flag.Args()

	if dsn == nil || *dsn == "" {
		ydbDSN := os.Getenv("YDB_DSN")
		dsn = &ydbDSN
	}

	if len(scripts) == 0 {
		scriptsEnv := os.Getenv("YDB_SCRIPTS")

		scripts = strings.Split(scriptsEnv, " ")
	}

	if dsn == nil || *dsn == "" || len(scripts) == 0 || scripts[0] == "" {
		flag.Usage()
		return
	}

	log.Printf("dsn: %v\n", *dsn)
	log.Printf("scripts: %v with len: %d\n", scripts, len(scripts))

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(*timeout))

	defer cancel()

	db, err := ydb.Open(ctx, *dsn)
	if err != nil {
		log.Fatalf("failed to connect to YDB: %v", err)
		return
	}

	for i := 0; ; i++ {
		if i >= 3 {
			return
		}

		log.Printf("attempt #%d\n", i)
		res, err := db.Query().Query(ctx, "SELECT * FROM `.sys_health/test`")

		if err == nil {
			_ = res.Close(ctx)
			break
		}

		log.Println("failed to query .sys_health/test")

		time.Sleep(10 * time.Second)
	}

	executeScripts(db, ctx)
}

func executeScripts(db *ydb.Driver, ctx context.Context) {
	for _, script := range scripts {
		buf, err := os.ReadFile(script)
		if err != nil {
			log.Fatalf("failed to read script %q: %v", script, err)
			return
		}

		sql := string(buf)

		result, err := db.Scripting().Execute(ctx, sql, nil)
		if err != nil {
			log.Fatalf("failed to execute script %q: %v", script, err)
			return
		}

		err = result.Close()
		if err != nil {
			log.Fatalf("failed to close result %q: %v", script, err)
		}
	}
}

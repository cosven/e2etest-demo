package workload

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"math/rand"
	"sync"

	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"go.uber.org/zap"
)

type AppendWorkload struct {
	DB          *sql.DB
	Concurrency int
	Tables      int
	PadLength   int
}

// MustExec must execute sql or fatal
func MustExec(DB *sql.DB, query string, args ...interface{}) sql.Result {
	r, err := DB.Exec(query, args...)
	if err != nil {
		log.Fatal("Exec query err.",
			zap.String("query", query),
			zap.Error(err))
	}
	return r
}

func (c *AppendWorkload) Prepare() error {
	// Use 32 threads to create Tables.
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(tid int) {
			defer wg.Done()
			for j := 0; j < c.Tables; j++ {
				if j%32 == tid {
					sql := fmt.Sprintf("drop table if exists write_stress%d", j+1)
					MustExec(c.DB, sql)
					sql = fmt.Sprintf("create table write_stress%d(col1 bigint, col2 varchar(256), data longtext, key k(col1, col2))", j+1)
					MustExec(c.DB, sql)
				}
			}
		}(i)
	}
	wg.Wait()
	log.Info("Prepare ok.")
	return nil
}

func (c *AppendWorkload) Run(ctx context.Context) error {
	log.Info("AppendWorkload start")
	var wg sync.WaitGroup
	for i := 0; i < c.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				err := c.runClient(ctx)
				log.Error("Append row failed", zap.Error(err))
			}
		}()
	}

	wg.Wait()
	log.Info("Run ok!")
	return nil
}

func (c *AppendWorkload) runClient(ctx context.Context) error {
	rng := rand.New(rand.NewSource(rand.Int63()))

	col2 := make([]byte, 192)
	data := make([]byte, c.PadLength)
	for {
		col1 := rng.Int63()
		col2Len := rng.Intn(192)
		_, _ = rng.Read(col2[:col2Len])
		dataLen := rng.Intn(c.PadLength)
		_, _ = rng.Read(data[:dataLen])
		tid := rng.Int()%c.Tables + 1
		sql := fmt.Sprintf("insert into write_stress%d values (?, ?, ?)", tid)
		_, err := c.DB.ExecContext(ctx, sql, col1,
			base64.StdEncoding.EncodeToString(col2[:col2Len]),
			base64.StdEncoding.EncodeToString(data[:dataLen]))
		if err != nil {
			return errors.Trace(err)
		}
	}
}

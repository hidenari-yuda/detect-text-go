package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/hidenari-yuda/paychan-server/domain/config"
	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	"gopkg.in/guregu/null.v4"
)

// DB is a db instance which implement of interfaces.SQL
//
type DB struct {
	db          *sqlx.DB
	printsQuery bool
}

func NewDB(dbConfig config.DB, printsQuery bool) *DB {
	var (
		db            *sqlx.DB
		err           error
		count         = 1
		maxRetryCount = 15
		url           string
	)

	// 本番環境の場合は、CloudSQLにUnix接続する
	if dbConfig.InstanceUnixSocket != "" {
		url = fmt.Sprintf("%s:%s@unix(/%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_general_ci",
			dbConfig.User,
			dbConfig.Pass,
			dbConfig.InstanceUnixSocket,
			dbConfig.Name,
		)

		// ローカル環境の場合は、DockerのMySQLへTCP接続する
	} else {
		url = fmt.Sprintf("%s:%s@tcp([%s]:%d)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true",
			dbConfig.User,
			dbConfig.Pass,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Name,
		)
	}

	for {
		fmt.Println("Trying to connect DB...", url)
		if db, err = sqlx.Connect("mysql", url); err != nil {
			fmt.Println("Something wrong with connecting DB...", err.Error())
		} else {
			fmt.Println("Succeeded!")
			break
		}

		time.Sleep(1000 * time.Millisecond)

		if count < maxRetryCount {
			count++
			continue
		} else {
			panic("Failed to connect DB.")
		}
	}

	return &DB{db: db, printsQuery: printsQuery}
}

func (d *DB) Get(name string, dest interface{}, query string, args ...interface{}) error {
	if d.printsQuery {
		defer measureLatency(name, query, args...)()
	}

	err := d.db.Get(dest, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s:%w", err.Error(), entity.ErrNotFound)
		}
		return fmt.Errorf("%s:%w", err.Error(), entity.ErrDBError)
	}
	return nil
}

func (d *DB) Select(name string, dest interface{}, query string, args ...interface{}) error {
	if d.printsQuery {
		defer measureLatency(name, query, args...)()
	}
	err := d.db.Select(dest, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s:%w", err.Error(), entity.ErrNotFound)
		}
		return fmt.Errorf("%s:%w", err.Error(), entity.ErrDBError)
	}
	return nil
}

func (d *DB) Exec(name string, query string, args ...interface{}) (int64, error) {
	if d.printsQuery {
		defer measureLatency(name, query, args...)()
	}
	r, err := d.db.Exec(query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%s:%w", err.Error(), entity.ErrNotFound)
		}
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 { // Duplicate entry
				return 0, fmt.Errorf("%s:%w", err.Error(), entity.ErrDuplicateEntry)
			}
		}
		return 0, fmt.Errorf("%s:%w", err.Error(), entity.ErrDBError)
	}
	return r.LastInsertId()
}

func (d *DB) Begin() (*Tx, error) {
	tx, err := d.db.Beginx()
	if err != nil {
		return nil, err
	}

	return &Tx{tx: tx, printsQuery: d.printsQuery}, nil
}

func (d *DB) MigrateUp(dir string) error {
	migrations := &migrate.FileMigrationSource{
		Dir: dir,
	}

	_, err := migrate.Exec(d.db.DB, "mysql", migrations, migrate.Up)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (d *DB) MigrateDown(dir string) error {
	migrations := &migrate.FileMigrationSource{
		Dir: dir,
	}

	_, err := migrate.Exec(d.db.DB, "mysql", migrations, migrate.Down)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Tx is a transactional db instance which implement of interfaces.SQL
//
type Tx struct {
	tx          *sqlx.Tx
	printsQuery bool
}

func (t *Tx) Get(name string, dest interface{}, query string, args ...interface{}) error {
	if t.printsQuery {
		defer measureLatency(name, query, args...)()
	}

	err := t.tx.Get(dest, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s:%w", err.Error(), entity.ErrNotFound)
		}
		return fmt.Errorf("%s:%w", err.Error(), entity.ErrDBError)
	}
	return err
}

func (t *Tx) Select(name string, dest interface{}, query string, args ...interface{}) error {
	if t.printsQuery {
		defer measureLatency(name, query, args...)()
	}

	err := t.tx.Select(dest, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s:%w", err.Error(), entity.ErrNotFound)
		}
		return fmt.Errorf("%s:%w", err.Error(), entity.ErrDBError)
	}
	return err
}

func (t *Tx) Exec(name, query string, args ...interface{}) (int64, error) {
	if t.printsQuery {
		defer measureLatency(name, query, args...)()
	}

	r, err := t.tx.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%s:%w", err.Error(), entity.ErrNotFound)
		}
		return 0, fmt.Errorf("%s:%w", err.Error(), entity.ErrDBError)
	}
	return r.LastInsertId()
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func printQuery(prints bool, query string) {
	if !prints {
		return
	}
	log.Println(query)
}

func debugQueryString(query string, args ...interface{}) string {
	var buffer bytes.Buffer
	nArgs := len(args)
	for i, part := range strings.Split(query, "?") {
		buffer.WriteString(part)
		if i < nArgs {
			switch a := args[i].(type) {
			case int64, int, uint:
				buffer.WriteString(fmt.Sprintf("%d", a))
			case float32, float64:
				buffer.WriteString(fmt.Sprintf("%f", a))
			case bool:
				buffer.WriteString(fmt.Sprintf("%t", a))
			case sql.NullBool:
				if a.Valid {
					buffer.WriteString(fmt.Sprintf("%t", a.Bool))
				} else {
					buffer.WriteString("NULL")
				}
			case sql.NullInt64:
				if a.Valid {
					buffer.WriteString(fmt.Sprintf("%d", a.Int64))
				} else {
					buffer.WriteString("NULL")
				}
			case sql.NullString:
				if a.Valid {
					buffer.WriteString(fmt.Sprintf("%q", a.String))
				} else {
					buffer.WriteString("NULL")
				}
			case sql.NullFloat64:
				if a.Valid {
					buffer.WriteString(fmt.Sprintf("%f", a.Float64))
				} else {
					buffer.WriteString("NULL")
				}
			case null.Int:
				if a.Valid {
					buffer.WriteString(fmt.Sprintf("%d", a.Int64))
				} else {
					buffer.WriteString("NULL")
				}
			case null.Bool:
				if a.Valid {
					buffer.WriteString(fmt.Sprintf("%t", a.Bool))
				} else {
					buffer.WriteString("NULL")
				}
			case null.Float:
				if a.Valid {
					buffer.WriteString(fmt.Sprintf("%f", a.Float64))
				} else {
					buffer.WriteString("NULL")
				}
			case null.String:
				if a.Valid {
					buffer.WriteString(fmt.Sprintf("%s", a.String))
				} else {
					buffer.WriteString("NULL")
				}
			case null.Time:
				if a.Valid {
					buffer.WriteString(fmt.Sprintf("%s", a.Time.String()))
				} else {
					buffer.WriteString("NULL")
				}
			default:
				buffer.WriteString(fmt.Sprintf("%q", a))
			}
		}
	}
	return strings.ReplaceAll(strings.ReplaceAll(buffer.String(), "\n", " "), "\t", "")
}

func measureLatency(name, query string, args ...interface{}) func() {
	s := time.Now()
	return func() {
		l := time.Since(s)
		m := struct {
			Type    string `json:"type"`
			Name    string `json:"name"`
			Latency int64  `json:"latency"`
			Query   string `json:"query"`
		}{
			Type:    "sql",
			Name:    name,
			Latency: l.Milliseconds(),
			Query:   debugQueryString(query, args...),
		}

		b, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
	}
}

// func connectWithConnector() (*sql.DB, error) {
// 	mustGetenv := func(k string) string {
// 		v := os.Getenv(k)
// 		if v == "" {
// 			log.Fatalf("Warning: %s environment variable not set.", k)
// 		}
// 		return v
// 	}
// 	// Note: Saving credentials in environment variables is convenient, but not
// 	// secure - consider a more secure solution such as
// 	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
// 	// keep secrets safe.
// 	var (
// 		dbUser                 = mustGetenv("DB_USER")                  // e.g. 'my-db-user'
// 		dbPwd                  = mustGetenv("DB_PASS")                  // e.g. 'my-db-password'
// 		dbName                 = mustGetenv("DB_NAME")                  // e.g. 'my-database'
// 		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
// 		usePrivate             = os.Getenv("PRIVATE_IP")
// 	)

// 	d, err := cloudsqlconn.NewDialer(context.Background())
// 	if err != nil {
// 		return nil, fmt.Errorf("cloudsqlconn.NewDialer: %v", err)
// 	}
// 	var opts []cloudsqlconn.DialOption
// 	if usePrivate != "" {
// 		opts = append(opts, cloudsqlconn.WithPrivateIP())
// 	}
// 	mysql.RegisterDialContext("cloudsqlconn",
// 		func(ctx context.Context, addr string) (net.Conn, error) {
// 			return d.Dial(ctx, instanceConnectionName, opts...)
// 		})

// 	dbURI := fmt.Sprintf("%s:%s@cloudsqlconn(localhost:3306)/%s?parseTime=true",
// 		dbUser, dbPwd, dbName)

// 	dbPool, err := sql.Open("mysql", dbURI)
// 	if err != nil {
// 		return nil, fmt.Errorf("sql.Open: %v", err)
// 	}
// 	return dbPool, nil
// }

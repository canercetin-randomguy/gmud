package mysql

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

const (
    username = "user"
    password = "password"
    hostname = "127.0.0.1:3306"
    dbname   = "characters"
)

func dsn(dbName string) string {
    return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func dbConnection() (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn(""))
    if err != nil {
        log.Printf("Error %s when opening DB\n", err)
        return nil, err
    }
    //defer db.Close()

    ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelfunc()
    res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
    if err != nil {
        log.Printf("Error %s when creating DB\n", err)
        return nil, err
    }
    no, err := res.RowsAffected()
    if err != nil {
        log.Printf("Error %s when fetching rows", err)
        return nil, err
    }
    log.Printf("rows affected %d\n", no)

    db.Close()
    db, err = sql.Open("mysql", dsn(dbname))
    if err != nil {
        log.Printf("Error %s when opening DB", err)
        return nil, err
    }
    //defer db.Close()

    db.SetMaxOpenConns(20)
    db.SetMaxIdleConns(20)
    db.SetConnMaxLifetime(time.Minute * 5)

    ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelfunc()
    err = db.PingContext(ctx)
    if err != nil {
        log.Printf("Errors %s pinging DB", err)
        return nil, err
    }
    log.Printf("Connected to DB %s successfully\n", dbname)
    return db, nil
}

func createCharacterTable(db *sql.DB) error {
    query := `CREATE TABLE IF NOT EXISTS Character (
		`name` varchar(20) NOT NULL DEFAULT 'x',
		`id` varchar(20) NOT NULL DEFAULT NULL,
		`class` varchar(20) NOT NULL DEFAULT NULL,
		`race` varchar(20) NOT NULL DEFAULT NULL,
		`level` int(3) NOT NULL DEFAULT '1'`
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;

    ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelfunc()
    res, err := db.ExecContext(ctx, query)
    if err != nil {
        log.Printf("Error %s when creating Character table", err)
        return err
    }
    rows, err := res.RowsAffected()
    if err != nil {
        log.Printf("Error %s when getting rows affected", err)
        return err
    }
    log.Printf("Rows affected when creating table: %d", rows)
    return nil
}

func main() {
    db, err := dbConnection()
    if err != nil {
        log.Printf("Error %s when getting db connection", err)
        return
    }
    defer db.Close()
    log.Printf("Successfully connected to database")
    err = createCharacterTable(db)
    if err != nil {
        log.Printf("Create Character table failed with error %s", err)
        return
    }
}

ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
defer cancelfunc()
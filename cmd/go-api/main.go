package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Accel-Hack/go-api/internal/app/infra/repository"
	"github.com/Accel-Hack/go-api/internal/app/server"
	"github.com/Accel-Hack/go-api/internal/app/usercase/sample"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"xorm.io/xorm"
)

type GoAPICmd struct {
	flags  *flag.FlagSet
	MySQL  MySQLOption
	Server ServerOption
	Log    LogOption
}

type MySQLOption struct {
	User     string
	Password string
	Addr     string
	Database string
	Table    string
	DSN      string
}

type ServerOption struct {
	Host string
	Port string
}

type LogOption struct {
	Level SlogLevel
}

type SlogLevel struct {
	slog.Level
}

func (l *SlogLevel) Set(s string) error {
	return l.UnmarshalText([]byte(s))
}

func (c *GoAPICmd) Usage() {
	fmt.Fprintf(c.flags.Output(), "Usage of go-api:\nA go-api requires '-mysql.addr' or '-mysql.dsn' (which is prioritized over '-mysql.addr').\n\n")
	c.flags.PrintDefaults()
}

var cmd = &GoAPICmd{
	flags:  flag.NewFlagSet("go-api", flag.ExitOnError),
	MySQL:  MySQLOption{},
	Server: ServerOption{},
	Log:    LogOption{Level: SlogLevel{slog.LevelInfo}},
}

func init() {
	cmd.flags.Var(&cmd.Log.Level, "log.level", "Logging level one of [DEBUG INFO WARN ERROR]")
	cmd.flags.StringVar(&cmd.Server.Host, "server.host", "localhost", "Host to serve")
	cmd.flags.StringVar(&cmd.Server.Port, "server.port", "8080", "Port to serve")
	cmd.flags.StringVar(&cmd.MySQL.User, "mysql.user", "root", "Username")
	cmd.flags.StringVar(&cmd.MySQL.Password, "mysql.password", "", "Password")
	cmd.flags.StringVar(&cmd.MySQL.Addr, "mysql.addr", "localhost:3566", "MySQL URL. Required if mysql.dsn is empty")
	cmd.flags.StringVar(&cmd.MySQL.Database, "mysql.database", "YOUR_APPLICATION", "Database name")
	cmd.flags.StringVar(&cmd.MySQL.Table, "mysql.table", "SAMPLE", "Table name")
	cmd.flags.StringVar(&cmd.MySQL.DSN, "mysql.dsn", "", `Data source name format defined as follow: `+
		`"[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]"`)
}

func main() {
	cmd.flags.Usage = cmd.Usage
	cmd.flags.Parse(os.Args[1:])
	ctx := context.Background()
	if err := cmd.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

func (c *GoAPICmd) Run(ctx context.Context) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: c.Log.Level.Level,
	},
	))
	dsn := (&mysql.Config{
		User:   c.MySQL.User,
		Passwd: c.MySQL.Password,
		Net:    "tcp",
		Addr:   c.MySQL.Addr,
		DBName: c.MySQL.Database,
	}).FormatDSN()
	if c.MySQL.DSN != "" {
		dsn = c.MySQL.DSN
	}
	log.Printf("connect MySql to %q", dsn)
	xormEngine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		return err
	}
	repo := repository.NewSampleXorm(xormEngine, c.MySQL.Table)
	usecase := sample.Usecase{Repository: repo}
	handler := server.InternalSampleHandler{Usecase: usecase, Logger: logger}
	mux := mux.NewRouter()
	handler.Route(mux)
	addr := net.JoinHostPort(c.Server.Host, c.Server.Port)
	s := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	serverChan := make(chan error)
	go func() {
		defer close(serverChan)
		log.Println("Listen on " + addr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverChan <- err
		}
	}()
	sigCtx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()
	select {
	case err := <-serverChan:
		log.Fatal(err)
	case <-sigCtx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		s.Shutdown(shutdownCtx)
		log.Println("Server is shutting down")
	}
	return nil
}

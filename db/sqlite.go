package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/danielkvist/botio/proto"

	"github.com/pkg/errors"

	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

// SQLite wraps a sql.DB client for SQLite
// to satisfy the DB interface.
type SQLite struct {
	Path            string
	Table           string
	client          *sql.DB
	MaxConns        int
	MaxConnLifetime time.Duration
}

// Connect tries to connect to a SQLite database. If it fails it returns a non-nil error. It also tries to create a table for the commands if not exists.
func (sq *SQLite) Connect() error {
	client, err := sql.Open("sqlite3", sq.Path)
	if err != nil {
		return errors.Wrapf(err, "while opening SQLite3 DB on %q", sq.Path)
	}

	client.SetMaxOpenConns(sq.MaxConns)
	client.SetConnMaxLifetime(sq.MaxConnLifetime)

	sq.client = client
	if err := sq.client.Ping(); err != nil {
		return errors.Wrapf(err, "while opening a connection the SQLite DB")
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (command TEXT NOT NULL PRIMARY KEY, response TEXT NOT NULL);", sq.Table)
	stmt, err := sq.client.Prepare(query)
	if err != nil {
		return errors.Wrapf(err, "while preparing SQL query")
	}
	defer stmt.Close()

	if _, err := stmt.Exec(); err != nil {
		log.Println(sq)
		return errors.Wrapf(err, "while creating table %q", sq.Table)
	}

	return nil
}

// Add receives a *proto.BotCommand and adds it to the table designated. If
// something goes wrong while executing the SQL statement it returns a non-nil
// error.
func (sq *SQLite) Add(cmd *proto.BotCommand) error {
	query := fmt.Sprintf("INSERT INTO %s (command, response) VALUES (?, ?)", sq.Table)
	stmt, err := sq.client.Prepare(query)
	if err != nil {
		return errors.Wrapf(err, "while preparing SQL query")
	}
	defer stmt.Close()

	el := cmd.GetCmd().GetCommand()
	val := cmd.GetResp().GetResponse()
	if _, err := stmt.Exec(el, val); err != nil {
		return errors.Wrapf(err, "while adding command %q to table %q", el, sq.Table)
	}

	return nil
}

// Get reveives a *proto.Command and returns the respective *proto.BotCommand
// if exists in the designated table. If not exists or there is any problem
// while executing the SQL statement it returns a non-nil error.
func (sq *SQLite) Get(cmd *proto.Command) (*proto.BotCommand, error) {
	query := fmt.Sprintf("SELECT response FROM %s WHERE command = ?", sq.Table)
	stmt, err := sq.client.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(err, "while preparing SQL query")
	}
	defer stmt.Close()

	el := cmd.GetCommand()
	row := stmt.QueryRow(el)

	var response string
	if err := row.Scan(&response); err != nil {
		return nil, errors.Wrapf(err, "while scanning DB for command %q", el)
	}

	return &proto.BotCommand{
		Cmd: &proto.Command{
			Command: el,
		},
		Resp: &proto.Response{
			Response: response,
		},
	}, nil
}

// GetAll ranges over all the entries of the designated table for the commands
// and returns a *proto.BotCommand with all the *proto.BotCommands found. If
// something goes wrong while executing the SQL statement or while
// getting some *proto.BotCommand it returns a non-nil error.
func (sq *SQLite) GetAll() (*proto.BotCommands, error) {
	query := fmt.Sprintf("SELECT command, response FROM %s", sq.Table)
	stmt, err := sq.client.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(err, "while preparing SQL query")
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, errors.Wrapf(err, "while extracting commands from table %q", sq.Table)
	}
	defer rows.Close()

	var commands []*proto.BotCommand
	for rows.Next() {
		var command string
		var response string
		if err := rows.Scan(&command, &response); err != nil {
			return nil, errors.Wrapf(err, "while getting command from table %q", sq.Table)
		}

		commands = append(commands, &proto.BotCommand{
			Cmd: &proto.Command{
				Command: command,
			},
			Resp: &proto.Response{
				Response: response,
			},
		})
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "while getting commands from table %q", sq.Table)
	}

	return &proto.BotCommands{
		Commands: commands,
	}, nil
}

// Remove removes the received *proto.BotCommand from the designated table. It returns
// a non-nil error if something goes wrong while executing the SQL statement.
func (sq *SQLite) Remove(cmd *proto.Command) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE command = ?", sq.Table)
	stmt, err := sq.client.Prepare(query)
	if err != nil {
		return errors.Wrapf(err, "while preparing SQL query")
	}
	defer stmt.Close()

	el := cmd.GetCommand()
	if _, err := stmt.Exec(el); err != nil {
		return errors.Wrapf(err, "while removing command %q from table %q", el, sq.Table)
	}

	return nil
}

// Update updates the *proto.Response of an existing *proto.BotCommand with the
// *proto.Response of the received *proto.BotCommand. If something goes wrong while
// executing the SQL statement it returns a non-nil error.
func (sq *SQLite) Update(cmd *proto.BotCommand) error {
	query := fmt.Sprintf("UPDATE %s SET response=? WHERE command=?", sq.Table)
	stmt, err := sq.client.Prepare(query)
	if err != nil {
		return errors.Wrapf(err, "while preparing SQL query")
	}
	defer stmt.Close()

	el := cmd.GetCmd().GetCommand()
	val := cmd.GetResp().GetResponse()

	if _, err := stmt.Exec(val, el); err != nil {
		return errors.Wrapf(err, "while updating command %q on table %q", el, sq.Table)
	}

	return nil
}

// Close tries to close the connection to the SQLite database. If it fails
// it returns a non-nil error.
func (sq *SQLite) Close() error {
	if err := sq.client.Close(); err != nil {
		return errors.Wrapf(err, "while closing connection to SQLite3 DB")
	}

	return nil
}

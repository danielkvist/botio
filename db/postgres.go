package db

import (
	"database/sql"
	"fmt"

	"github.com/danielkvist/botio/proto"

	// postgres driver
	_ "github.com/lib/pq"
)

// Postgres wraps a sql.DB client for PostgreSQL and
// satisfies the DB interface.
type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
	Table    string
	client   *sql.DB
}

// Connect tries to connect to a PostgreSQL database. If it fails it returns
// a non-nil error. It also tries to create a table for the commands if not exist.
func (ps *Postgres) Connect() error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", ps.Host, ps.Port, ps.User, ps.Password, ps.DB)
	client, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("while validating arguments to connect to DB: %v", err)
	}

	ps.client = client
	if err := ps.client.Ping(); err != nil {
		return fmt.Errorf("while opening a connection to DB: %v", err)
	}

	statement := `
		CREATE IF NOT EXISTS $1 (
			command TEXT NOT NULL PRIMARY KEY,
			response TEXT NOT NULL
		);`

	if _, err := ps.client.Exec(statement, ps.Table); err != nil {
		return fmt.Errorf("while creating a table for commands: %v", err)
	}

	return nil
}

// Add receives a *proto.BotCommand and adds it to the
// table designated. If something goes wrong executing the
// SQL statement it returns a non-nil error.
func (ps *Postgres) Add(cmd *proto.BotCommand) error {
	statement := `INSERT INTO $1 (command, response) VALUES ($2, $3);`

	el := cmd.GetCmd().GetCommand()
	val := cmd.GetResp().GetResponse()
	if _, err := ps.client.Exec(statement, ps.Table, el, val); err != nil {
		return fmt.Errorf("while adding command %q: %v", el, err)
	}

	return nil
}

// Get receives a *proto.Command and returns the respective *proto.BotCommand
// if exists in the designated table. If not or there is any problem
// while executing the SQL statement it returns a non-nil error.
func (ps *Postgres) Get(cmd *proto.Command) (*proto.BotCommand, error) {
	el := cmd.GetCommand()

	statement := `SELECT * FROM $1 WHERE command=$2;`
	row := ps.client.QueryRow(statement, ps.Table, el)

	var command *proto.BotCommand
	if err := row.Scan(command.Cmd.Command, command.Resp.Response); err != nil {
		return nil, fmt.Errorf("while getting command %q: %v", el, err)
	}

	return command, nil
}

// GetAll ranges over all the entries of the designated table for the commands
// and returns a *proto.BotCommands with all the *proto.BotCommand found.
// If something goes wrong while executing the SQL statement or while
// getting some command it returns a non-nil error.
func (ps *Postgres) GetAll() (*proto.BotCommands, error) {
	statement := `SELECT * FROM $1;`
	rows, err := ps.client.Query(statement, ps.Table)
	if err != nil {
		return nil, fmt.Errorf("while getting commands from DB: %v", err)
	}
	defer rows.Close()

	var commands []*proto.BotCommand
	for rows.Next() {
		var command *proto.BotCommand
		if err := rows.Scan(command.Cmd.Command, command.Resp.Response); err != nil {
			return nil, fmt.Errorf("while getting command: %v", err)
		}

		commands = append(commands, command)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("while getting commands: %v", err)
	}

	return &proto.BotCommands{
		Commands: commands,
	}, nil
}

// Remove removes a *proto.BotCommand from the designated table. It returns a
// non-nil error if there is some problem while executing the
// SQL statement or deleting the command.
func (ps *Postgres) Remove(cmd *proto.Command) error {
	el := cmd.GetCommand()

	statement := `DELETE FROM $1 WHERE command=$2;`
	if _, err := ps.client.Exec(statement, ps.Table, el); err != nil {
		return fmt.Errorf("while removing command %q: %v", el, err)
	}

	return nil
}

// Update updates the Response of an existing *proto.BotCommand
// with the Response of the received *proto.BotCommand. If
// there is any error while executing the SQL statement
// it returns a non-nil error.
func (ps *Postgres) Update(cmd *proto.BotCommand) error {
	el := cmd.GetCmd().GetCommand()
	val := cmd.GetResp().GetResponse()

	statement := `
	UPDATE $1
	SET response=$2
	WHERE command=$3;`

	if _, err := ps.client.Exec(statement, val, el); err != nil {
		return fmt.Errorf("while updating command %q: %v", el, err)
	}

	return nil
}

// Close tries to close the connection to the PostgreSQL database.
// If fails it returns a non-nil error.
func (ps *Postgres) Close() error {
	if err := ps.client.Close(); err != nil {
		return fmt.Errorf("while closing connection to DB: %v", err)
	}

	return nil
}

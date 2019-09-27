package db

import (
	"database/sql"
	"fmt"

	"github.com/danielkvist/botio/models"
)

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
	Table    string
	client   *sql.DB
}

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

	return nil
}

func (ps *Postgres) Add(el, val string) (*models.Command, error) {
	statement := `
	INSERT INTO $1 (command, response )
	VALUES ($2, $3);`

	if _, err := ps.client.Exec(statement, ps.Table, el, val); err != nil {
		return nil, fmt.Errorf("while adding command %q: %v", el, err)
	}

	return &models.Command{
		Cmd:      el,
		Response: val,
	}, nil
}

func (ps *Postgres) Get(el string) (*models.Command, error) {
	statement := `SELECT * FROM $1 WHERE command=$2;`
	row := ps.client.QueryRow(statement, ps.Table, el)

	var command *models.Command
	if err := row.Scan(command.Cmd, command.Response); err != nil {
		return nil, fmt.Errorf("while getting command %q: %v", el, err)
	}

	return command, nil
}

func (ps *Postgres) GetAll() ([]*models.Command, error) {
	statement := `SELECT * FROM $1;`
	rows, err := ps.client.Query(statement, ps.Table)
	if err != nil {
		return nil, fmt.Errorf("while getting commands from DB: %v", err)
	}
	defer rows.Close()

	var commands []*models.Command
	for rows.Next() {
		var command *models.Command
		if err := rows.Scan(command.Cmd, command.Response); err != nil {
			return nil, fmt.Errorf("while getting command for list of commands: %v", err)
		}

		commands = append(commands, command)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("while getting commands: %v", err)
	}

	return commands, nil
}

func (ps *Postgres) Remove(el string) error {
	statement := `DELETE FROM $1 WHERE command=$2;`
	if _, err := ps.client.Exec(statement, ps.Table, el); err != nil {
		return fmt.Errorf("while removing command %q: %v", el, err)
	}

	return nil
}

func (ps *Postgres) Update(el, val string) (*models.Command, error) {
	statement := `
	UPDATE $1
	SET command=$2
		response=$3
	WHERE command=$4;`

	if _, err := ps.client.Exec(statement, el, val); err != nil {
		return nil, fmt.Errorf("while updating command %q: %v", el, err)
	}

	return &models.Command{
		Cmd:      el,
		Response: val,
	}, nil
}

func (ps *Postgres) Close() error {
	if err := ps.client.Close(); err != nil {
		return fmt.Errorf("while closing connection to DB: %v", err)
	}

	return nil
}

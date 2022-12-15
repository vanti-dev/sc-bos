package db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Enrollment struct {
	Name        string
	Description string
	Address     string
	Cert        []byte
}

func GetEnrollment(ctx context.Context, tx pgx.Tx, name string) (en Enrollment, err error) {
	// language=postgresql
	query := `
		SELECT description, address, cert
		FROM enrollment
		WHERE name = $1;
    `

	row := tx.QueryRow(ctx, query, name)
	var descNull *string
	err = row.Scan(&descNull, &en.Address, &en.Cert)
	if descNull != nil {
		en.Description = *descNull
	}
	en.Name = name
	return
}

func CreateEnrollment(ctx context.Context, tx pgx.Tx, en Enrollment) error {
	// language=postgresql
	query := `
		INSERT INTO enrollment (name, description, address, cert) 
		VALUES ($1, $2, $3, $4);
	`

	var descNull *string
	if en.Description != "" {
		descNull = &en.Description
	}

	_, err := tx.Exec(ctx, query, en.Name, descNull, en.Address, en.Cert)
	return err
}

func ListEnrollments(ctx context.Context, tx pgx.Tx) ([]Enrollment, error) {
	// language=postgresql
	query := `
		SELECT name, description, address, cert
		FROM enrollment
		ORDER BY name;
	`

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []Enrollment
	for rows.Next() {
		var en Enrollment
		err = rows.Scan(&en.Name, &en.Description, &en.Address, &en.Cert)
		if err != nil {
			return nil, err
		}

		enrollments = append(enrollments, en)
	}
	return enrollments, nil
}

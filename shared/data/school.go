package data

import (
	"embed"
	"encoding/csv"
	"fmt"
	"strings"
)

//go:embed school_*.csv
var school embed.FS

type School struct {
	accreditations map[string]string
	levels         map[string]string
	ownerships     map[string]string
}

func LoadSchool() (*School, error) {
	// Open the CSV files
	saf, err := school.Open("school_accreditations.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to open file (school_accreditations.csv): %w", err)
	}
	defer saf.Close()

	slf, err := school.Open("school_levels.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to open file (school_levels.csv): %w", err)
	}
	defer slf.Close()

	sof, err := school.Open("school_ownerships.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to open file (school_ownerships.csv): %w", err)
	}
	defer sof.Close()

	// Read the CSV files
	sar, err := csv.NewReader(saf).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read file (school_accreditations.csv): %w", err)
	}

	slr, err := csv.NewReader(slf).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read file (school_levels.csv): %w", err)
	}

	sor, err := csv.NewReader(sof).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read file (school_ownerships.csv): %w", err)
	}

	// Extract data from the CSV files
	accreditations := make(map[string]string, len(sar))
	for i, record := range sar {
		if i == 0 {
			continue
		}
		if len(record) != 2 {
			return nil, fmt.Errorf("invalid row %d: %+v", i, record)
		}
		accreditations[strings.ToLower(record[1])] = record[1]
	}

	levels := make(map[string]string, len(slr))
	for i, record := range slr {
		if i == 0 {
			continue
		}
		if len(record) != 2 {
			return nil, fmt.Errorf("invalid row %d: %+v", i, record)
		}
		levels[strings.ToLower(record[1])] = record[1]
	}

	ownerships := make(map[string]string, len(sor))
	for i, record := range sor {
		if i == 0 {
			continue
		}
		if len(record) != 2 {
			return nil, fmt.Errorf("invalid row %d: %+v", i, record)
		}
		ownerships[strings.ToLower(record[1])] = record[1]
	}

	return &School{
		accreditations: accreditations,
		levels:         levels,
		ownerships:     ownerships,
	}, nil
}

func (s *School) Accreditations() map[string]string {
	accreditations := make(map[string]string, len(s.accreditations))
	for k, v := range s.accreditations {
		accreditations[k] = v
	}
	return accreditations
}

func (s *School) Levels() map[string]string {
	levels := make(map[string]string, len(s.levels))
	for k, v := range s.levels {
		levels[k] = v
	}
	return levels
}

func (s *School) Ownerships() map[string]string {
	ownerships := make(map[string]string, len(s.ownerships))
	for k, v := range s.ownerships {
		ownerships[k] = v
	}
	return ownerships
}

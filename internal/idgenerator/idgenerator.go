package idgenerator

import (
	"errors"
	"time"

	"github.com/sony/sonyflake"
)

type IdGenConfig struct {
	StartTime      time.Time
	MachineID      func() (uint16, error)
	CheckMachineID func(uint16) bool
}

type IdGenerator struct {
	snowflake *sonyflake.Sonyflake
}

func NewIdGenerator(conf IdGenConfig) (*IdGenerator, error) {
	// Create a new Sonyflake instance
	sf := sonyflake.NewSonyflake(sonyflake.Settings{
		// Start time is important - choose a recent timestamp
		StartTime: time.Now().UTC(),
		// Optional: Provide a machine ID
		MachineID: conf.MachineID,
	})
	if sf == nil {
		return nil, errors.New("sonyflake couldn't be created")
	}

	return &IdGenerator{sf}, nil
}

func (idg *IdGenerator) GenerateId() (uint64, error) {
	id, err := idg.snowflake.NextID()
    if err != nil {
        return 0, err
    }

	return id, nil
}
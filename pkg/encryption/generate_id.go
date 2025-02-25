package encryption

import (
	"os"
	"strconv"
	"time"

	"github.com/godruoyi/go-snowflake"
	"github.com/yorukot/go-template/pkg/logger"
)

// Init snowflake to set MachineID and start time
func init() {
	machineID := os.Getenv("MACHINE_ID")
	num, err := strconv.Atoi(machineID)

	if err != nil {
		logger.Log.Fatal("Error to init snowflake invalid MACHINE_ID")
	}

	snowflake.SetMachineID(uint16(num))
	snowflake.SetStartTime(time.Date(2024, 10, 24, 0, 0, 0, 0, time.UTC))
}

// Generate new snowflake ID
func GenerateID() uint64 {
	return snowflake.ID()
}

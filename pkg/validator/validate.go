package validator

import (
	"fmt"
	"github.com/microlib/simple"
	"os"
	"strconv"
	"strings"
)

// checkEnvars - private function, iterates through each item and checks the required field
func checkEnvar(item string, logger *simple.Logger) error {
	name := strings.Split(item, ",")[0]
	required, _ := strconv.ParseBool(strings.Split(item, ",")[1])
	logger.Trace(fmt.Sprintf("Input paramaters -> name %s : required %t", name, required))
	if os.Getenv(name) == "" {
		if required {
			logger.Error(fmt.Sprintf("%s envar is mandatory please set it", name))
			return fmt.Errorf(fmt.Sprintf("%s envar is mandatory please set it", name))
		}

		logger.Error(fmt.Sprintf("%s envar is empty please set it", name))
	}
	return nil
}

// ValidateEnvars : public call that groups all envar validations
// These envars are set via the openshift template
func ValidateEnvars(logger *simple.Logger) error {
	items := []string{
		"LOG_LEVEL,false",
		"SERVER_PORT,true",
		"COUCHBASE_HOST,true",
		"COUCHBASE_DATABASE,true",
		"COUCHBASE_USER,true",
		"COUCHBASE_PASSWORD,true",
		"VERSION,true",
		"NAME,true",
		"COUCHBASE_BUCKET,true",
	}
	for x := range items {
		if err := checkEnvar(items[x], logger); err != nil {
			return err
		}
	}
	return nil
}

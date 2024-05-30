package environment

import "os"

func GetEnvVar(name string) string {
	return os.Getenv(name)
}

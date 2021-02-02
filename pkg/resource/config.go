package resource

import (
	"fmt"
	"os"
	"strconv"
)

var (
	Port       string
	AppVersion string
	AppName    string
	Framework  string

	OtlpEndpoint string

	Mode             string
	debug            string
	profile          string
	Base             string
	MaxNumber        int
	MaxBond          int
	tracePropagation string
	logMethod        string
	Number           int32
	Module           string
	Component        string
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
func Configure() {
	// General Settings
	Port = getEnv("PORT", "3330")
	AppName = getEnv("SERVICE_NAME", "treactor-app")
	AppVersion = getEnv("SERVICE_VERSION", "0.0")
	Framework = "golang"
	// Reactor Specific Settings
	Mode = getEnv("TREACTOR_MODE", "local")
	Module = getEnv("TREACTOR_MODULE", "treactor")
	Component = getEnv("TREACTOR_COMPONENT", "app")
	// Reactor Fixed Settings
	Base = "/treact"

	MaxNumber, _ = strconv.Atoi(getEnv("TREACTOR_MAX_NUMBER", "103"))
	MaxBond, _ = strconv.Atoi(getEnv("TREACTOR_MAX_BOND", "5"))
	n, _ := strconv.Atoi(getEnv("TREACTOR_NUMBER", "0"))
	Number = int32(n)

	OtlpEndpoint = getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "")

	tracePropagation = getEnv("TREACTOR_TRACE_PROPAGATION", "w3c")
	logMethod = os.Getenv("TREACTOR_LOG_METHOD")
}

func IsLocalMode() bool {
	return "local" == Mode
}

func IsKubernetesMode() bool {
	return "cluster" == Mode
}

func MoleculeUrl(molecule string) string {
	if Mode == "cluster" {
		if Module == "bond" {
			if Component == "n" {
				return fmt.Sprintf("http://bond-n/treact/bonds/n?molecule=%s&execute=1", molecule)
			}
			next,_ := strconv.Atoi(Component)
			next++
			if next > MaxBond {
				return fmt.Sprintf("http://bond-n/treact/bonds/n?molecule=%s&execute=1", molecule)
			}
			return fmt.Sprintf("http://bond-%d/treact/bonds/%d?molecule=%s&execute=1", next, next, molecule)
		}
		return fmt.Sprintf("http://bond-1/treact/bonds/1?molecule=%s&execute=1", molecule)
	} else {
		return fmt.Sprintf("http://localhost:%s/treact/bonds/n?molecule=%s&execute=1", Port, molecule)
	}
}

func TracePropagation() {
	//switch traceInternal {
	//case "b3":
	//	return &b3.HTTPFormat{}
	//default:
	//	return &b3.HTTPFormat{}
	//}
	return
}

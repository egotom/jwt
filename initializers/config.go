package initializers
import (
	"os"
    "log"
	"jwt/models"
    "gopkg.in/yaml.v3"
)

var Config models.Config
func LoadConfig(fn string){
	body, err := os.ReadFile(fn)
	if err != nil {
		log.Fatalf("initializers.LoadConfig() error: %v", err)
	}
	err = yaml.Unmarshal(body, &Config)
	if err != nil {
		log.Fatalf("initializers.LoadConfig() error: %v", err)
	}
	// log.Println("--- Config:\n%v\n\n", Config)
}
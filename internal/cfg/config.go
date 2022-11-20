package cfg

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	UrtsiPort        string            `yaml:"urtsiPort"`
	UrtsiDelay       uint8             `yaml:"urtsiCommandDelay"`
	GrafikEyePort    string            `yaml:"grafikEyePort"`
	LumagenPort      string            `yaml:"lumagenPort"`
	LumagenActions   []LumagenAction   `yaml:"lumagenActions"`
	GrafikEyeActions []GrafikEyeAction `yaml:"grafikEyeActions"`
}

type LumagenAction struct {
	Aspects                []uint8 `yaml:"aspects"`
	Framerates             []uint8 `yaml:"framerates"`
	GrafikEyePressButtonId uint8   `yaml:"grafikEyePressButtonId"`
	UrtsiCommand           string  `yaml:"urtsiCommand"`
}

type GrafikEyeAction struct {
	ButtonId     uint8  `yaml:"buttonId"` // see constants in GrafikEye module
	Event        uint8  `yaml:"eventId"`  // press (3) release (4)
	UrtsiCommand string `yaml:"urtsiCommand"`
}

func (c *Config) GrafikEyeConfigured() bool {
	if len(c.GrafikEyePort) > 0 {
		return true
	}
	return false
}

func (c *Config) LumagenConfigured() bool {
	if len(c.LumagenPort) > 0 {
		return true
	}
	return false
}

func (c *Config) UrtsiConfigured() bool {
	if len(c.UrtsiPort) > 0 {
		return true
	}
	return false
}

func (a *LumagenAction) FramerateIn(framerate uint8) bool {
	for _, f := range a.Framerates {
		if f == framerate {
			return true
		}
	}
	return false
}

func (a *LumagenAction) AspectIn(aspect uint8) bool {
	for _, f := range a.Aspects {
		if f == aspect {
			return true
		}
	}
	return false
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

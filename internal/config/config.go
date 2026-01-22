package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the complete application configuration
type Config struct {
	GrindPlan        GrindPlan        `yaml:"grind_plan"`
	SMTP             SMTP             `yaml:"smtp"`
	Reminders        Reminders        `yaml:"reminders"`
	CustomProblemSet CustomProblemSet `yaml:"custom_problem_set"`
	Logger           Logger           `yaml:"logger"`
}

// GrindPlan contains study plan settings
type GrindPlan struct {
	Weeks        int `yaml:"weeks"`
	HoursPerWeek int `yaml:"hours_per_week"`
}

// SMTP contains email configuration
type SMTP struct {
	Provider string `yaml:"provider"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

// Reminders contains reminder notification settings
type Reminders struct {
	DailyReminder      bool   `yaml:"daily_reminder"`
	DailyTime          string `yaml:"daily_time"`
	WeeklyOverview     bool   `yaml:"weekly_overview"`
	WeeklyDay          string `yaml:"weekly_day"`
	WeeklyTime         string `yaml:"weekly_time"`
	StruggleReminder   bool   `yaml:"struggle_reminder"`
	StruggleDaysBefore int    `yaml:"struggle_days_before"`
}

// CustomProblemSet contains custom problem set configuration
type CustomProblemSet struct {
	Enabled      bool   `yaml:"enabled"`
	Path         string `yaml:"path,omitempty"`
	Weeks        int    `yaml:"weeks,omitempty"`
	HoursPerWeek int    `yaml:"hours_per_week,omitempty"`
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(home, ".g7c", "config.yaml")
	return configPath, nil
}

// Logger contains logging configuration
type Logger struct {
	Debug bool `yaml:"debug"`
}

func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) SaveConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// ParseDailyTime parses the daily reminder time
func (r *Reminders) ParseDailyTime() (time.Time, error) {
	return time.Parse("15:04", r.DailyTime)
}

// ParseWeeklyTime parses the weekly reminder time
func (r *Reminders) ParseWeeklyTime() (time.Time, error) {
	return time.Parse("15:04", r.WeeklyTime)
}

// GetSMTPPassword gets password from env var if not set in config
func (s *SMTP) GetSMTPPassword() string {
	if s.Password != "" {
		return s.Password
	}
	return os.Getenv("G7C_SMTP_PASSWORD")
}

// ExpandPath expands ~ in the custom problem set path
func (c *CustomProblemSet) ExpandPath() (string, error) {
	if c.Path == "" {
		return "", nil
	}

	if c.Path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, c.Path[2:]), nil
	}

	return c.Path, nil
}

// GetDefault returns a Config with default values
func GetDefault() *Config {
	return &Config{
		GrindPlan: GrindPlan{
			Weeks:        10,
			HoursPerWeek: 5,
		},
		SMTP: SMTP{
			Provider: "",
			Email:    "",
			Password: "",
		},
		Reminders: Reminders{
			DailyReminder:      false,
			DailyTime:          "09:00",
			WeeklyOverview:     false,
			WeeklyDay:          "monday",
			WeeklyTime:         "09:00",
			StruggleReminder:   false,
			StruggleDaysBefore: 1,
		},
		CustomProblemSet: CustomProblemSet{
			Enabled: false,
		},
		Logger: Logger{
			Debug: false,
		},
	}
}

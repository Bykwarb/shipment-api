package model

type Config struct {
	Server struct {
		Host string
		Port string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	Filter struct {
		ExpectedNumElements      int     `yaml:"expected_num_elements"`
		FalsePositiveProbability float64 `yaml:"false_positive_probability"`
	}
}

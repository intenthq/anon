package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Run("if the file doesn't exist", func(t *testing.T) {
		conf, err := loadConfig("non-existing-file")
		assert.Nil(t, conf, "should return nil if the file doesn't exist")
		assert.Error(t, err, "should return the error if the file doesn't exist")
	})
	t.Run("if the json can't be decoded", func(t *testing.T) {
		conf, err := loadConfig("config_invalid_test.json")
		assert.Nil(t, conf, "should return nil if the json can't be decoded")
		assert.Error(t, err, "should return the error if the json can't be decoded")
	})
	t.Run("default config values", func(t *testing.T) {
		conf, err := loadConfig("config_defaults_test.json")
		require.Nil(t, err, "should return no error if the config can be loaded")
		assert.Equal(t, Config{
			Csv: CsvConfig{
				Delimiter: ",",
			},
			Sampling: SamplingConfig{
				Mod:      1,
				IDColumn: 0,
			},
			Actions: []ActionConfig{},
		}, *conf, "should fill the config with the default values")
	})
	t.Run("if the config can be loaded", func(t *testing.T) {
		gte := 0.0
		lt := 100.0
		output := "0-100"
		conf, err := loadConfig("config_test.json")
		require.Nil(t, err, "should return no error if the config can be loaded")
		assert.Equal(t, Config{
			Csv: CsvConfig{
				Delimiter: "|",
			},
			Sampling: SamplingConfig{
				Mod:      77,
				IDColumn: 84,
			},
			Actions: []ActionConfig{
				ActionConfig{
					Name: "hash",
				},
				ActionConfig{
					Name: "outcode",
				},
				ActionConfig{
					Name: "year",
					DateConfig: DateConfig{
						Format: "20060102",
					},
				},
				ActionConfig{
					Name: "ranges",
					RangeConfig: []RangeConfig{
						RangeConfig{
							Gte:    &gte,
							Lt:     &lt,
							Output: &output,
						},
					},
				},
				ActionConfig{
					Name: "nothing",
				},
			},
		}, *conf, "should return the config properly decoded")
	})
}

func TestValidateConfig(t *testing.T) {
	t.Run("range config validation", func(t *testing.T) {
		num := 2.0
		output := "0-100"
		t.Run("range has at least one of lt, lte, gt, gte", func(t *testing.T) {
			conf := Config{
				Actions: []ActionConfig{
					ActionConfig{
						Name: "range",
						RangeConfig: []RangeConfig{
							RangeConfig{
								Output: &output,
							},
						},
					},
				},
			}
			assert.Error(t, conf.validate(), "should return an error")
		})
		t.Run("range contains both lt and lte", func(t *testing.T) {
			conf := Config{
				Actions: []ActionConfig{
					ActionConfig{
						Name: "range",
						RangeConfig: []RangeConfig{
							RangeConfig{
								Lt:     &num,
								Lte:    &num,
								Output: &output,
							},
						},
					},
				},
			}
			assert.Error(t, conf.validate(), "should return an error")
		})
		t.Run("range contains both gt and gte", func(t *testing.T) {
			conf := Config{
				Actions: []ActionConfig{
					ActionConfig{
						Name: "range",
						RangeConfig: []RangeConfig{
							RangeConfig{
								Gt:     &num,
								Gte:    &num,
								Output: &output,
							},
						},
					},
				},
			}
			assert.Error(t, conf.validate(), "should return an error")
		})
		t.Run("range without output defined", func(t *testing.T) {
			conf := Config{
				Actions: []ActionConfig{
					ActionConfig{
						Name: "range",
						RangeConfig: []RangeConfig{
							RangeConfig{
								Lt:  &num,
								Gte: &num,
							},
						},
					},
				},
			}
			assert.Error(t, conf.validate(), "should return an error")
		})
		t.Run("config contains a correct range", func(t *testing.T) {
			conf := Config{
				Actions: []ActionConfig{
					ActionConfig{
						Name: "range",
						RangeConfig: []RangeConfig{
							RangeConfig{
								Lt:     &num,
								Gte:    &num,
								Output: &output,
							},
						},
					},
				},
			}
			assert.NoError(t, conf.validate(), "should not return an error")
		})
	})
}

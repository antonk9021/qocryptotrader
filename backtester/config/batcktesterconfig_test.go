package config

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/antonk9021/qocryptotrader/backtester/common"
	"github.com/antonk9021/qocryptotrader/common/file"
	"github.com/antonk9021/qocryptotrader/encoding/json"
)

func TestLoadBacktesterConfig(t *testing.T) {
	t.Parallel()
	cfg, err := GenerateDefaultConfig()
	if err != nil {
		t.Error(err)
	}
	testConfig, err := json.Marshal(cfg)
	if err != nil {
		t.Error(err)
	}
	dir := t.TempDir()
	f := filepath.Join(dir, "test.config")
	err = file.Write(f, testConfig)
	if err != nil {
		t.Error(err)
	}
	_, err = ReadBacktesterConfigFromPath(f)
	if err != nil {
		t.Error(err)
	}

	_, err = ReadBacktesterConfigFromPath("test")
	if !errors.Is(err, common.ErrFileNotFound) {
		t.Errorf("received '%v' expected '%v'", err, common.ErrFileNotFound)
	}
}

func TestGenerateDefaultConfig(t *testing.T) {
	t.Parallel()
	cfg, err := GenerateDefaultConfig()
	if err != nil {
		t.Error(err)
	}
	if !cfg.PrintLogo {
		t.Errorf("received '%v' expected '%v'", cfg.PrintLogo, true)
	}
}

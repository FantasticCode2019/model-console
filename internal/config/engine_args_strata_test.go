package config

import (
	"reflect"
	"testing"
)

func TestLoadStrata(t *testing.T) {
	cfg, err := Load(mapEnv(minimalEnv(map[string]string{
		"ENGINE_KIND":  "strata",
		"MODEL_SOURCE": "hf://ISTA-DASLab/Qwen3.8-Flash-Next-GSQ-RCO-GGUF --include IQ2_XS/a.gguf --include IQ2_XS/b.gguf --subdir IQ2_XS",
		"ENGINE_ARGS":  "--context=32768 --parallel 2 --low-ram off --kv int8 --future-option 'two words'",
	})))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Engine.Kind != EngineStrata || cfg.Engine.URL != "http://strata:8080" {
		t.Fatalf("unexpected engine: %+v", cfg.Engine)
	}
	if got, ok := DeriveContextSize(EngineStrata, cfg.Engine.Args); !ok || got != 32768 {
		t.Fatalf("context = (%d, %v)", got, ok)
	}
	if got, ok := DeriveMaxConcurrency(EngineStrata, cfg.Engine.Args); !ok || got != 2 {
		t.Fatalf("concurrency = (%d, %v)", got, ok)
	}
	if got, ok := DerivePoolTokens(EngineStrata, cfg.Engine.Args); ok {
		t.Fatalf("Strata does not report a KV pool, got %d", got)
	}
	if !reflect.DeepEqual(cfg.Engine.Args.Unknown, []string{"--future-option", "two words"}) {
		t.Fatalf("unknown flags lost: %v", cfg.Engine.Args.Unknown)
	}
}

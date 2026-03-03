package logger

import "testing"

func TestLogger(t *testing.T) {
	if err := InitLogger(
		WithLevel(TraceLevel),
		WithFileLog(true),
		WithLogDir("./logs"),
		WithLogFile("app.log"),
		WithMaxSize(1),
		WithJSONFormat(true),
	); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10000; i++ {
		Trace().Str("foo", "bar").Msg("Hello World")
		Debug().Str("foo", "bar").Msg("Hello World")
		Info().Str("foo", "bar").Msg("Hello World")
		Warn().Str("foo", "bar").Msg("Hello World")
		Error().Str("foo", "bar").Msg("Hello World")
	}
}

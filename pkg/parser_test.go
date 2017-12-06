package pkg_test

import (
	"bytes"
	"testing"

	"github.com/influx6/box/pkg"
	"github.com/influx6/faux/tests"
)

func TestMessageParser(t *testing.T) {
	var parser pkg.OpParser
	parser.Secret = "wreckage"

	t.Logf("\tWhen parsing message with incorrect secret")
	{
		raw := []byte("wreckoge#232uFR5   MSG   welcome to the league!")
		_, err := parser.Parse(raw)
		if err == nil {
			tests.Failed("Should have failed to parse message")
		}
		tests.Passed("Should have failed to parse message")

		if err != pkg.ErrInvalidSecret {
			tests.FailedWithError(err, "Expected to have received error 'ErrInvalidSecret'")
		}
		tests.Passed("Should have received error 'ErrInvalidSecret'")
	}

	t.Logf("\tWhen parsing message with no OP")
	{
		raw := []byte("wreckoge#232uFR5  welcome to the league of heroes!")
		_, err := parser.Parse(raw)
		if err == nil {
			tests.Failed("Should have failed to parse message")
		}
		tests.Passed("Should have failed to parse message")

		if err != pkg.ErrInvalidOpName {
			tests.FailedWithError(err, "Expected to have received error 'ErrInvalidOpName'")
		}
		tests.Passed("Should have received error 'ErrInvalidOpName'")
	}

	t.Logf("\tWhen parsing message with correct secret")
	{
		raw := []byte("wreckage#232uFR5   MSG   welcome to the league!")
		op, err := parser.Parse(raw)
		if err != nil {
			tests.FailedWithError(err, "Should have successfully parsed message")
		}
		tests.Passed("Should have successfully parsed message")

		if op.Identity != "232uFR5" {
			tests.Info("Received: %+q", op.Identity)
			tests.Info("Expected: %+q", "232uFR5")
			tests.Failed("Should have successfully matched identity")
		}
		tests.Passed("Should have successfully matched identity")

		if op.Op != "MSG" {
			tests.Info("Received: %+q", op.Op)
			tests.Info("Expected: %+q", "MSG")
			tests.Failed("Should have successfully matched operation name")
		}
		tests.Passed("Should have successfully matched operation name")

		body := []byte("welcome to the league!")
		if bytes.Equal(op.Body, body) {
			tests.Info("Received: %+q", op.Body)
			tests.Info("Expected: %+q", body)
			tests.Failed("Should have successfully matched body")
		}
		tests.Passed("Should have successfully matched body")

		if bytes.Equal(op.Raw, raw) {
			tests.Info("Received: %+q", op.Raw)
			tests.Info("Expected: %+q", raw)
			tests.Failed("Should have successfully matched body")
		}
		tests.Passed("Should have successfully matched raw")
	}
}

package test

import (
	"github.com/bhbosman/go-sqlize/cmd"
	"os"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	cmd.RootCmd.SetArgs([]string{"./work/ifStatement.go"})
	sb := new(strings.Builder)
	cmd.RootCmd.SetOut(sb)
	err := cmd.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	//if sb.String() != "ifStatement.go" {
	//	require.Equal(t, sb.String(), "ifStatement.go")
	//}
}

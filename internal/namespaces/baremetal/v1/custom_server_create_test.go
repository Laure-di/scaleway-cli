package baremetal_test

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v1"
	baremetalSDK "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
)

// All test below should succeed to create an instance.
func Test_CreateServer(t *testing.T) {
	// Simple use cases
	t.Run("Simple", func(t *testing.T) {
		t.Run("Default", core.Test(&core.TestConfig{
			Commands: baremetal.GetCommands(),
			Cmd:      "scw baremetal server create zone=" + region + " type=" + offerName + " -w",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
			AfterFunc: core.ExecAfterCmd("scw baremetal server delete {{ .CmdResult.ID }} zone=" + region),
		},
		))

		t.Run("With name", core.Test(&core.TestConfig{
			Commands: baremetal.GetCommands(),
			Cmd:      "scw baremetal server create name=test-create-server-with-name zone=" + region + " type=" + offerName + " -w",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.Equal(t, "test-create-server-with-name", ctx.Result.(*baremetalSDK.Server).Name)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: core.ExecAfterCmd("scw baremetal server delete {{ .CmdResult.ID }} zone=" + region),
		}))

		t.Run("Tags", core.Test(&core.TestConfig{
			Commands: baremetal.GetCommands(),
			Cmd:      "scw baremetal server create tags.0=prod tags.1=blue zone=" + region + " type=" + offerName + " -w",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.Equal(t, "prod", ctx.Result.(*baremetalSDK.Server).Tags[0])
					assert.Equal(t, "blue", ctx.Result.(*baremetalSDK.Server).Tags[1])
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: core.ExecAfterCmd("scw baremetal server delete {{ .CmdResult.ID }} zone=" + region),
		}))
	})
}

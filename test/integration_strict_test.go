package test_test

import (
	"testing"

	"github.com/gruntwork-io/terragrunt/internal/strict"
	"github.com/gruntwork-io/terragrunt/test/helpers"
	"github.com/gruntwork-io/terragrunt/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStrictMode(t *testing.T) {
	t.Parallel()

	helpers.CleanupTerraformFolder(t, testFixtureEmptyState)

	tc := []struct {
		name           string
		controls       []string
		strictMode     bool
		expectedStderr string
		expectedError  error
	}{
		{
			name:           "plan-all",
			controls:       []string{},
			strictMode:     false,
			expectedStderr: "The `plan-all` command is deprecated and will be removed in a future version. Use `terragrunt run-all plan` instead.",
			expectedError:  nil,
		},
		{
			name:           "plan-all with plan-all strict control",
			controls:       []string{"plan-all"},
			strictMode:     false,
			expectedStderr: "",
			expectedError:  strict.StrictControls[strict.PlanAll].Error,
		},
		{
			name:           "plan-all with multiple strict controls",
			controls:       []string{"plan-all", "apply-all"},
			strictMode:     false,
			expectedStderr: "",
			expectedError:  strict.StrictControls[strict.PlanAll].Error,
		},
		{
			name:           "plan-all with strict mode",
			controls:       []string{},
			strictMode:     true,
			expectedStderr: "",
			expectedError:  strict.StrictControls[strict.PlanAll].Error,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpEnvPath := helpers.CopyEnvironment(t, testFixtureEmptyState)
			rootPath := util.JoinPath(tmpEnvPath, testFixtureEmptyState)

			args := "--terragrunt-non-interactive --terragrunt-log-level debug --terragrunt-working-dir " + rootPath
			if tt.strictMode {
				args = "--strict-mode " + args
			}

			for _, control := range tt.controls {
				args = " --strict-control " + control + " " + args
			}

			_, stderr, err := helpers.RunTerragruntCommandWithOutput(t, "terragrunt plan-all "+args)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			assert.Contains(t, stderr, tt.expectedStderr)
		})
	}
}

func TestRootTerragruntHCLStrictMode(t *testing.T) {
	t.Parallel()

	helpers.CleanupTerraformFolder(t, testFixtureFindParent)

	tc := []struct {
		name           string
		controls       []string
		strictMode     bool
		expectedStderr string
		expectedError  error
	}{
		{
			name:           "root terragrunt.hcl",
			strictMode:     false,
			expectedStderr: strict.StrictControls[strict.RootTerragruntHCL].Warning,
		},
		{
			name:          "root terragrunt.hcl with root-terragrunt-hcl strict control",
			controls:      []string{"root-terragrunt-hcl"},
			strictMode:    false,
			expectedError: strict.StrictControls[strict.RootTerragruntHCL].Error,
		},
		{
			name:          "root terragrunt.hcl with strict mode",
			strictMode:    true,
			expectedError: strict.StrictControls[strict.RootTerragruntHCL].Error,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpEnvPath := helpers.CopyEnvironment(t, testFixtureFindParent)
			rootPath := util.JoinPath(tmpEnvPath, testFixtureFindParent, "app")

			args := "--terragrunt-non-interactive --terragrunt-log-level debug --terragrunt-working-dir " + rootPath
			if tt.strictMode {
				args = "--strict-mode " + args
			}

			for _, control := range tt.controls {
				args = " --strict-control " + control + " " + args
			}

			_, stderr, err := helpers.RunTerragruntCommandWithOutput(t, "terragrunt plan "+args)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Contains(t, stderr, tt.expectedStderr)
		})
	}
}

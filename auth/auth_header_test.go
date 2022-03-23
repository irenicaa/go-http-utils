package auth

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type env struct {
	name     string
	value    string
	isNotSet bool
}

func updateEnvs(envs []env) error {
	for _, env := range envs {
		if env.isNotSet {
			if err := os.Unsetenv(env.name); err != nil {
				return fmt.Errorf("unable to unset an environment variable: %w", err)
			}
		} else {
			if err := os.Setenv(env.name, env.value); err != nil {
				return fmt.Errorf("unable to set an environment variable: %w", err)
			}
		}
	}

	return nil
}

func unsetEnvs(envs []env) error {
	var envsForUnsetting []env
	for _, envInstance := range envs {
		envForUnsetting := env{name: envInstance.name, isNotSet: true}
		envsForUnsetting = append(envsForUnsetting, envForUnsetting)
	}

	return updateEnvs(envsForUnsetting)
}

func TestMakeBasicAuthHeader(t *testing.T) {
	type args struct {
		usernameEnv env
		passwordEnv env
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "without the username and the password",
			args: args{
				usernameEnv: env{name: "TEST_USERNAME", isNotSet: true},
				passwordEnv: env{name: "TEST_PASSWORD", isNotSet: true},
			},
			want: "",
		},
		{
			name: "with the empty username and the empty password",
			args: args{
				usernameEnv: env{name: "TEST_USERNAME", value: ""},
				passwordEnv: env{name: "TEST_PASSWORD", value: ""},
			},
			want: "",
		},
		{
			name: "with the nonempty username and the nonempty password",
			args: args{
				usernameEnv: env{name: "TEST_USERNAME", value: "username"},
				passwordEnv: env{name: "TEST_PASSWORD", value: "password"},
			},
			want: "Basic dXNlcm5hbWU6cGFzc3dvcmQ=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := updateEnvs([]env{tt.args.usernameEnv, tt.args.passwordEnv})
			require.NoError(t, err)
			defer unsetEnvs([]env{tt.args.usernameEnv, tt.args.passwordEnv})

			got := MakeBasicAuthHeader(
				tt.args.usernameEnv.name,
				tt.args.passwordEnv.name,
			)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMakeBearerAuthHeader(t *testing.T) {
	type args struct {
		tokenEnv env
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "without the token",
			args: args{
				tokenEnv: env{name: "TEST_TOKEN", isNotSet: true},
			},
			want: "",
		},
		{
			name: "with the empty token",
			args: args{
				tokenEnv: env{name: "TEST_TOKEN", value: ""},
			},
			want: "",
		},
		{
			name: "with the nonempty token",
			args: args{
				tokenEnv: env{name: "TEST_TOKEN", value: "token"},
			},
			want: "Bearer token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := updateEnvs([]env{tt.args.tokenEnv})
			require.NoError(t, err)
			defer unsetEnvs([]env{tt.args.tokenEnv})

			got := MakeBearerAuthHeader(tt.args.tokenEnv.name)

			assert.Equal(t, tt.want, got)
		})
	}
}

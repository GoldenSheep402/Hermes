package config

import "testing"

func TestGenYamlConfigWithMod(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	type args struct {
		path  string
		force bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test", args{"./5config.yaml", true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenYamlConfig(tt.args.path, tt.args.force); (err != nil) != tt.wantErr {
				t.Errorf("GenYamlConfigWithMod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

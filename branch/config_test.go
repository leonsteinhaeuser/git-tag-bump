package branch

import (
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "file found",
			args: args{
				path: "../config.yaml",
			},
			want: &Config{
				Major: Identifier{
					Branch: BranchIdentifier{
						Name: RegExIdentifier{
							RegEx: `^(feat|feature|enh|enhanc|enhancement|fix|bugfix|chore\([a-z0-9-]+\))!/`,
						},
					},
				},
				Minor: Identifier{
					Branch: BranchIdentifier{
						Name: RegExIdentifier{
							RegEx: "^(feat|feature)/",
						},
					},
				},
				Patch: Identifier{
					Branch: BranchIdentifier{
						Name: RegExIdentifier{
							RegEx: `^(enh|enhanc|enhancement|fix|bugfix|chore\([a-z0-9-]+\))/`,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "file not found",
			args: args{
				path: "testdata/config.yaml",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfig(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentifier_match(t *testing.T) {
	type fields struct {
		Branch BranchIdentifier
	}
	type args struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "match",
			fields: fields{
				Branch: BranchIdentifier{
					Name: RegExIdentifier{
						RegEx: "^[a-z]+!/",
					},
				},
			},
			args: args{
				value: "fix!/",
			},
			want: true,
		},
		{
			name: "not match",
			fields: fields{
				Branch: BranchIdentifier{
					Name: RegExIdentifier{
						RegEx: "^[a-z]+!/",
					},
				},
			},
			args: args{
				value: "fix/",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Identifier{
				Branch: tt.fields.Branch,
			}
			if got := i.match(tt.args.value); got != tt.want {
				t.Errorf("Identifier.match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBranchIdentifier_match(t *testing.T) {
	type fields struct {
		Name RegExIdentifier
	}
	type args struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "match",
			fields: fields{
				Name: RegExIdentifier{
					RegEx: "^[a-z]+!/",
				},
			},
			args: args{
				value: "fix!/",
			},
			want: true,
		},
		{
			name: "not match",
			fields: fields{
				Name: RegExIdentifier{
					RegEx: "^[a-z]+!/",
				},
			},
			args: args{
				value: "fix/",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := &BranchIdentifier{
				Name: tt.fields.Name,
			}
			if got := bi.match(tt.args.value); got != tt.want {
				t.Errorf("BranchIdentifier.match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegExIdentifier_match(t *testing.T) {
	type fields struct {
		Regex string
	}
	type args struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "match",
			fields: fields{
				Regex: "^[a-z]+!/",
			},
			args: args{
				value: "fix!/",
			},
			want: true,
		},
		{
			name: "not match",
			fields: fields{
				Regex: "^[a-z]+!/",
			},
			args: args{
				value: "fix/",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ri := RegExIdentifier{
				RegEx: tt.fields.Regex,
			}
			if got := ri.match(tt.args.value); got != tt.want {
				t.Errorf("RegExIdentifier.match() = %v, want %v", got, tt.want)
			}
		})
	}
}

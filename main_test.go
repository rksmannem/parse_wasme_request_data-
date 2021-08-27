package main

import "testing"

func TestProcessRequestDataForKey(t *testing.T) {
	type args struct {
		headers      [][2]string
		resolveSpecs ResolveSpecs
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid headers and resolve specs",
			args: args{
				headers: [][2]string{
					[2]string{"paramName", "key_in_header"},
					[2]string{"key_in_header", "1234"},
				},
				resolveSpecs: ResolveSpecs{
					ResolveSpec{
						ParamName:     "key_in_header",
						ParamLocation: InHeader,
					},
					ResolveSpec{
						ParamName:     "key_in_body",
						ParamLocation: InBody,
					},
				},
			},
			want:    "1234",
			wantErr: false,
		},
		{
			name: "paramName provided in header not found in resolve spec",
			args: args{
				headers: [][2]string{
					[2]string{"paramName", "key_not_in_header"},
					[2]string{"key_not_in_header", "1234"},
				},
				resolveSpecs: ResolveSpecs{
					ResolveSpec{
						ParamName:     "key_in_header",
						ParamLocation: InHeader,
					},
					ResolveSpec{
						ParamName:     "key_in_body",
						ParamLocation: InBody,
					},
				},
			},
			want:    EmptyString,
			wantErr: true,
		},
		{
			name: "paramName not found headers",
			args: args{
				headers: [][2]string{
					[2]string{"key_in_header", "1234"},
				},
				resolveSpecs: ResolveSpecs{
					ResolveSpec{
						ParamName:     "key_in_header",
						ParamLocation: InHeader,
					},
					ResolveSpec{
						ParamName:     "key_in_body",
						ParamLocation: InBody,
					},
				},
			},
			want:    EmptyString,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProcessRequestDataForKey(tt.args.headers, tt.args.resolveSpecs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessRequestDataForKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ProcessRequestDataForKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

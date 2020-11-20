package main

import "testing"

func TestFromFlagOrEnv(t *testing.T) {
	type args struct {
		flag string
		env  string
	}
	tests := []struct {
		name      string
		args      args
		wantValue string
		wantOk    bool
	}{
		{name: "1", args: args{flag: "test", env: "test"}, wantValue: "test", wantOk: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := FromFlagOrEnv(tt.args.flag, tt.args.env)
			if gotValue != tt.wantValue {
				t.Errorf("FromFlagOrEnv() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("FromFlagOrEnv() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

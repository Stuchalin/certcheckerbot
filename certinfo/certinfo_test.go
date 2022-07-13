package certinfo

import (
	"regexp"
	"testing"
)

func TestGetCertInfo(t *testing.T) {
	type args struct {
		URL            string
		printFullChain bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test Not URL",
			args: args{
				URL:            "NotValidURL",
				printFullChain: false,
			},
			want: "Cannot check cert from URL NotValidURL\\..*",
		},
		{
			name: "test valid google URL",
			args: args{
				URL:            "google.com",
				printFullChain: false,
			},
			want: "DNSNames: .*google\\.com.*\nIssuer Name: .*\nExpiry: \\d\\d\\d\\d-\\d\\d-\\d\\d\nCommon Name: .*\n",
		},
		{
			name: "test valid google URL full chain",
			args: args{
				URL:            "google.com",
				printFullChain: true,
			},
			want: "DNSNames: .*google\\.com.*\nIssuer Name: .*\nExpiry: \\d\\d\\d\\d-\\d\\d-\\d\\d\nCommon Name: (.|\n)*" +
				"DNSNames: \\[\\]\nIssuer Name: .*OU=Root CA.*\nExpiry: \\d\\d\\d\\d-\\d\\d-\\d\\d\nCommon Name: .*Root.*\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCertInfo(tt.args.URL, tt.args.printFullChain)
			res, err := regexp.MatchString(tt.want, got)
			if err != nil {
				t.Errorf("GetCertInfo() - regex error: %s", err)
			}
			if !res {
				t.Errorf("GetCertInfo() = %v, regex pattern = %v", got, tt.want)
			}
		})
	}
}

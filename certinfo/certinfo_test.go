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
		name       string
		args       args
		want       string
		wantErr    bool
		errMessage string
	}{
		{
			name: "test Not URL",
			args: args{
				URL:            "NotValidURL",
				printFullChain: false,
			},
			wantErr:    true,
			errMessage: "check certificate error - cannot check cert from URL NotValidURL\\..*",
		},
		{
			name: "test valid google URL",
			args: args{
				URL:            "google.com",
				printFullChain: false,
			},
			want: "DNSNames: .*google\\.com.*\nIssuer Name: .*\nExpiry: \\d\\d\\d\\d-\\d\\d-\\d\\d\nCommon Name: .*\n\n",
		},
		{
			name: "test valid google URL full chain",
			args: args{
				URL:            "google.com",
				printFullChain: true,
			},
			want: "DNSNames: .*google\\.com.*\nIssuer Name: .*\nExpiry: \\d\\d\\d\\d-\\d\\d-\\d\\d\nCommon Name: (.|\n)*" +
				"DNSNames: \\[\\]\nIssuer Name: .*OU=Root CA.*\nExpiry: \\d\\d\\d\\d-\\d\\d-\\d\\d\nCommon Name: .*Root.*\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCertInfo(tt.args.URL, tt.args.printFullChain)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetCertInfo() - expected Error, got nil")
				}
				res, err2 := regexp.MatchString(tt.errMessage, err.Error())
				if err2 != nil {
					t.Errorf("GetCertInfo() - regex error: %s", err2)
				}
				if !res {
					t.Errorf("GetCertInfo() = %v, regex pattern = %v", err, tt.errMessage)
				}
			}
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

func TestGetCertsInfo(t *testing.T) {
	type args struct {
		URLs           string
		printFullChain bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test single valid URL",
			args: args{
				URLs:           "google.com",
				printFullChain: false,
			},
			want: "DNSNames: .*google\\.com.*\nIssuer Name: .*\nExpiry: \\d\\d\\d\\d-\\d\\d-\\d\\d\nCommon Name: .*\n\n",
		},
		{
			name: "test few valid URLs",
			args: args{
				URLs:           "google.com github.com wikipedia.com",
				printFullChain: false,
			},
			want: "DNSNames: .*google\\.com.*\nIssuer Name: .*\nExpiry: \\d\\d\\d\\d-\\d\\d-\\d\\d\nCommon Name: .*\n\n" +
				"DNSNames: .*github\\.com.*\nIssuer Name: .*\nExpiry: \\d\\d\\d\\d-\\d\\d-\\d\\d\nCommon Name: .*\n\n" +
				"DNSNames: .*wikipedia\\.com.*\nIssuer Name: .*\nExpiry: \\d\\d\\d\\d-\\d\\d-\\d\\d\nCommon Name: .*\n\n",
		},
		{
			name: "test valid and fail URLs",
			args: args{
				URLs:           "google.com notValidDomain wikipedia.com",
				printFullChain: false,
			},
			want: "DNSNames: .*google\\.com.*\nIssuer Name: .*\nExpiry: \\d\\d\\d\\d-\\d\\d-\\d\\d\nCommon Name: .*\n\n" +
				"Cannot check cert from URL notValidDomain\\..*\n\n" +
				"DNSNames: .*wikipedia\\.com.*\nIssuer Name: .*\nExpiry: \\d\\d\\d\\d-\\d\\d-\\d\\d\nCommon Name: .*\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCertsInfo(tt.args.URLs, tt.args.printFullChain)
			res, err := regexp.MatchString(tt.want, got)
			if err != nil {
				t.Errorf("GetCertsInfo() - regex error: %s", err)
			}
			if !res {
				t.Errorf("GetCertsInfo() = %v, regex pattern = %v", got, tt.want)
			}
		})
	}
}

package pixiv

import "testing"

func Test_getValidPixivImageURL(t *testing.T) {
	tests := []struct {
		name     string
		imageURL string
		want     string
		wantErr  bool
	}{
		{
			name:     "Canonical Master URL",
			imageURL: "i.pximg.net/img-master/img/2024/01/01/00/00/00/12345_p0_master1200.jpg",
			want:     "https://i.pximg.net/img-master/img/2024/01/01/00/00/00/12345_p0_master1200.jpg",
		},
		{
			name:     "Protocol Relative, Uppercase Host And Extension (Normalized)",
			imageURL: "//I.PXIMG.NET/a.PNG",
			want:     "https://i.pximg.net/a.PNG",
		},
		{
			name:     "Query And Fragment Dropped By Rebuild",
			imageURL: "i.pximg.net/a.jpg?w=200#anchor",
			want:     "https://i.pximg.net/a.jpg",
		},

		// One case per guard. Rejects short-circuit, so each needs its own input.
		{
			name:     "Empty Input",
			imageURL: "",
			wantErr:  true,
		},
		{
			name:     "Explicit Scheme (Mixed Case)",
			imageURL: "HtTpS://i.pximg.net/a.jpg",
			wantErr:  true,
		},
		{
			name:     "Control Character (Unparseable)",
			imageURL: "i.pximg.net/a.jpg\rHost: evil.com",
			wantErr:  true,
		},
		{
			name:     "Suffix Confusion",
			imageURL: "i.pximg.net.evil.com/a.jpg",
			wantErr:  true,
		},
		{
			name:     "Valid Host As Userinfo",
			imageURL: "i.pximg.net@evil.com/a.jpg", // real host is evil.com
			wantErr:  true,
		},
		{
			name:     "Custom Port",
			imageURL: "i.pximg.net:8080/a.jpg", // Hostname() strips the port
			wantErr:  true,
		},
		{
			name:     "Encoded Question Mark Fakes Extension",
			imageURL: "i.pximg.net/evil.svg%3f.jpg", // parsed.Path decodes to ".jpg"
			wantErr:  true,
		},
		{
			name:     "Traversal Onto Disallowed Extension",
			imageURL: "i.pximg.net/a.jpg/../evil.svg",
			wantErr:  true,
		},
		{
			name:     "Disallowed Extension",
			imageURL: "i.pximg.net/evil.svg",
			wantErr:  true,
		},
		{
			name:     "Empty Path (Root)",
			imageURL: "i.pximg.net",
			wantErr:  true,
		},
		{
			name:     "Raw Space Encoded By Parser",
			imageURL: "i.pximg.net/a b.jpg", // EscapedPath() becomes "a%20b.jpg"
			wantErr:  true,
		},
		{
			name:     "Raw Null Byte Injection",
			imageURL: "i.pximg.net/a.jpg\x00.png",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getValidPixivImageURL(tt.imageURL)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("got %q, want error", got)
				}

				if got != "" {
					t.Errorf("got %q alongside error, want empty string", got)
				}

				return
			}

			if err != nil {
				t.Fatalf("unwanted error: %v", err)
			}

			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

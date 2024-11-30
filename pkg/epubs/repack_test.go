package epubs

import (
	"os"
	"testing"
)

func Test_RepackEpub(t *testing.T) {
	type args struct {
		existingEpubPath string
		updatedEpubPath  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "repackZip",
			args: args{
				existingEpubPath: "texts/dahl.epub",
				updatedEpubPath:  "texts/dahl_rep.epub",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RepackEpub(tt.args.existingEpubPath, tt.args.updatedEpubPath, nil); (err != nil) != tt.wantErr {
				t.Errorf("repackZip() error = %v, wantErr %v", err, tt.wantErr)
			}
			// check if the file was created
			if _, err := os.Stat(tt.args.updatedEpubPath); os.IsNotExist(err) {
				t.Errorf("repackZip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

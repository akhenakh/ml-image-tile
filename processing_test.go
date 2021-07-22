package main

import (
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path"
	"testing"

	log "github.com/go-kit/kit/log"
)

func Test_processImageBimg(t *testing.T) {
	testFile := "./testdata/testimg.png"

	type args struct {
		logger  log.Logger
		srcPath string
		resize  int
		width   int
		height  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"happy path", args{
			log.NewLogfmtLogger(os.Stdout),
			testFile,
			0,
			224,
			224,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := os.TempDir()
			t.Log("tmpDir", path.Join(tmpDir, "testdata"))

			os.Mkdir(path.Join(tmpDir, "testdata"), 0755)

			if err := processImageBimg(
				tt.args.logger,
				tt.args.srcPath,
				tmpDir,
				tt.args.resize,
				tt.args.width,
				tt.args.height,
			); (err != nil) != tt.wantErr {
				t.Errorf("processImageBimg() error = %v, wantErr %v", err, tt.wantErr)
			}

			// clean up only if we succeeded
			// os.RemoveAll(tmpDir)
		})
	}
}

package main

import "testing"

func TestShouldSkipFile(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		// Binary extensions — should be skipped
		{"exe file", "terraform-provider.exe", true},
		{"dll file", "lib/native.dll", true},
		{"so file", "/usr/lib/libssl.so", true},
		{"bin file", "tools/runner.bin", true},
		{"object file", "build/main.o", true},
		{"archive file", "lib/libfoo.a", true},
		{"dylib file", "lib/libbar.dylib", true},
		{"pyc file", "__pycache__/mod.pyc", true},
		{"class file", "com/app/Main.class", true},
		{"jar file", "libs/dep.jar", true},
		{"zip file", "archive.zip", true},
		{"tar file", "backup.tar", true},
		{"gz file", "backup.tar.gz", true},
		{"png file", "assets/logo.png", true},
		{"jpg file", "assets/photo.jpg", true},
		{"gif file", "assets/anim.gif", true},
		{"pdf file", "docs/manual.pdf", true},
		{"ico file", "favicon.ico", true},
		{"woff file", "fonts/font.woff", true},
		{"woff2 file", "fonts/font.woff2", true},
		{"ttf file", "fonts/font.ttf", true},

		// Case insensitivity
		{"EXE uppercase", "app.EXE", true},
		{"Dll mixed case", "lib/Native.Dll", true},

		// Source/text files — not skipped
		{"go file", "main.go", false},
		{"json file", "config.json", false},
		{"txt file", "notes.txt", false},
		{"yaml file", "config.yaml", false},
		{"py file", "script.py", false},
		{"env file", ".env", false},
		{"sh file", "deploy.sh", false},
		{"md file", "README.md", false},
		{"no extension", "Makefile", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSkipFile(tt.path)
			if got != tt.want {
				t.Errorf("shouldSkipFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}


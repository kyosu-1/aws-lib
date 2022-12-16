package file

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetRelativeFilePaths(t *testing.T) {
	// テスト用のディレクトリを作成する(tmpDirを利用)
	// テスト用のファイルを作成する
	// テスト用のファイルの相対パスを取得する

	testCases := []struct {
		name     string
		dirPath  string
		createRelativeFilePaths []string
		expected []string
		wantErr  bool
	}{
		{
			name:     "正常系",
			dirPath:  "testDir",
			createRelativeFilePaths: []string{
				"testFile1",
				"testFile2",
			},
			expected: []string{
				"testFile1",
				"testFile2",
			},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// テスト用のディレクトリを作成する(tmpDirを利用)
			// テスト用のファイルを作成する
			// テスト用のファイルの相対パスを取得する
			// 取得した結果が期待値と一致するか確認する

			// テスト用のディレクトリを作成する(tmpDirを利用)
			tmpDir, err := os.MkdirTemp("./", tc.dirPath)
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			// テスト用のファイルを作成する
			for _, relativeFilePath := range tc.createRelativeFilePaths {
				_, err := os.Create(tmpDir + "/" + relativeFilePath)
				if err != nil {
					t.Fatal(err)
				}
			}

			// テスト用のファイルの相対パスを取得する
			relativeFilePaths, err := GetRelativeFilePaths(tmpDir)

			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.expected, relativeFilePaths)
		})
	}
}

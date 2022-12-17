package file

import (
	"os"
	"path/filepath"
)

// 指定したディレクトリ内のファイルの相対パスを取得する
// 例えば、指定したディレクトリが"/home/user"で、
// ディレクトリ内にあるファイルが"/home/user/testFile1"と"/home/user/testFile2"の場合、
// []string{"user/testFile1", "user/testFile2"}を返す
func GetRelativeFilePaths(dirPath string) ([]string, error) {
	var relativeFilePaths []string

	// ディレクトリ内のファイルの絶対パスを取得する
	absoluteFilePaths, err := getAbsoluteFilePaths(dirPath)
	if err != nil {
		return nil, err
	}

	// 絶対パスを相対パスに変換する
	for _, absoluteFilePath := range absoluteFilePaths {
		relativeFilePath, err := filepath.Rel(dirPath, absoluteFilePath)
		if err != nil {
			return nil, err
		}
		relativeFilePaths = append(relativeFilePaths, relativeFilePath)
	}

	return relativeFilePaths, nil
}

// 指定したディレクトリ内のファイルの絶対パスを取得する
func getAbsoluteFilePaths(dirPath string) ([]string, error) {
	var absoluteFilePaths []string

	// ディレクトリ内のファイルの絶対パスを取得する
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			absoluteFilePaths = append(absoluteFilePaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return absoluteFilePaths, nil
}

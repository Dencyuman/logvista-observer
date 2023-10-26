package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateLogvistaDir(dirPath string) (string, error) {
	logvistaDirPath := filepath.Join(dirPath, ".logvista")
	// ディレクトリが存在しない場合は作成する
	// os.Statの実行とos.IsNotExistを同時に行うことで変数errのスコープがif文内に限定される
	if _, err := os.Stat(logvistaDirPath); os.IsNotExist(err) {
		err := os.Mkdir(logvistaDirPath, 0755)
		if err != nil {
			return "", fmt.Errorf("ディレクトリの作成に失敗しました: %v", err)
		}
		fmt.Printf("%s ディレクトリを作成しました\n", logvistaDirPath)
	} else if err != nil {
		return "", fmt.Errorf("ディレクトリの確認中にエラーが発生しました: %v", err)
	}
	return logvistaDirPath, nil
}

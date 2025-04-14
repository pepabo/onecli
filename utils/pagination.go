package utils

// Paginate はページネーションをサポートする関数を実行し、全結果を取得します
func Paginate[T any](fetch func(offset int) ([]T, error)) ([]T, error) {
	var allResults []T
	offset := 0

	for {
		results, err := fetch(offset)
		if err != nil {
			return nil, err
		}

		// 結果が空の場合は終了
		if len(results) == 0 {
			break
		}

		allResults = append(allResults, results...)
		offset += len(results)
	}

	return allResults, nil
}

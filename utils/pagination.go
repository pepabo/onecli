package utils

func Paginate[T any](fetch func(page int) ([]T, error), limit int) ([]T, error) {
	var allResults []T
	page := 1

	for {
		results, err := fetch(page)
		if err != nil {
			return nil, err
		}

		// 結果が空の場合は終了
		if len(results) == 0 {
			break
		}

		allResults = append(allResults, results...)

		// 取得した結果がlimitより少ない場合は終了
		if len(results) < limit {
			break
		}

		page++
	}

	return allResults, nil
}

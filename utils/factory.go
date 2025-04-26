package utils

func VndbRequestData(brand string) map[string]interface{} {
	return map[string]interface{}{
		"filters": []interface{}{
			"developer",
			"=",
			[]string{
				"search", "=", brand,
			},
		},
		"fields": "title, alttitle, rating, image.url",
	}
}

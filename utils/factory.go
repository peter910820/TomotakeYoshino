package utils

// create struct to use brand name to search producer id for vndb
func VndbProducerRequest(brand string) map[string]interface{} {
	return map[string]interface{}{
		"filters": []string{
			"search", "=", brand,
		},
		"results": 100,
		"fields":  "id, name",
	}
}

func VndbVnRequest(brandId string) map[string]interface{} {
	return map[string]interface{}{
		"filters": []interface{}{
			"developer",
			"=",
			[]string{
				"search", "=", brandId,
			},
		},
		"results": 100,
		"fields":  "title, alttitle, rating, image.url",
	}
}

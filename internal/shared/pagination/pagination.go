package pagination

import "strconv"

type Params struct {
	Limit  int32
	Offset int32
}

func Parse(limitRaw, offsetRaw string) Params {
	limit := int32(20)
	offset := int32(0)

	if parsed, err := strconv.Atoi(limitRaw); err == nil && parsed > 0 && parsed <= 100 {
		limit = int32(parsed)
	}
	if parsed, err := strconv.Atoi(offsetRaw); err == nil && parsed >= 0 {
		offset = int32(parsed)
	}

	return Params{Limit: limit, Offset: offset}
}

package utils

import "math"

type OffsetPaginationMeta struct {
	Page        int `json:"page"`
	Limit       int `json:"limit"`
	PageCount   int `json:"pageCount"`
	RecordCount int `json:"recordCount"`
}

type GeneratedOffsetPaginationMeta struct {
	Limit        int                  `json:"limit"`
	Offset       int                  `json:"offset"`
	OffPageLimit bool                 `json:"offPageLimit"`
	Meta         OffsetPaginationMeta `json:"meta"`
}

type PaginationOffsetMetaOptions struct {
	Page        int
	Limit       int
	RecordCount int
}

func GenerateOffsetPaginationMeta(options PaginationOffsetMetaOptions) GeneratedOffsetPaginationMeta {
	page := max(1, options.Page)

	pageCount := int(math.Ceil(float64(options.RecordCount) / float64(options.Limit)))

	offset := (page - 1) * options.Limit

	offPageLimit := page > pageCount

	return GeneratedOffsetPaginationMeta{
		Limit:        options.Limit,
		Offset:       offset,
		OffPageLimit: offPageLimit,
		Meta: OffsetPaginationMeta{
			Page:        page,
			Limit:       options.Limit,
			PageCount:   pageCount,
			RecordCount: options.RecordCount,
		},
	}
}

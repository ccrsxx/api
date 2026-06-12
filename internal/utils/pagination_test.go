package utils_test

import (
	"testing"

	"github.com/ccrsxx/api/internal/utils"
)

func TestGenerateOffsetPaginationMeta(t *testing.T) {
	t.Run("First Page", func(t *testing.T) {
		result := utils.GenerateOffsetPaginationMeta(utils.PaginationOffsetMetaOptions{
			Page:        1,
			Limit:       10,
			RecordCount: 42,
		})

		if result.Limit != 10 {
			t.Fatalf("got limit %d, want 10", result.Limit)
		}

		if result.Offset != 0 {
			t.Fatalf("got offset %d, want 0", result.Offset)
		}

		if result.OffPageLimit {
			t.Error("got offPageLimit true, want false")
		}

		if result.Meta.Page != 1 {
			t.Fatalf("got meta.page %d, want 1", result.Meta.Page)
		}

		if result.Meta.Limit != 10 {
			t.Fatalf("got meta.limit %d, want 10", result.Meta.Limit)
		}

		if result.Meta.PageCount != 5 {
			t.Fatalf("got meta.pageCount %d, want 5", result.Meta.PageCount)
		}

		if result.Meta.RecordCount != 42 {
			t.Errorf("got meta.recordCount %d, want 42", result.Meta.RecordCount)
		}
	})

	t.Run("Middle Page Offset", func(t *testing.T) {
		result := utils.GenerateOffsetPaginationMeta(utils.PaginationOffsetMetaOptions{
			Page:        3,
			Limit:       10,
			RecordCount: 42,
		})

		if result.Offset != 20 {
			t.Fatalf("got offset %d, want 20", result.Offset)
		}

		if result.OffPageLimit {
			t.Fatal("got offPageLimit true, want false")
		}

		if result.Meta.Page != 3 {
			t.Errorf("got meta.page %d, want 3", result.Meta.Page)
		}
	})

	t.Run("PageCount Ceiling Division", func(t *testing.T) {
		result := utils.GenerateOffsetPaginationMeta(utils.PaginationOffsetMetaOptions{
			Page:        1,
			Limit:       10,
			RecordCount: 41,
		})

		// ceil(41/10) = 5, not 4
		if result.Meta.PageCount != 5 {
			t.Errorf("got meta.pageCount %d, want 5", result.Meta.PageCount)
		}
	})

	t.Run("Exact Division PageCount", func(t *testing.T) {
		result := utils.GenerateOffsetPaginationMeta(utils.PaginationOffsetMetaOptions{
			Page:        1,
			Limit:       10,
			RecordCount: 40,
		})

		// ceil(40/10) = 4, no extra page
		if result.Meta.PageCount != 4 {
			t.Errorf("got meta.pageCount %d, want 4", result.Meta.PageCount)
		}
	})

	t.Run("OffPageLimit True When Page Exceeds PageCount", func(t *testing.T) {
		result := utils.GenerateOffsetPaginationMeta(utils.PaginationOffsetMetaOptions{
			Page:        6,
			Limit:       10,
			RecordCount: 42,
		})

		// pageCount is 5, page 6 is out of bounds
		if !result.OffPageLimit {
			t.Error("got offPageLimit false, want true")
		}
	})

	t.Run("Page Clamps to 1 When Zero", func(t *testing.T) {
		result := utils.GenerateOffsetPaginationMeta(utils.PaginationOffsetMetaOptions{
			Page:        0,
			Limit:       10,
			RecordCount: 42,
		})

		if result.Meta.Page != 1 {
			t.Fatalf("got meta.page %d, want 1", result.Meta.Page)
		}

		if result.Offset != 0 {
			t.Errorf("got offset %d, want 0", result.Offset)
		}
	})

	t.Run("Page Clamps to 1 When Negative", func(t *testing.T) {
		result := utils.GenerateOffsetPaginationMeta(utils.PaginationOffsetMetaOptions{
			Page:        -5,
			Limit:       10,
			RecordCount: 42,
		})

		if result.Meta.Page != 1 {
			t.Fatalf("got meta.page %d, want 1", result.Meta.Page)
		}

		if result.Offset != 0 {
			t.Errorf("got offset %d, want 0", result.Offset)
		}
	})

	t.Run("Zero RecordCount", func(t *testing.T) {
		result := utils.GenerateOffsetPaginationMeta(utils.PaginationOffsetMetaOptions{
			Page:        1,
			Limit:       10,
			RecordCount: 0,
		})

		if result.Meta.PageCount != 0 {
			t.Fatalf("got meta.pageCount %d, want 0", result.Meta.PageCount)
		}

		if result.Meta.RecordCount != 0 {
			t.Errorf("got meta.recordCount %d, want 0", result.Meta.RecordCount)
		}
	})
}

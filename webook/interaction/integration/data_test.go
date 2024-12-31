package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"learn_go/webook/interaction/integration/startup"
	"learn_go/webook/internal/domain"
	"testing"
	"time"
)

func TestData(t *testing.T) {
	db := startup.NewDB()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	now := time.Now()

	err := db.WithContext(ctx).Create(&domain.Article{
		ID:      1,
		Title:   "article1",
		Content: "this is article 1",
		Author: domain.Author{
			ID: 1001,
		},
		Status: domain.ArticleStatusPublished,
		CTime:  now,
	}).Error
	assert.NoError(t, err)
}

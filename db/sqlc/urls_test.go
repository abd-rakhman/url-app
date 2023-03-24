package db

import (
	"context"
	"testing"

	"github.com/abd-rakhman/url-app/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateURL(t *testing.T) {
	// test for user given hashID
	randomUniqueHashID := utils.RandomString(16)
	url, err := testQueries.CreateUrl(context.Background(), CreateUrlParams{
		HashID: randomUniqueHashID,
		Url:    "test_url_link",
	})
	require.NoError(t, err)
	require.NotEmpty(t, url)
	require.Equal(t, randomUniqueHashID, url.HashID)
	require.Equal(t, "test_url_link", url.Url)

	// test for duplicate of the user given hashID
	_, err = testQueries.CreateUrl(context.Background(), CreateUrlParams{
		HashID: randomUniqueHashID,
		Url:    "test_url_link",
	})
	require.Error(t, err)
}

func TestGetURLByHashId(t *testing.T) {
	// create new url
	randomUniqueHashID := utils.RandomString(16)
	url, err := testQueries.CreateUrl(context.Background(), CreateUrlParams{
		HashID: randomUniqueHashID,
		Url:    "test_url_link",
	})
	require.NoError(t, err)
	require.NotEmpty(t, url)
	require.Equal(t, randomUniqueHashID, url.HashID)
	require.Equal(t, "test_url_link", url.Url)

	// test if the url is exists
	url, err = testQueries.GetUrlByHashId(context.Background(), randomUniqueHashID)
	require.NoError(t, err)
	require.NotEmpty(t, url)
	require.Equal(t, randomUniqueHashID, url.HashID)
	require.Equal(t, "test_url_link", url.Url)

	// test if the random url is not exists
	_, err = testQueries.GetUrlByHashId(context.Background(), utils.RandomString(16))
	require.Error(t, err)
}

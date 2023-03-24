package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/abd-rakhman/url-app/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateNewURLTx(t *testing.T) {
	store := NewStore(testDB)

	// test for user given hashID
	randomUniqueHashID := utils.RandomString(16)
	url, err := store.CreateNewURLTx(context.Background(), CreateNewURLRequest{
		HashID: randomUniqueHashID,
		URL:    "test_url_link",
	})
	require.NoError(t, err)
	require.NotEmpty(t, url)
	require.Equal(t, randomUniqueHashID, url.HashID)
	require.Equal(t, "test_url_link", url.Url)

	// test for duplicate of the user given hashID
	_, err = store.CreateNewURLTx(context.Background(), CreateNewURLRequest{
		HashID: randomUniqueHashID,
		URL:    "test_url_link",
	})
	require.Error(t, err)

	// test for many random hashID
	n := 100

	urlChan := make(chan Url, n)
	errChan := make(chan error, n)

	for i := 0; i < n; i++ {
		go func(i int) {
			url, err := store.CreateNewURLTx(context.Background(), CreateNewURLRequest{
				URL: fmt.Sprintf("test_url_link%d", i),
			})
			urlChan <- url
			errChan <- err
		}(i)
	}

	for i := 0; i < n; i++ {
		url := <-urlChan
		err := <-errChan
		require.NoError(t, err)
		require.NotEmpty(t, url)

		// test if the url is exists
		getUrl, err := store.GetUrlByHashId(context.Background(), url.HashID)
		require.NoError(t, err)
		require.NotEmpty(t, getUrl)
		require.Equal(t, url.Url, getUrl.Url)
		require.Equal(t, url.HashID, getUrl.HashID)
	}

	// test with several goroutines on the same hashID
	n = 10
	randomUniqueHashID = utils.RandomString(16)
	for i := 0; i < n; i++ {
		go func() {
			_, err := store.CreateNewURLTx(context.Background(), CreateNewURLRequest{
				HashID: randomUniqueHashID,
				URL:    "test_url_link",
			})
			errChan <- err
		}()
	}

	k := 0
	for i := 0; i < n; i++ {
		err := <-errChan
		if err != nil {
			k++
		}
	}
	require.Equal(t, n-1, k)

	// check the last several urls on the same hashID
	url, err = store.GetUrlByHashId(context.Background(), randomUniqueHashID)
	require.NoError(t, err)
	require.NotEmpty(t, url)
	require.Equal(t, randomUniqueHashID, url.HashID)
	require.Equal(t, "test_url_link", url.Url)

	// check for max length of hashID
	_, err = store.CreateNewURLTx(context.Background(), CreateNewURLRequest{
		HashID: utils.RandomString(20),
		URL:    "test_url_link",
	})
	require.Error(t, err)

	// check for expiresAt with existingHashID
	currentTime := time.Now().Unix()
	fmt.Println(currentTime)
	randomUniqueHashID = utils.RandomString(16)
	result, err := store.CreateNewURLTx(context.Background(), CreateNewURLRequest{
		HashID:    randomUniqueHashID,
		URL:       "test_url_link",
		ExpiresAt: currentTime + 3600,
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, randomUniqueHashID, result.HashID)
	require.Equal(t, "test_url_link", result.Url)
	require.Equal(t, currentTime+3600, result.ExpiresAt.Unix())

	// check for expiresAt with randomHashID
	result, err = store.CreateNewURLTx(context.Background(), CreateNewURLRequest{
		URL:       "test_url_link",
		ExpiresAt: currentTime + 3600,
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, currentTime+3600, result.ExpiresAt.Unix())

	// now, I should be able to extract it
	getUrl, err := store.GetUrlByHashId(context.Background(), result.HashID)
	require.NoError(t, err)
	require.NotEmpty(t, getUrl)
	require.Equal(t, result.Url, getUrl.Url)
	require.Equal(t, result.HashID, getUrl.HashID)
	require.Equal(t, result.ExpiresAt.Unix(), getUrl.ExpiresAt.Unix())

	// check for expired expiresAt
	result, err = store.CreateNewURLTx(context.Background(), CreateNewURLRequest{
		URL:       "test_url_link",
		ExpiresAt: currentTime + 10,
	})
	require.NoError(t, err)
	time.Sleep(10 * time.Second)

	// now, I should not be able to extract it
	// getUrl, err = store.GetUrlByHashId(context.Background(), result.HashID)
	// require.Error(t, err)
	// require.Empty(t, getUrl)
}

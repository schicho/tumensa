package tumensa

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

func buildCacheFilePath(timestamp time.Time) string {
	year, week := timestamp.ISOWeek()
	filename := fmt.Sprintf("tumensa-gql-week-%v-%v", year, week)
	return path.Join(os.TempDir(), filename)
}

func GetCachedGQLResponse(timestamp time.Time) (cacheFile *os.File, ok bool) {
	cacheFilePath := buildCacheFilePath(timestamp)

	cacheFile, err := os.Open(cacheFilePath)
	if err != nil {
		return nil, false
	}
	return cacheFile, true
}

func CacheGQLResponse(r io.Reader, timestamp time.Time) error {
	cacheFilePath := buildCacheFilePath(timestamp)

	cacheFile, err := os.Create(cacheFilePath)
	if err != nil {
		return err
	}
	defer cacheFile.Close()

	_, err = io.Copy(cacheFile, r)
	if err != nil {
		return err
	}
	return nil
}

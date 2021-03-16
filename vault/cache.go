package vault

import (
	"fmt"
	"time"

	"github.com/maahsome/vault-view/common"
	"github.com/sirupsen/logrus"
)

const expireMinutes=1

// Cache - the interface
type Cache struct {
	VaultClient Client
	CachePaths  map[string]CachePath
	CacheDatas  map[string]CacheData
}

// CachePath - hold items read in cache
type CachePath struct {
	CacheTime time.Time
	Paths     map[string]Paths
}

// CacheData - key/value data cache
type CacheData struct {
	CacheTime time.Time
	Data      DataRecord
}

var cachePaths map[string]CachePath
var cacheDatas map[string]CacheData

// NewCache - Create a new caching instance
func NewCache(client Client) *Cache {

	cachePaths := make(map[string]CachePath, 0)
	cacheDatas := make(map[string]CacheData, 0)
	return &Cache{VaultClient: client, CachePaths: cachePaths, CacheDatas: cacheDatas}

}

// CachePathExists - Check for the existance of a path in cache
func (c *Cache) CachePathExists(path string) bool {
	if _, ok := c.CachePaths[path]; ok {
		return true
	}
	return false
}

// UpdateCachePath - Add a path item to the cache list
func (c *Cache) UpdateCachePath(path string, paths map[string]Paths) (bool, error) {

	c.CachePaths[path] = CachePath{
		CacheTime: time.Now(),
		Paths:     paths,
	}
	return true, nil
}

// CacheDataExist - Check to see if we have the key/value pairs in cache
// We added the max parameter in order to support testing, as we can pass in a 0
func (c *Cache) CacheDataExist(path string, max float64) bool {
	if val, ok := c.CacheDatas[path]; ok {
		diff := time.Now().Sub(val.CacheTime)
		if diff.Minutes() > max {
			common.Logger.WithFields(logrus.Fields{
				"unit":     "cache",
				"function": "existance",
			}).Info(fmt.Sprintf("CACHE/expired for: %s", path))
			return false
		}
		return true
	}
	return false
}

// GetCacheData - Fetch the key/value data
func (c *Cache) GetCacheData(path string) DataRecord {
	if _, ok := c.CacheDatas[path]; ok {
		// base key exists, add the values
		return c.CacheDatas[path].Data
	}
	return DataRecord{}
}

// UpdateCacheData - Add key/value data to the cache
func (c *Cache) UpdateCacheData(path string, data DataRecord) (bool, error) {

	// base key exists, add the values
	c.CacheDatas[path] = CacheData{
		CacheTime: time.Now(),
		Data:      data,
	}

	return true, nil
}

// PreloadFolderPaths - called as a goroutine for each of the Folders at the
// current level.
func (c *Cache) PreloadFolderPaths(folderPath string) error {

	common.Logger.WithFields(logrus.Fields{
		"unit":        "cache",
		"function":    "preload",
		"folder_path": folderPath,
	}).Debug("Pre-Cache Keys for Downlevel Folders")
	pathList, err := c.VaultClient.GetPaths(folderPath)
	if err != nil {
		common.Logger.WithFields(logrus.Fields{
			"unit":       "cache",
			"function":   "preload",
			"folder_key": folderPath,
		}).Warn("Unable to list keys in keypath")
	}

	vaultPaths := make(map[string]Paths)
	for _, val := range pathList {
		if val.Type == vaultData {
			// Add this key to the list
			vaultPaths[fmt.Sprintf("%s%s", folderPath, val.Path)] = Paths{
				Type:    val.Type,
				Path:    val.Path,
				Parent:  folderPath,
				Version: 0,
			}
		}
	}
	c.PreloadPaths(vaultPaths)

	return nil
}

// PreloadPaths - called as a goroutine to pre-populate key/values for each KEY item
func (c *Cache) PreloadPaths(paths map[string]Paths) error {
	var data DataRecord
	var derr error

	common.Logger.WithFields(logrus.Fields{
		"unit":       "cache",
		"function":   "preload",
		"path_count": len(paths),
	}).Debug("Preload Paths")
	for k, v := range paths {
		if v.Type == vaultData {
			if !c.CacheDataExist(k, expireMinutes) {
				common.Logger.WithFields(logrus.Fields{
					"unit":     "cache",
					"function": "preload",
					"key":      k,
					"parent":   v.Parent,
					"path":     v.Path,
				}).Debug("Pre-Cache Data")

				data, derr = c.VaultClient.GetData(k)
				if derr != nil {
					common.Logger.WithFields(logrus.Fields{
						"unit":     "cache",
						"function": "preload",
						"key":      k,
					}).Warn("Failed to load key data")
				}
				if ok, err := c.UpdateCacheData(k, data); !ok {
					common.Logger.WithFields(logrus.Fields{
						"unit":     "cache",
						"function": "preload",
					}).WithError(err).Error("Bad UpdateCacheKeyValues")
				}
			}
		}
	}
	return nil
}

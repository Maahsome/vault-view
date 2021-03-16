package vault_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/maahsome/vault-view/vault"
	"github.com/maahsome/vault-view/common"
)

var _ = Describe("Cache", func() {
	var (
        mockVault Client
		cache *Cache
		newPaths map[string]Paths
		newData string
    )

	common.NewLogger("Warning", "")

	BeforeEach(func() {
		mockVault = NewVaultMock()
		cache = NewCache(mockVault)
		newPaths = make(map[string]Paths, 0)

		newPaths["/test/1"] = Paths{
			Type:     "Folder",
			Path:     "1",
			Parent:   "test",
			FullPath: "test/1",
		}

		newData = `{
			"data": {
				"data": {
					"value": "vault-test1"
				},
				"metadata": {
					"created_time": "2019-08-30T02:43:30.607941986Z",
					"deletion_time": "",
					"destroyed": false,
					"version": 1
				}
			}
		}`
	})

	Describe("NewCache", func() {

		It("creates a vault cache", func() {
			var dataRecord DataRecord

			marshErr := json.Unmarshal([]byte(newData), &dataRecord)
			if marshErr != nil {
				// Do we just fail the test here?
				Expect(marshErr).To(Equal(nil))
			}

			cache.UpdateCachePath("/test", newPaths)
			cache.UpdateCacheData("/test/1", dataRecord)

			Expect(len(cache.CachePaths)).To(Equal(1))
			Expect(len(cache.CacheDatas)).To(Equal(1))
		})
	})
	Describe("UpdateCachePath", func() {
		It("updates a cache path", func() {
		cache.UpdateCachePath("/test", newPaths)
		Expect(cache.CachePathExists("/test")).To(Equal(true))
		})
	})
	Describe("CachePathExists", func() {
		It("check path existance", func() {
			cache.UpdateCachePath("/test", newPaths)
			Expect(cache.CachePathExists("/noexist")).To(Equal(false))
			Expect(cache.CachePathExists("/test")).To(Equal(true))
		})
	})
	Describe("CacheDataExists", func() {
		It("check data existance", func() {
			var dataRecord DataRecord
			marshErr := json.Unmarshal([]byte(newData), &dataRecord)
			if marshErr != nil {
				Expect(marshErr).To(Equal(nil))
			}

			cache.UpdateCacheData("/test/1", dataRecord)

			Expect(cache.CacheDataExist("/test/1", 1)).To(Equal(true))
			Expect(cache.CacheDataExist("/test/2", 1)).To(Equal(false))
			Expect(cache.CacheDataExist("/test/1", 0)).To(Equal(false))
		})
	})
	Describe("GetCacheData", func() {
		It("get data from path", func() {
			var dataRecord DataRecord
			marshErr := json.Unmarshal([]byte(newData), &dataRecord)
			if marshErr != nil {
				Expect(marshErr).To(Equal(nil))
			}

			cache.UpdateCacheData("/test/1", dataRecord)

			Expect(cache.GetCacheData("/test/1")).To(Equal(dataRecord))
			Expect(cache.GetCacheData("/test/2")).To(Equal(DataRecord{}))
		})
	})
	Describe("PreloadData", func() {
		It("preload data", func() {

			newPaths["/test/2"] = Paths{
				Type:     "Folder",
				Path:     "2",
				Parent:   "/test",
				FullPath: "/test/2",
			}

			cache.PreloadPaths(newPaths)

			if vaultData, err := mockVault.GetData("/test/1"); err != nil {
				Expect(err).To(Equal(nil))
			} else {
				Expect(cache.GetCacheData("/test/1")).To(Equal(vaultData))
			}

			if vaultData, err := mockVault.GetData("/test/2"); err != nil {
				Expect(err).To(Equal(nil))
			} else {
				Expect(cache.GetCacheData("/test/2")).To(Equal(vaultData))
			}
		})
	})

})

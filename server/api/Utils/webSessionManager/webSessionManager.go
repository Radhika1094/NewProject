package webSessionManager

import (
	"MahaVan/MahaVanServer/api/constants"

	cachehelper "corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/cachehelper"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/logginghelper"
)

var SessionHelper cachehelper.CacheGCHelper

// var CachedLoggedInLearners cachehelper.CacheGCHelper

func Init() {
	logginghelper.LogDebug("started")

	SessionHelper = cachehelper.CacheGCHelper{}

	SessionHelper.Setup(constants.SESSION_HELPER_MAX_CACHE_ENTRIES, constants.EXPIRATION_DURATION)

	// CachedLoggedInLearners = cachehelper.CacheGCHelper{}

	// CachedLoggedInLearners.Setup(1000, 120)
	logginghelper.LogDebug("ended")
}

package service

import "sync"

// Concurrent session store using sync.Map
var sessionStore sync.Map

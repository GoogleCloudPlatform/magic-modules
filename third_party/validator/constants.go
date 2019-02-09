package google

import (
	"github.com/hashicorp/terraform/helper/mutexkv"
)

// Global MutexKV
var mutexKV = mutexkv.NewMutexKV()



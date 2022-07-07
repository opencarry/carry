package helper

import (
	v1 "github.com/opencarry/carry/pkg/apis/carry.i/v1"
	"github.com/opencarry/carry/pkg/util/sets"
)

var standardQuotaResources = sets.NewString(
	string(v1.ResourceCPU),
	string(v1.ResourceMemory),
)

func IsStandardQuotaResourceName(str string) bool {
	return standardQuotaResources.Has(str)
}

package detectors

import "shuru-hoja/pkg/types"

type Detector interface {
    Detect(info types.FileInfo) *types.ScanResult
}

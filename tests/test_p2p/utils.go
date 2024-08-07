package test_p2p

func getExpectedFileSize(chunks [][]byte) float64 {
	totalSize := 0
	for _, chunk := range chunks {
		totalSize += len(chunk)
	}
	return float64(totalSize)
}

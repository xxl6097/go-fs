//go:build !darwin

package utils

// AdhocSign 只有 darwin 需要重签名，其他平台空实现。
func AdhocSign(string) error {
	return nil
}

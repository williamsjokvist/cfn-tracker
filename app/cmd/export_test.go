package cmd

import "time"

// SetPollTimeoutForTest は外部テストから Poll の上限を短縮し、元の値へ戻す関数を返す。
// このファイルは _test.go のためテストビルドでのみコンパイルされ、
// 本番バイナリにテスト専用 API は含まれない。
func SetPollTimeoutForTest(timeout time.Duration) func() {
	original := pollTimeout
	pollTimeout = timeout
	return func() { pollTimeout = original }
}

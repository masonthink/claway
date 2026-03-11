package service

import (
	"fmt"
	"strings"
)

// computeUnifiedDiff produces a simple unified-diff-style string
// comparing oldContent and newContent line by line.
// Returns an empty string if the contents are identical.
func computeUnifiedDiff(oldContent, newContent string) string {
	oldLines := splitLines(oldContent)
	newLines := splitLines(newContent)

	// Use a simple LCS-based diff to produce unified output.
	lcs := longestCommonSubsequence(oldLines, newLines)

	var hunks []string
	oi, ni, li := 0, 0, 0

	for li < len(lcs) {
		// Advance past lines that don't match the next LCS element.
		for oi < len(oldLines) && oldLines[oi] != lcs[li] {
			hunks = append(hunks, fmt.Sprintf("-%s", oldLines[oi]))
			oi++
		}
		for ni < len(newLines) && newLines[ni] != lcs[li] {
			hunks = append(hunks, fmt.Sprintf("+%s", newLines[ni]))
			ni++
		}
		// Context line (in the LCS).
		hunks = append(hunks, fmt.Sprintf(" %s", lcs[li]))
		oi++
		ni++
		li++
	}

	// Remaining lines after the last LCS element.
	for oi < len(oldLines) {
		hunks = append(hunks, fmt.Sprintf("-%s", oldLines[oi]))
		oi++
	}
	for ni < len(newLines) {
		hunks = append(hunks, fmt.Sprintf("+%s", newLines[ni]))
		ni++
	}

	// If there are no removals or additions, contents are identical.
	hasDiff := false
	for _, h := range hunks {
		if h[0] == '+' || h[0] == '-' {
			hasDiff = true
			break
		}
	}
	if !hasDiff {
		return ""
	}

	return strings.Join(hunks, "\n")
}

// splitLines splits a string into lines, preserving empty trailing lines.
func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

// longestCommonSubsequence returns the LCS of two string slices.
func longestCommonSubsequence(a, b []string) []string {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] >= dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}

	// Backtrack to find the actual LCS.
	result := make([]string, 0, dp[m][n])
	i, j := m, n
	for i > 0 && j > 0 {
		if a[i-1] == b[j-1] {
			result = append(result, a[i-1])
			i--
			j--
		} else if dp[i-1][j] >= dp[i][j-1] {
			i--
		} else {
			j--
		}
	}

	// Reverse the result.
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}

	return result
}

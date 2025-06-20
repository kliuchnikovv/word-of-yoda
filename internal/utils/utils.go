package utils

import "math/bits"

// HasLeadingZeroBits проверяет, что в hash (byte slice) первые difficulty бит равны 0
func HasLeadingZeroBits(hash []byte, difficulty int) bool {
	// Идем по байтам, считаем ведущие нули
	bitsNeeded := difficulty
	for i := 0; i < len(hash) && bitsNeeded > 0; i++ {
		// берем текущий байт
		zeroBits := bits.LeadingZeros8(hash[i]) // число ведущих нулевых бит в этом байте (0..8)
		if zeroBits >= bitsNeeded {
			// все требуемые нули поместились в этот байт
			return true
		}
		if zeroBits < 8 {
			// встретился ненулевой бит до конца байта
			return false
		}
		// весь байт нулевой, уменьшаем требуемые биты
		bitsNeeded -= 8
	}
	// Если дошли до конца hash или bitsNeeded == 0
	return bitsNeeded <= 0
}

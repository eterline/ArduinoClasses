package qweweproto

func SeparatePacketBytes(data []byte) (separated []byte, shift int, allLen int) {
	const headerSize = 3 // 1 byte start + 2 bytes length

	if len(data) < headerSize {
		return nil, 0, 0
	}

	for i := 0; i <= len(data)-headerSize; i++ {
		if data[i] != byte(PKT_QW_START) {
			continue
		}

		if i+2 >= len(data) {
			// Ждём ещё байты для length
			break
		}

		packetLen := int(data[i+1])<<8 | int(data[i+2])
		totalLen := packetLen + 3 // + start byte + 2 length bytes

		if i+totalLen <= len(data) {
			// Есть весь пакет — возвращаем
			return data[i : i+totalLen], i, totalLen - 1
		}

		// Пакет ещё не полностью получен — ждём
		break
	}

	// Ничего не нашли
	return nil, 0, 0
}

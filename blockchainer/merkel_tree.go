package blockchainer

import "crypto/sha256"

//This func calculating result hash on Merkel tree
func merkelMaker(seedling [][32]byte) [32]byte {
	seedlingLen := len(seedling)
	if seedlingLen == 1 {
		return seedling[0]
	}

	var spring [][32]byte
	if seedlingLen%2 == 0 {
		spring = make([][32]byte, seedlingLen/2)
		for i := 0; i < seedlingLen/2; i++ {
			tmp := seedling[i*2+1][:]
			spring[i] = sha256.Sum256(append(seedling[i*2][:], tmp...)[:])
		}
	} else {
		spring = make([][32]byte, seedlingLen/2+1)
		for i := 0; i < seedlingLen/2; i++ {
			tmp := seedling[i*2+1][:]
			spring[i] = sha256.Sum256(append(seedling[i*2][:], tmp...)[:])
		}
		spring[seedlingLen/2] = seedling[seedlingLen-1]
	}

	return merkelMaker(spring)
}

func GetMerkelHash(messages []string) [32]byte {

	//Making hashes array for calculating the tree
	hashes := [][32]byte{}
	for _, msg := range messages {
		hashes = append(hashes, sha256.Sum256([]byte(msg)))
	}

	return merkelMaker(hashes)
}

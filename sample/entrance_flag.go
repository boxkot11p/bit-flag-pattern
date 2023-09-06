package sample

type EntranceFlag int64

const (
	EntranceFlag_UNSPECIFIED EntranceFlag = 0b000
	EntranceFlag_NORMAL      EntranceFlag = 0b001
    EntranceFlag_SPECIAL     EntranceFlag = 0b010
	EntranceFlag_PREMIUM     EntranceFlag = 0b100
)

func HasEntranceFlag(dst int64, src EntranceFlag) bool {
	return dst & int64(src) == int64(src)
}
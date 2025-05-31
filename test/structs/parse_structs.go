package test_parser_structs

type L3 struct {
	Sub_sub_info int64
	Flooaat      float64
	Buul         bool
	Listofnums   []int64
}

type L2 struct {
	Sub_info []L3
}

type L1 struct {
	Name      string
	Info      L2
	Moreinfo  []L2
	Emptyinfo []L2 // Keep this empty
}

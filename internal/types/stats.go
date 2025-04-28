package types

const StorageKeySchemeYieldRate = "scheme_yield_rate."

type SchemeYieldRate struct {
	Scheme Scheme
	Yiled  int
	Total  int
}

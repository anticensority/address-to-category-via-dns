package types

type AddressToIntCat struct {
	Hostnames map[string]int
	Ipv4 map[string]int
	Ipv6 map[string]int
	Ipv4Subnets [][]int
}


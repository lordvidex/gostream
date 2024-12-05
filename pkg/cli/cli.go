package cli

import "github.com/urfave/cli/v3"

// MergeChains combines multiple cli.ValueSourceChains into one
func MergeChains(chains ...cli.ValueSourceChain) cli.ValueSourceChain {
	ch := cli.NewValueSourceChain()

	for _, chain := range chains {
		ch.Chain = append(ch.Chain, chain.Chain...)
	}

	return ch
}

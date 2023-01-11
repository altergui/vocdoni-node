package api

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestAPIHelpers_encodeEVMResultsArgs(t *testing.T) {
	type args struct {
		electionId            common.Hash
		organizationId        common.Address
		censusRoot            common.Hash
		sourceContractAddress common.Address
		results               [][]*big.Int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "encodeEVMResultsArgs0",
			args: args{
				electionId:            common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
				organizationId:        common.HexToAddress("0x0000000000000000000000000000000000000001"),
				censusRoot:            common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
				sourceContractAddress: common.HexToAddress("0x0000000000000000000000000000000000000001"),
				results: [][]*big.Int{
					{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
					{big.NewInt(4), big.NewInt(5), big.NewInt(6)},
				},
			},
			want:    "0x000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000006",
			wantErr: false,
		},
		{
			name: "encodeEVMResultsArgs1",
			args: args{
				electionId:            common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
				organizationId:        common.HexToAddress("0x326C977E6efc84E512bB9C30f76E30c160eD06FB"),
				censusRoot:            common.HexToHash("0xff00000000000000000000000000000000000000000000000000000000000002"),
				sourceContractAddress: common.HexToAddress("0xCC79157eb46F5624204f47AB42b3906cAA40eaB7"),
				results: [][]*big.Int{
					{big.NewInt(0), big.NewInt(0), big.NewInt(0)},
				},
			},
			want:    "0x0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000326c977e6efc84e512bb9c30f76e30c160ed06fbff00000000000000000000000000000000000000000000000000000000000002000000000000000000000000cc79157eb46f5624204f47ab42b3906caa40eab700000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeEVMResultsArgs(tt.args.electionId, tt.args.organizationId, tt.args.censusRoot, tt.args.sourceContractAddress, tt.args.results)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeEVMResultsArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("encodeEVMResultsArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

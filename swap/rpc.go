package swap

import (
	"context"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gopartyparrot/goparrot_buy/config"
)

type TokenAccountInfo struct {
	Mint    solana.PublicKey
	Account solana.PublicKey
}

func GetTokenAccountsBalance(
	ctx context.Context,
	clientRPC rpc.Client,
	accounts ...solana.PublicKey,
) (map[string]uint64, error) {
	res, err := clientRPC.GetMultipleAccounts(ctx, accounts...)
	if err != nil {
		return nil, err
	}
	tokenAccounts := map[string]uint64{}
	for i, a := range res.Value {
		if a.Owner.Equals(solana.TokenProgramID) {
			ta := token.Account{}
			err = bin.NewBinDecoder(a.Data.GetBinary()).Decode(&ta)
			if err != nil {
				return nil, err
			}
			tokenAccounts[accounts[i].String()] = ta.Amount
		} else {
			tokenAccounts[accounts[i].String()] = a.Lamports
		}
	}
	return tokenAccounts, nil
}

func GetTokenAccountsFromMints(
	ctx context.Context,
	clientRPC rpc.Client,
	owner solana.PublicKey,
	mints ...solana.PublicKey,
) (map[string]solana.PublicKey, map[string]solana.PublicKey, error) {

	duplicates := map[string]bool{}
	tokenAccounts := []solana.PublicKey{}
	tokenAccountInfos := []TokenAccountInfo{}
	for _, m := range mints {
		if ok := duplicates[m.String()]; ok {
			continue
		}
		duplicates[m.String()] = true
		a, _, err := solana.FindAssociatedTokenAddress(owner, m)
		if err != nil {
			return nil, nil, err
		}
		// Use owner address for NativeSOL mint
		if m.String() == config.NativeSOL {
			a = owner
		}
		tokenAccounts = append(tokenAccounts, a)
		tokenAccountInfos = append(tokenAccountInfos, TokenAccountInfo{
			Mint:    m,
			Account: a,
		})
	}

	res, err := clientRPC.GetMultipleAccounts(ctx, tokenAccounts...)
	if err != nil {
		return nil, nil, err
	}

	missingAccounts := map[string]solana.PublicKey{}
	existingAccounts := map[string]solana.PublicKey{}
	for i, a := range res.Value {
		tai := tokenAccountInfos[i]
		if a == nil {
			missingAccounts[tai.Mint.String()] = tai.Account
			continue
		}
		if tai.Mint.String() == config.NativeSOL {
			existingAccounts[tai.Mint.String()] = owner
			continue
		}
		var ta token.Account
		err = bin.NewBinDecoder(a.Data.GetBinary()).Decode(&ta)
		if err != nil {
			return nil, nil, err
		}
		existingAccounts[tai.Mint.String()] = tai.Account
	}

	return existingAccounts, missingAccounts, nil
}

func ExecuteInstructions(
	ctx context.Context,
	clientRPC *rpc.Client,
	signers []solana.PrivateKey,
	inst ...solana.Instruction,
) (*solana.Signature, error) {

	recent, err := clientRPC.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}

	tx, err := solana.NewTransaction(
		inst,
		recent.Value.Blockhash,
		solana.TransactionPayer(signers[0].PublicKey()),
	)
	if err != nil {
		return nil, err
	}

	// For debug, print tx accounts
	// tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Instruction"))

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			for _, payer := range signers {
				if payer.PublicKey().Equals(key) {
					return &payer
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	sig, err := clientRPC.SendTransactionWithOpts(
		ctx,
		tx,
		true, // set false to simulate transaction first
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return nil, err
	}

	return &sig, nil
}

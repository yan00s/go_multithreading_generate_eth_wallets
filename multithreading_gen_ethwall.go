package main

import (
	"crypto/ecdsa"
	"crypto/sha512"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/umahmood/mnemonic"
	"golang.org/x/crypto/pbkdf2"
)

func create_mnemonic() string {
	m, err := mnemonic.New(mnemonic.DefaultConfig) // default 128 bits
	if err != nil {
		log.Fatal(err)
	}
	words, err := m.Words()
	if err != nil {
		log.Fatal(err)
	}
	result := strings.Join(words, " ")
	return result
}

func get_result_time(bef_time int64) string {
	result_time := time.Now().UnixMilli()
	result_second := (result_time - bef_time) / 1000
	result_milisec := (result_time - bef_time) % 1000
	result := fmt.Sprintf("result: %v.%v seconds", result_second, result_milisec)
	return result
}

func write_result(result_str string) {
	time_now := time.Now().Format("15-04-05")
	name_file := fmt.Sprintf("result_wallets_%v.txt", time_now)
	os.WriteFile(name_file, []byte(result_str), 0644)
}

func NewSeed(mnemonic, password string) []byte {
	return pbkdf2.Key([]byte(mnemonic), []byte("mnemonic"+password), 2048, 64, sha512.New)
}

func MustParseDerivationPath(path string) (accounts.DerivationPath, error) {
	return accounts.ParseDerivationPath(path)
}

func SeedPathToECDSA(seed []byte, path accounts.DerivationPath) (*ecdsa.PrivateKey, error) {
	key, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	for _, n := range path {
		key, err = key.Derive(n)
		if err != nil {
			return nil, err
		}
	}

	keyEC, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}

	return keyEC.ToECDSA(), nil
}

func MnemonicPathToECDSA(mnemonic, password, pathRaw string) (*ecdsa.PrivateKey, error) {
	seed := NewSeed(mnemonic, password)

	path, err := MustParseDerivationPath(pathRaw)
	if err != nil {
		return nil, err
	}

	return SeedPathToECDSA(seed, path)
}

func gen_wallets(count int, finish_chan chan int, res_chan chan string) {
	for i := count; i > 0; i-- {
		mnemonic := create_mnemonic()
		privateKey, err := MnemonicPathToECDSA(mnemonic, "", "m/44'/60'/0'/0/0")
		if err != nil {
			log.Fatalln(err)
		}
		privateKeyBytes := crypto.FromECDSA(privateKey)
		addr := crypto.PubkeyToAddress(privateKey.PublicKey)
		privateKeyHex := hexutil.Encode(privateKeyBytes)

		result := fmt.Sprintf("Mnemonic: %v\nAddressEth: %v\nPrivateKey: %v\n", mnemonic, addr, privateKeyHex)
		res_chan <- result

	}
	finish_chan <- 1
}

func wait_finish(finish_chan chan int, need_count int) {
	for i := 0; i < need_count; i++ {
		re := <-finish_chan
		_ = re
	}
}

func get_result(res_chan chan string, generated_count int) []string {
	results_wallets := []string{}

	for i := 0; i < generated_count; i++ {
		result := <-res_chan
		results_wallets = append(results_wallets, result)
	}
	return results_wallets
}

func main() {
	thread_count := 1
	count_mnemon := 0

	fmt.Printf("How many wallets to create? ")
	fmt.Scanf("%v", &count_mnemon)
	fmt.Printf("How many threads to run? ")
	fmt.Scanf("%v", &thread_count)

	finish_chan := make(chan int, thread_count)
	res_chan := make(chan string, count_mnemon)

	bef_time := time.Now().UnixMilli()

	count_per_th := count_mnemon / thread_count
	max_count := count_per_th * thread_count
	need_more := count_mnemon - max_count

	for i := 0; i < thread_count; i++ {
		if need_more > 0 {
			go gen_wallets(count_per_th+need_more, finish_chan, res_chan)
			need_more = 0
		} else {
			go gen_wallets(count_per_th, finish_chan, res_chan)
		}
	}

	wait_finish(finish_chan, thread_count)
	results_wallets := get_result(res_chan, count_mnemon)

	fmt.Println(get_result_time(bef_time))
	write_result(strings.Join(results_wallets, "\n"))
	fmt.Println("count:", len(results_wallets))
}

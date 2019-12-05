package bbrpc

const (
	_tPassphrase = "123" //默认的测试密码，rpc unlockkey等用
)

var (
	tCryptonightAddr = AddrKeypair{ //冷钱包挖矿地址
		Keypair: Keypair{Privkey: "eadae10eb384b4d090c10bf2469ee359e32c179026f616ebdf38318ccda5a068", Pubkey: "639ddcfda6e7357cb6543ecb328d6abd130daedaa26beb09bf0e34260f583d77"},
		Address: "1ewyng3s66g7by2fbdehdnbgd2eypn39jscz59dkw6qktdzewknhsmk4t"}
	tCryptonightKey = AddrKeypair{ //挖矿私钥
		Keypair: Keypair{Privkey: "174c4fabefc9573c9cd506dff7f6cb0c54ecaaa63cc0d9f53da7e9c133a01c3a", Pubkey: "ea34707897e6a9ec8e4038179a75fb29d1204a3a5bc4bd07c1c6454d3feac3f7"},
		Address: "1yz1ymftd8q3c21xxrhdkmjh0t4mzpxct2ww413qcn7k9ey3g6knfswbw"}

	//几组地址
	tAddr0 = AddrKeypair{
		Keypair: Keypair{Privkey: "195cd69eff4580ad2430f92d2c86865c596e72edb33f40df5d41c97883241c7c", Pubkey: "a7386f6cbe769fda91462637393970850ae7528d2cee5214c26cc4b27c014a65"},
		Address: "1cn502z5jrhpc452jxrp8tmq71a2q0e9s6wk4d4etkxvbwv3f72ksbkdn"}
	tAddr1 = AddrKeypair{
		Keypair: Keypair{Privkey: "3de774bfb200a46f6d969f5e080572859bc5d7b297fdb34471f55be3326b2153", Pubkey: "1fb8c0c79a506fd8fcca12065331110ae4aedceb2eac38f75379174c6a5b1bff"},
		Address: "1zwdnptjc2xwn7xsrngqeqq5ewg512cak0r9cnz6rdx89nhy0q0fstv2y"}
	tAddr2 = AddrKeypair{
		Keypair: Keypair{Privkey: "8c49b0f3788e07025303ef763e55d14781c09d43cb749628d26280f8f6912336", Pubkey: "5910534ab7629ccb73659df42afc3c382597223a9caa4040a687dbebbe1aa88a"},
		Address: "1ham1nfqbve3tcg20nae3m8mq4mw3sz1ayjepawybkhhbejjk21cvjnx3"}
	tAddr3 = AddrKeypair{
		Keypair: Keypair{Privkey: "5dd0705adf24f1177cedf2795521748358ec08b2d46ddb659f4f68e870433e60", Pubkey: "e4dcb0b8282298a43d5f8c5cbdd3bc27e7f6a44bf0be04e38301655c09038fdb"},
		Address: "1ve7g62awcm0r7rr4qvr4q97pwwkvsmxxbj65yfd4k0h2he5gvkj3d8dz"}
)

package types

type Node struct {
	account    Account
	peers      []Peer
	blockchain BlockChain
	vote       Vote
	privateKey string
}

type Account struct {
	publicKey   string
	permissions []TableAuth
	balance     Balance
	bills       []Bill
}

type Balance struct{}

type Bill struct{}

type Casino struct {
	votes        []Vote
	payments     []Payment
	downpayments []Payment
}

type Payment struct {
	nodeID string
	pay    int
}

type TableAuth struct {
	table      string
	permission string
	info       string
}

type Peer struct {
	account Account
	vote    Vote
}

type BlockChain struct {
	blocks []Block
}

type Transaction struct {
	transType string
	value     string
	signature string
}

type Block struct {
	transactions  []Transaction
	prevPublicKey string
	index         int
}

type Vote struct {
	stake      int
	position   bool
	confidence int
	account    Account
}

type MetaData struct {
	hashes       Account
	currentBlock *Block
}

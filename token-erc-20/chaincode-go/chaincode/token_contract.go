package chaincode

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Define key names for options
const totalSupplyKey = "totalSupply"

// Define objectType names for prefix
const allowancePrefix = "allowance"

// SmartContract provides functions for transferring tokens between accounts
type SmartContract struct {
	contractapi.Contract
}

//定义htlc
type Htlc struct {
	Id      string `json:"Id"`
	Sender  string `json:"sender"`
	Amount  int    `json:"amount"`
	Premage int
	// Get ID of submitting client identity `json:"pre_image"`
	HashValue string `json:"hash_value"`
	State     int    `json:"state"`
}

//define user's wallet
type Asset struct {
	Name string `json:"Name"`
	//Price    string `json:"Price"`
	//Quantity string `json:"Quantity"`
	Money   int    `json:"Money"`
	Address string `json:"Address"` //what's Address's data type? equal Id? int? How to create the Address?
}

type Asset_address struct {
	Name_ad string `json:"Name"`
	Pwd     string `json:"Pwd"`
	//Price_ad    string `json:"Price"`
	//Quantity_ad string `json:"Quantity_ad"`
	//Time_ad     string `json:"Time"`
}

// event provides an organized struct for emitting events
type event struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value int    `json:"value"`
}

func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, name string, pwd string) error {
	//creat address
	asset_address := Asset_address{
		Name_ad: name,
		Pwd:     pwd,
		//Price_ad:    price,
		//Quantity_ad: quantity,
		//Time_ad:     time,
	}
	addresString, err := json.Marshal(asset_address)
	if err != nil {
		return err
	}
	addressBytes := sha256.Sum256(addresString)
	address := hex.EncodeToString(addressBytes[:])

	//verify the correctness of the asset
	exists, err := s.AssetExists(ctx, address)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", address)
	}
	//init asset
	money := 0
	asset := Asset{
		Name:    name,
		Address: address,
		//Price:    price,
		//Quantity: quantity,
		Money: money,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(address, assetJSON)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, address string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(address)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

//ReadAsset returns the asset stored in the world state with given address
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, address string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(address)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", address)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// Mint creates new tokens and adds them to minter's account balance
// This function triggers a Transfer event
func (s *SmartContract) Mint(ctx contractapi.TransactionContextInterface, amount int, address string) error {
	//at begin,verify amout of correction
	if amount <= 0 {
		return fmt.Errorf("mint amount must be a positive integer")
	}

	//get user's asset by address
	asset, err := s.ReadAsset(ctx, address) //auto type
	//mint token
	asset.Money += amount
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(address, assetJSON)
	if err != nil {
		return nil
	}
	return nil
}

// Check minter authorization - this sample assumes Org1 is the central banker with privilege to mint new tokens
// 	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
// 	if err != nil {
// 		return fmt.Errorf("failed to get MSPID: %v", err)
// 	}
// 	if clientMSPID != "Org1MSP" {
// 		return fmt.Errorf("client is not authorized to mint new tokens")
// 	}

// 	// Get ID of submitting client identity
// 	minter, err := ctx.GetClientIdentity().GetID()
// 	if err != nil {
// 		return fmt.Errorf("failed to get client id: %v", err)
// 	}

// 	if amount <= 0 {
// 		return fmt.Errorf("mint amount must be a positive integer")
// 	}

// 	currentBalanceBytes, err := ctx.GetStub().GetState(minter)
// 	if err != nil {
// 		return fmt.Errorf("failed to read minter account %s from world state: %v", minter, err)
// 	}

// 	var currentBalance int

// 	// If minter current balance doesn't yet exist, we'll create it with a current balance of 0
// 	if currentBalanceBytes == nil {
// 		currentBalance = 0
// 	} else {
// 		currentBalance, _ = strconv.Atoi(string(currentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.
// 	}

// 	updatedBalance := currentBalance + amount

// 	err = ctx.GetStub().PutState(minter, []byte(strconv.Itoa(updatedBalance)))
// 	if err != nil {
// 		return err
// 	}

// 	// Update the totalSupply
// 	totalSupplyBytes, err := ctx.GetStub().GetState(totalSupplyKey)
// 	if err != nil {
// 		return fmt.Errorf("failed to retrieve total token supply: %v", err)
// 	}

// 	var totalSupply int

// 	// If no tokens have been minted, initialize the totalSupply
// 	if totalSupplyBytes == nil {
// 		totalSupply = 0
// 	} else {
// 		totalSupply, _ = strconv.Atoi(string(totalSupplyBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.
// 	}

// 	// Add the mint amount to the total supply and update the state
// 	totalSupply += amount
// 	err = ctx.GetStub().PutState(totalSupplyKey, []byte(strconv.Itoa(totalSupply)))
// 	if err != nil {
// 		return err
// 	}

// 	// Emit the Transfer event
// 	transferEvent := event{"0x0", minter, amount}
// 	transferEventJSON, err := json.Marshal(transferEvent)
// 	if err != nil {
// 		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
// 	}
// 	err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
// 	if err != nil {
// 		return fmt.Errorf("failed to set event: %v", err)
// 	}

// 	log.Printf("minter account %s balance updated from %d to %d", minter, currentBalance, updatedBalance)

// 	return nil
// }

// // Burn redeems tokens the minter's account balance
// // This function triggers a Transfer event
func (s *SmartContract) Burn(ctx contractapi.TransactionContextInterface, amount int) error {

	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to burn new tokens
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != "Org1MSP" {
		return fmt.Errorf("client is not authorized to mint new tokens")
	}

	// Get ID of submitting client identity
	minter, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	if amount <= 0 {
		return errors.New("burn amount must be a positive integer")
	}

	currentBalanceBytes, err := ctx.GetStub().GetState(minter)
	if err != nil {
		return fmt.Errorf("failed to read minter account %s from world state: %v", minter, err)
	}

	var currentBalance int

	// Check if minter current balance exists
	if currentBalanceBytes == nil {
		return errors.New("The balance does not exist")
	}

	currentBalance, _ = strconv.Atoi(string(currentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.

	updatedBalance := currentBalance - amount

	err = ctx.GetStub().PutState(minter, []byte(strconv.Itoa(updatedBalance)))
	if err != nil {
		return err
	}

	// Update the totalSupply
	totalSupplyBytes, err := ctx.GetStub().GetState(totalSupplyKey)
	if err != nil {
		return fmt.Errorf("failed to retrieve total token supply: %v", err)
	}

	// If no tokens have been minted, throw error
	if totalSupplyBytes == nil {
		return errors.New("totalSupply does not exist")
	}

	totalSupply, _ := strconv.Atoi(string(totalSupplyBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.

	// Subtract the burn amount to the total supply and update the state
	totalSupply -= amount
	err = ctx.GetStub().PutState(totalSupplyKey, []byte(strconv.Itoa(totalSupply)))
	if err != nil {
		return err
	}

	// Emit the Transfer event
	transferEvent := event{minter, "0x0", amount}
	transferEventJSON, err := json.Marshal(transferEvent)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	log.Printf("minter account %s balance updated from %d to %d", minter, currentBalance, updatedBalance)

	return nil
}

// Transfer transfers tokens from client account to recipient account
// // recipient account must be a valid clientID as returned by the ClientID() function
// // This function triggers a Transfer event
func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, sender string, recipient string, amount int) error {

	// Get ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	err = s.transferHelper(ctx, sender, recipient, amount)
	if err != nil {
		return fmt.Errorf("failed to transfer: %v", err)
	}

	// Emit the Transfer event
	transferEvent := event{clientID, recipient, amount}
	transferEventJSON, err := json.Marshal(transferEvent)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	return nil
}

// BalanceOf returns the balance of the given account
func (s *SmartContract) BalanceOf(ctx contractapi.TransactionContextInterface, account string) (int, error) {
	balanceBytes, err := ctx.GetStub().GetState(account)
	if err != nil {
		return 0, fmt.Errorf("failed to read from world state: %v", err)
	}
	if balanceBytes == nil {
		return 0, fmt.Errorf("the account %s does not exist", account)
	}

	balance, _ := strconv.Atoi(string(balanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.

	return balance, nil
}

// ClientAccountBalance returns the balance of the requesting client's account
func (s *SmartContract) ClientAccountBalance(ctx contractapi.TransactionContextInterface) (int, error) {

	// Get ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return 0, fmt.Errorf("failed to get client id: %v", err)
	}

	balanceBytes, err := ctx.GetStub().GetState(clientID)
	if err != nil {
		return 0, fmt.Errorf("failed to read from world state: %v", err)
	}
	if balanceBytes == nil {
		return 0, fmt.Errorf("the account %s does not exist", clientID)
	}

	balance, _ := strconv.Atoi(string(balanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.

	return balance, nil
}

// ClientAccountID returns the id of the requesting client's account
// In this implementation, the client account ID is the clientId itself
// Users can use this function to get their own account id, which they can then give to others as the payment address
func (s *SmartContract) ClientAccountID(ctx contractapi.TransactionContextInterface) (string, error) {

	// Get ID of submitting client identity
	clientAccountID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to get client id: %v", err)
	}

	return clientAccountID, nil
}

// TotalSupply returns the total token supply
func (s *SmartContract) TotalSupply(ctx contractapi.TransactionContextInterface) (int, error) {

	// Retrieve total supply of tokens from state of smart contract
	totalSupplyBytes, err := ctx.GetStub().GetState(totalSupplyKey)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve total token supply: %v", err)
	}

	var totalSupply int

	// If no tokens have been minted, return 0
	if totalSupplyBytes == nil {
		totalSupply = 0
	} else {
		totalSupply, _ = strconv.Atoi(string(totalSupplyBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.
	}

	log.Printf("TotalSupply: %d tokens", totalSupply)

	return totalSupply, nil
}

// Approve allows the spender to withdraw from the calling client's token account
// The spender can withdraw multiple times if necessary, up to the value amount
// This function triggers an Approval event
func (s *SmartContract) Approve(ctx contractapi.TransactionContextInterface, spender string, value int) error {

	// Get ID of submitting client identity
	owner, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	// Create allowanceKey
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{owner, spender})
	if err != nil {
		return fmt.Errorf("failed to create the composite key for prefix %s: %v", allowancePrefix, err)
	}

	// Update the state of the smart contract by adding the allowanceKey and value
	err = ctx.GetStub().PutState(allowanceKey, []byte(strconv.Itoa(value)))
	if err != nil {
		return fmt.Errorf("failed to update state of smart contract for key %s: %v", allowanceKey, err)
	}

	// Emit the Approval event
	approvalEvent := event{owner, spender, value}
	approvalEventJSON, err := json.Marshal(approvalEvent)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Approval", approvalEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	log.Printf("client %s approved a withdrawal allowance of %d for spender %s", owner, value, spender)

	return nil
}

// Allowance returns the amount still available for the spender to withdraw from the owner
func (s *SmartContract) Allowance(ctx contractapi.TransactionContextInterface, owner string, spender string) (int, error) {

	// Create allowanceKey
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{owner, spender})
	if err != nil {
		return 0, fmt.Errorf("failed to create the composite key for prefix %s: %v", allowancePrefix, err)
	}

	// Read the allowance amount from the world state
	allowanceBytes, err := ctx.GetStub().GetState(allowanceKey)
	if err != nil {
		return 0, fmt.Errorf("failed to read allowance for %s from world state: %v", allowanceKey, err)
	}

	var allowance int

	// If no current allowance, set allowance to 0
	if allowanceBytes == nil {
		allowance = 0
	} else {
		allowance, err = strconv.Atoi(string(allowanceBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.
	}

	log.Printf("The allowance left for spender %s to withdraw from owner %s: %d", spender, owner, allowance)

	return allowance, nil
}

// TransferFrom transfers the value amount from the "from" address to the "to" address
// This function triggers a Transfer event
func (s *SmartContract) TransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, value int) error {

	// // Get ID of submitting client identity
	// spender, err := ctx.GetClientIdentity().GetID()
	// if err != nil {
	// 	return fmt.Errorf("failed to get client id: %v", err)
	// }

	// // Create allowanceKey
	// allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{from, spender})
	// if err != nil {
	// 	return fmt.Errorf("failed to create the composite key for prefix %s: %v", allowancePrefix, err)
	// }

	// // Retrieve the allowance of the spender
	// currentAllowanceBytes, err := ctx.GetStub().GetState(allowanceKey)
	// if err != nil {
	// 	return fmt.Errorf("failed to retrieve the allowance for %s from world state: %v", allowanceKey, err)
	// }

	// var currentAllowance int
	// currentAllowance, _ = strconv.Atoi(string(currentAllowanceBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.

	// // Check if transferred value is less than allowance
	// if currentAllowance < value {
	// 	return fmt.Errorf("spender does not have enough allowance for transfer")
	// }

	// // Initiate the transfer
	// err = s.transferHelper(ctx, from, to, value)
	// if err != nil {
	// 	return fmt.Errorf("failed to transfer: %v", err)
	// }

	// // Decrease the allowance
	// updatedAllowance := currentAllowance - value
	// err = ctx.GetStub().PutState(allowanceKey, []byte(strconv.Itoa(updatedAllowance)))
	// if err != nil {
	// 	return err
	// }

	// Emit the Transfer event
	transferEvent := event{from, to, value}
	transferEventJSON, err := json.Marshal(transferEvent)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	//log.Printf("spender %s allowance updated from %d to %d", spender, currentAllowance, updatedAllowance)

	return nil
}

// Helper Functions

// transferHelper is a helper function that transfers tokens from the "from" address to the "to" address
// Dependant functions include Transfer and TransferFrom
func (s *SmartContract) transferHelper(ctx contractapi.TransactionContextInterface, from string, to string, value int) error {

	if from == to {
		return fmt.Errorf("cannot transfer to and from same client account")
	}

	if value < 0 { // transfer of 0 is allowed in ERC-20, so just validate against negative amounts
		return fmt.Errorf("transfer amount cannot be negative")
	}

	//get sender assets
	senderAsset, err := s.ReadAsset(ctx, from)

	//get reveive assets
	receiveAsset, err := s.ReadAsset(ctx, to)
	if err != nil {
		return fmt.Errorf("failed to read receive asset %s from world state: %v", to, err)
	}

	//transaction
	senderAsset.Money -= value
	receiveAsset.Money += value

	//update
	assetJSON, err := json.Marshal(senderAsset)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(from, assetJSON)
	if err != nil {
		return nil
	}

	assetJSON2, err := json.Marshal(receiveAsset)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(to, assetJSON2)
	if err != nil {
		return nil
	}
	return nil
	// fromCurrentBalanceBytes, err := ctx.GetStub().GetState(from)
	// if err != nil {
	// 	return fmt.Errorf("failed to read client account %s from world state: %v", from, err)
	// }

	// if fromCurrentBalanceBytes == nil {
	// 	return fmt.Errorf("client account %s has no balance", from)
	// }

	// fromCurrentBalance, _ := strconv.Atoi(string(fromCurrentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.

	// if fromCurrentBalance < value {
	// 	return fmt.Errorf("client account %s has insufficient funds", from)
	// }

	// toCurrentBalanceBytes, err := ctx.GetStub().GetState(to)
	// if err != nil {
	// 	return fmt.Errorf("failed to read recipient account %s from world state: %v", to, err)
	// }

	// var toCurrentBalance int
	// // If recipient current balance doesn't yet exist, we'll create it with a current balance of 0
	// if toCurrentBalanceBytes == nil {
	// 	toCurrentBalance = 0
	// } else {
	// 	toCurrentBalance, _ = strconv.Atoi(string(toCurrentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.
	// }

	// fromUpdatedBalance := fromCurrentBalance - value
	// toUpdatedBalance := toCurrentBalance + value

	// err = ctx.GetStub().PutState(from, []byte(strconv.Itoa(fromUpdatedBalance)))
	// if err != nil {
	// 	return err
	// }

	// err = ctx.GetStub().PutState(to, []byte(strconv.Itoa(toUpdatedBalance)))
	// if err != nil {
	// 	return err
	// }

	// log.Printf("client %s balance updated from %d to %d", from, fromCurrentBalance, fromUpdatedBalance)
	// log.Printf("recipient %s balance updated from %d to %d", to, toCurrentBalance, toUpdatedBalance)

	// return nil
}

//根据交易id查询htlc是否存在
func (s *SmartContract) HtlcExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	htlcJSON, err := ctx.GetStub().GetState(id)
	if err != nil { //如果发生错误就显示错误
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return htlcJSON != nil, nil
}

//hash锁函数 卖家拥有token 返回一个交易id放入到账本中去 买家通过id可以查的htmc的信息
//交易id为买家的fabric的地址 这样id就不用传输 error 这里不能将买方地址直接作为交易id 因为买方不止这一次的交易 而地址只有一个
func (s *SmartContract) CreateHash(ctx contractapi.TransactionContextInterface, id string, amount int, premage int, address string) error {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	exists, err := s.HtlcExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}
	if err != nil {
		return fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != "Org1MSP" {
		return fmt.Errorf("client is not authorized to mint new tokens")
	}
	//卖方获取自己的交易地址
	//sender, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get senderclient id: %v", err)
	}
	//确保交易金额的正确
	if amount <= 0 {
		return fmt.Errorf("sender amount must be a positive integer")
	}

	//get asset
	asset, err := s.ReadAsset(ctx, address) //auto type

	//对原像hash加密
	hashByte := sha256.Sum256([]byte(strconv.Itoa(premage)))

	hashCode := hex.EncodeToString(hashByte[:])
	state := 0
	//将初始化htlc hashvalue为idByte；
	htlc := Htlc{
		Id:        id,
		Sender:    asset.Address,
		Amount:    amount,
		Premage:   0, // 先设置为空，在卖方在以太坊上面取得资产后 再由卖方自己更新premage 但是这样的话有个问题 如果卖方不更新呢 这里应该是买方的函数里面设置的
		HashValue: hashCode,
		State:     state,
	}
	//将htll转为字符串的格式
	htlcByte, err := json.Marshal(htlc)
	if err != nil {
		return err
	}

	//解决双花的问题 如何对金额设置标志位 并且这个金额是怎么处理的 怎么处理的 怎么处理的 怎么处理的啊
	//两种方法：第一种直接销毁 然后在退回 但是在以太坊上面怎么弄还不知道 但是销毁的话还要将这个请求发送给mint账户 因为当前这个账户是不能销毁的 和mint的交互是怎么的过程也不知道 做这些干嘛啊 什么都不知道

	//第二种方法 创建一个中间账户
	//问题是这个账户是谁创建的 什么时候创建的 在那个人的基础上能不能进一步的修改 我们设定中间账户的原因是解决双花问题 但是好像那片论文里面没有提到双花问题
	//而哈希时间锁最主要的问题是 资金流的问题 被冻结的钱到哪里去了 必须有个归宿 暂时没有想到自己的方案

	//验证交易id是否存在 存在的就返回错误 然后继续创建 改变premage的值 creathash
	//无错误就将交易id放进账本数据库中

	return ctx.GetStub().PutState(id, htlcByte)
}

//根据交易id查询htlc的内容
func (s *SmartContract) QueryTransId(ctx contractapi.TransactionContextInterface, id string) (*Htlc, error) {
	htlcJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if htlcJSON == nil {
		return nil, fmt.Errorf("the htlc %s does not exist", id)
	}
	//定义变量对象
	var htlc Htlc
	//反序列化
	err = json.Unmarshal(htlcJSON, &htlc)
	if err != nil {
		return nil, err
	}

	return &htlc, nil
}

//  //验证买方hash的作用函数
//  func(h *Htlc) Verify(pwd string){
// 	pwdByte :=sha256.Sum256(pwd)
// 	if pwdByte ==
//  }

//资产转移函数 这里的情景为卖方已经在以太坊上面取得了资产 并且将哈希原像返回给了买方 买房通过哈希原值在fabric上面取得卖方的资产
func (s *SmartContract) AcrossTransfer(ctx contractapi.TransactionContextInterface, pwd int, id string, address string) error {
	//get bao(buy) address
	asset1, err := s.ReadAsset(ctx, address)

	//买方获取自己的交易地址
	receipt, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get receipetclient id: %v", err)
	}
	//call function 根据交易id查询相应的htlc
	//var htlc Htlc
	htlc, err := s.QueryTransId(ctx, id)
	if err != nil {
		return err
	}
	// htlcJSON,err := ctx.GetStub().GetState(id)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to read from world state: %v", err)
	// }
	// if htlcJSON == nil {
	// 	return nil, fmt.Errorf("the htlc %s does not exist", id)
	// }
	// //定义变量对象

	// //反序列化 得到htlc的信息
	// htlc = json.Unmarshal(htlcJSON, &htlc)
	//验证哈希是否正确

	pwdByte := sha256.Sum256([]byte(strconv.Itoa(pwd)))
	hashCode := hex.EncodeToString(pwdByte[:])
	if hashCode != htlc.HashValue {
		return fmt.Errorf("the hash does not match!") //if you change code when the produce is running, you must invoke it again!
	}
	// //验证htlc中的状态是否已经锁定
	if htlc.State != 0 {
		return fmt.Errorf("this htlc transaction state is error")
	}

	//匹配成功 转移资产
	err = s.transferHelper(ctx, htlc.Sender, asset1.Address, htlc.Amount)
	if err != nil {
		return fmt.Errorf("failed to transfer: %v", err)
	}
	//wirte state
	htlc.State = 1
	htlcByte, err := json.Marshal(htlc)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState(id, htlcByte)

	// Emit the Transfer event
	transferEvent := event{htlc.Sender, receipt, htlc.Amount}
	transferEventJSON, err := json.Marshal(transferEvent)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	return nil
}

// func (s *SmartContract) AcrossTransfer(ctx contractapi.TransactionContextInterface, receipt string, pwd string, id string) error {
// 	//买方获取自己的交易地址

// 	//var htlc Htlc
// 	htlc, err := s.QueryTransId(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	//验证哈希是否正确
// 	pwdByte := sha256.Sum256([]byte(pwd))
// 	if pwdByte != htlc.HashValue {
// 		return fmt.Errorf("the hash does not match")
// 	}
// 	// // //验证htlc中的状态是否已经锁定
// 	//  if htlc.State != "HashLOCK" {
// 	//  	return fmt.Errorf("this htlc transaction state is error")
// 	//  }

// 	//匹配成功 转移资产
// 	err = transferHelper(ctx, htlc.Sender, receipt, htlc.Amount)
// 	if err != nil {
// 		return fmt.Errorf("failed to transfer: %v", err)
// 	}

// 	// Emit the Transfer event
// 	transferEvent := event{htlc.Sender, receipt, htlc.Amount}
// 	transferEventJSON, err := json.Marshal(transferEvent)
// 	if err != nil {
// 		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
// 	}
// 	err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
// 	if err != nil {
// 		return fmt.Errorf("failed to set event: %v", err)
// 	}
// 	htlc.State = "Received"
// 	return nil
// }

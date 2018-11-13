/*
* Copyright 2018 The OpenWallet Authors
* This file is part of the OpenWallet library.
*
* The OpenWallet library is free software: you can redistribute it and/or modify
* it under the terms of the GNU Lesser General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* The OpenWallet library is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
* GNU Lesser General Public License for more details.
 */

package tron

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/blocktree/go-owcdrivers/addressEncoder"
	"github.com/blocktree/go-owcrypt"
	"github.com/shopspring/decimal"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/tronprotocol/grpc-gateway/core"
)

func getTxHash(tx *core.Transaction) (txHash []byte, err error) {

	txRaw, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		return nil, err
	}
	txHash = owcrypt.Hash(txRaw, 0, owcrypt.HASH_ALG_SHA256)
	return txHash, err
}

// CreateTransactionRef Done!
// Function: Create a transaction
//
// Java Reference:
// public static Transaction setReference(Transaction transaction, Block newestBlock) {
// 	long blockHeight = newestBlock.getBlockHeader().getRawData().getNumber();
// 	byte[] blockHash = getBlockHash(newestBlock).getBytes();
// 	byte[] refBlockNum = ByteArray.fromLong(blockHeight);
// 	Transaction.raw rawData = transaction.getRawData().toBuilder()
// 	    .setRefBlockHash(ByteString.copyFrom(ByteArray.subArray(blockHash, 8, 16)))
// 	    .setRefBlockBytes(ByteString.copyFrom(ByteArray.subArray(refBlockNum, 6, 8)))
// 	    .build();
// 	return transaction.toBuilder().setRawData(rawData).build();
// }
//
// public static Transaction createTransaction(byte[] from, byte[] to, long amount) {
// 	Transaction.Builder transactionBuilder = Transaction.newBuilder();
// 	Block newestBlock = WalletClient.getBlock(-1);
//
// 	Transaction.Contract.Builder contractBuilder = Transaction.Contract.newBuilder();
// 	Contract.TransferContract.Builder transferContractBuilder = Contract.TransferContract.newBuilder();
//
// 	transferContractBuilder.setAmount(amount);
// 	ByteString bsTo = ByteString.copyFrom(to);
// 	ByteString bsOwner = ByteString.copyFrom(from);
// 	transferContractBuilder.setToAddress(bsTo);
// 	transferContractBuilder.setOwnerAddress(bsOwner);
//
// 	try {
// 	  Any any = Any.pack(transferContractBuilder.build());
// 	  contractBuilder.setParameter(any);
// 	} catch (Exception e) {
// 	  return null;
// 	}
// 	contractBuilder.setType(Transaction.Contract.ContractType.TransferContract);
//
// 	transactionBuilder.getRawDataBuilder().addContract(contractBuilder)
// 	    .setTimestamp(System.currentTimeMillis())//timestamp should be in millisecond format
// 	    .setExpiration(newestBlock.getBlockHeader().getRawData().getTimestamp() + 10 * 60 * 60 * 1000);//exchange can set Expiration by needs
// 	Transaction transaction = transactionBuilder.build();
// 	Transaction refTransaction = setReference(transaction, newestBlock);
// 	return refTransaction;
// }
func (wm *WalletManager) CreateTransactionRef(to_address, owner_address string, amount int64) (txRawHex string, err error) {

	// addressEncoder.AddressDecode return 20 bytes of the center of Address
	to_address_bytes, err := addressEncoder.AddressDecode(to_address, addressEncoder.TRON_mainnetAddress)
	if err != nil {
		return "", err
	} else {
		to_address_bytes = append([]byte{0x41}, to_address_bytes...)
	}

	owner_address_bytes, err := addressEncoder.AddressDecode(owner_address, addressEncoder.TRON_mainnetAddress)
	if err != nil {
		return "", err
	} else {
		owner_address_bytes = append([]byte{0x41}, owner_address_bytes...)
	}

	// --------------------- Generate TX Contract ------------------------

	tc := &core.TransferContract{
		OwnerAddress: owner_address_bytes,
		ToAddress:    to_address_bytes,
		Amount:       amount * 1000000,
	}

	tcRaw, err := proto.Marshal(tc)
	if err != nil {
		return "", err
	}

	txContact := &core.Transaction_Contract{
		Type:         core.Transaction_Contract_TransferContract,
		Parameter:    &any.Any{Value: tcRaw, TypeUrl: "type.googleapis.com/protocol.TransferContract"},
		Provider:     nil,
		ContractName: nil,
	}

	// ----------------------- Get Reference Block ----------------------
	block, err := wm.GetNowBlock()
	if err != nil {
		return "", err
	}
	blockID, err := hex.DecodeString(block.GetBlockHashID())
	if err != nil {
		return txRawHex, err
	}
	refBlockBytes, refBlockHash := blockID[6:8], blockID[8:16]

	// -------------------- Set timestamp --------------------
	/*
		RFC 3339 date strings
		Example 4: Compute Timestamp from Java `System.currentTimeMillis()`.

		    long millis = System.currentTimeMillis();

		    Timestamp timestamp = Timestamp.newBuilder().setSeconds(millis / 1000)
			  .setNanos((int) ((millis % 1000) * 1000000)).build();
	*/
	timestamp := time.Now().UnixNano() / 1000000 // <int64

	// -------------------- Create Traction --------------------
	txRaw := &core.TransactionRaw{
		RefBlockBytes: refBlockBytes,
		RefBlockHash:  refBlockHash,
		Contract:      []*core.Transaction_Contract{txContact},
		Expiration:    timestamp + 10*60*60*60,
		// Timestamp:     timestamp,
	}
	tx := &core.Transaction{
		RawData: txRaw,
		// Signature: nil,
		// Ret:       nil,
	}

	// -------------------- TX Encoding --------------------
	if x, err := proto.Marshal(tx); err != nil {
		return "", err
	} else {
		txRawHex = hex.EncodeToString(x)
	}

	return txRawHex, nil
}

// Done!
// public static Transaction sign(Transaction transaction, ECKey myKey) {
// 	Transaction.Builder transactionBuilderSigned = transaction.toBuilder();
//
// 	byte[] hash = Sha256Hash.hash(transaction.getRawData().toByteArray());
//
// 	List<Contract> listContract = transaction.getRawData().getContractList();
//
// 	for (int i = 0; i < listContract.size(); i++) {
// 	  ECDSASignature signature = myKey.sign(hash);
// 	  ByteString bsSign = ByteString.copyFrom(signature.toByteArray());
// 	  transactionBuilderSigned.addSignature(bsSign);  //Each contract may be signed with a different private key in the future.
// 	}
//
// 	transaction = transactionBuilderSigned.build();
//
// 	return transaction;
// }
func (wm *WalletManager) SignTransactionRef(txRawhex string, privateKey string) (signedTxRaw string, err error) {

	tx := &core.Transaction{}
	if txRawBts, err := hex.DecodeString(txRawhex); err != nil {
		return signedTxRaw, err
	} else {
		if err := proto.Unmarshal(txRawBts, tx); err != nil {
			return signedTxRaw, err
		}
	}
	// tx.GetRawData().GetRefBlockBytes()

	txHash, err := getTxHash(tx)
	if err != nil {
		return signedTxRaw, err
	}
	fmt.Println("Tx Hash = ", hex.EncodeToString(txHash))

	pk, err := hex.DecodeString(privateKey)
	if err != nil {
		return signedTxRaw, err
	}

	for i, _ := range tx.GetRawData().GetContract() {

		// sign, ret := owcrypt.Signature(privateKey, nil, 0, txHash, uint16(len(txHash)), owcrypt.ECC_CURVE_SECP256K1)
		sign, ret := owcrypt.Tron_signature(pk, txHash)
		if ret != owcrypt.SUCCESS {
			err := errors.New(fmt.Sprintf("Signature[%d] faild!", i))
			log.Println(err)
			return signedTxRaw, err
		}
		tx.Signature = append(tx.Signature, sign)
	}

	if x, err := proto.Marshal(tx); err != nil {
		return signedTxRaw, err
	} else {
		signedTxRaw = hex.EncodeToString(x)
	}

	return signedTxRaw, nil

}

// Done!
//   /*
//    * 1. check hash
//    * 2. check double spent
//    * 3. check sign
//    * 4. check balance
//    */
//  public static boolean validTransaction(Transaction signedTransaction) {
// 	assert (signedTransaction.getSignatureCount() == signedTransaction.getRawData().getContractCount());
//
// 	List<Transaction.Contract> listContract = signedTransaction.getRawData().getContractList();
//
// 	byte[] hash = Sha256Hash.hash(signedTransaction.getRawData().toByteArray());
//
// 	int count = signedTransaction.getSignatureCount();
// 	if (count == 0) {
// 	  return false;
// 	}
//
// 	for (int i = 0; i < count; ++i) {
// 	  try {
// 	    Transaction.Contract contract = listContract.get(i);
//
// 	    byte[] owner = getOwner(contract);
// 	    byte[] address = ECKey.signatureToAddress(hash, getBase64FromByteString(signedTransaction.getSignature(i)));
// 	    if (!Arrays.equals(owner, address)) {
// 		return false;
// 	    }
//
// 	  } catch (SignatureException e) {
// 	    e.printStackTrace();
// 	    return false;
// 	  }
// 	}
// 	return true;
//  }
func (wm *WalletManager) ValidSignedTransactionRef(txHex string) error {

	tx := &core.Transaction{}
	if txBytes, err := hex.DecodeString(txHex); err != nil {
		return nil
	} else {
		if err := proto.Unmarshal(txBytes, tx); err != nil {
			return err
		}
	}

	if len(tx.GetSignature()) != len(tx.GetRawData().GetContract()) {
		err := errors.New("ValidSignedTransactionRef faild: no signature found!")
		log.Println(err)
		return err
	}

	listContracts := tx.RawData.GetContract()
	countSignature := len(tx.Signature)

	txHash, err := getTxHash(tx)
	if err != nil {
		return err
	}

	if countSignature == 0 {
		return errors.New("No signature found!")
	}

	for i := 0; i < countSignature; i++ {
		contract := listContracts[i]

		// Get the instance of TransferContract to get Owner Address for validate signature
		tc := &core.TransferContract{}
		if err := proto.Unmarshal(contract.Parameter.GetValue(), tc); err != nil {
			return err
		}

		owner_address_hex := hex.EncodeToString(tc.GetOwnerAddress())

		// pkBytes, ret := owcrypt.RecoverPubkey(tx.Signature[i], txHash, owcrypt.ECC_CURVE_SECP256K1|owcrypt.HASH_OUTSIDE_FLAG)
		pkBytes, ret := owcrypt.RecoverPubkey(tx.Signature[i], txHash, owcrypt.ECC_CURVE_SECP256K1)
		if ret != owcrypt.SUCCESS {
			err := errors.New("ValidSignedTransactionRef faild: owcryt.RecoverPubkey return err!")
			log.Println(err)
			return err
		}

		pkgen_address_bytes, err := createAddressByPkRef(pkBytes)
		if err != nil {
			log.Println(err)
			return err
		}
		pkgen_address_hex := hex.EncodeToString(pkgen_address_bytes[:len(pkgen_address_bytes)-4])

		// Check whether the address is equal between signature generating and contract owner pointed
		if pkgen_address_hex != owner_address_hex {
			return errors.New("Validate failed, signed address is not the owner address!")
		}
	}

	return nil
}

//SendTransaction 发送交易
func (wm *WalletManager) SendTransaction(walletID, to string, amount decimal.Decimal, password string, feesInSender bool) ([]string, error) {
	return nil, nil
}

// ------------------------------------------------------------------------------------------------------
func debugPrintTx(txRawhex string) {

	tx := &core.Transaction{}
	if txRawBts, err := hex.DecodeString(txRawhex); err != nil {
		fmt.Println(err)
	} else {
		if err := proto.Unmarshal(txRawBts, tx); err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv Print Test vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv")

	txHash, err := getTxHash(tx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Tx Hash = ", hex.EncodeToString(txHash))

	txRawD := tx.RawData
	txC := txRawD.GetContract()
	fmt.Println("txRawD.Contract = ")
	for _, c := range txC {
		fmt.Println("\tc.ContractName=", c.ContractName)
		fmt.Println("\tc.Provider   =", c.Provider)
		fmt.Println("\tc.Type       =", c.Type)
		fmt.Println("\tc.Parameter   =", c.Parameter)

		ts := &core.TransferContract{}
		proto.Unmarshal(c.Parameter.Value, ts)
		fmt.Println("\tts.OwnerAddress =", hex.EncodeToString(ts.OwnerAddress))
		fmt.Println("\tts.ToAddress =", hex.EncodeToString(ts.ToAddress))
		fmt.Println("\tts.Amount =", ts.Amount)
	}
	fmt.Println("txRawD.Data =  ", txRawD.Data)
	fmt.Println("txRawD.Auths =   ", txRawD.Auths)
	fmt.Println("txRawD.Scripts =   ", txRawD.Scripts)
	fmt.Println("txRawD.RefBlockBytes = ", hex.EncodeToString(txRawD.RefBlockBytes))
	fmt.Println("txRawD.RefBlockHash Bts = ", txRawD.RefBlockHash, "Len:", len(txRawD.RefBlockHash))
	fmt.Println("txRawD.RefBlockHash Hex = ", hex.EncodeToString(txRawD.RefBlockHash), "Len:", len(hex.EncodeToString(txRawD.RefBlockHash)))
	// dst := make([]byte, 32)
	// bs, err := base64.StdEncoding.Decode(dst, txRawD.RefBlockHash)
	// fmt.Println("txRawD.RefBlockHash base64Bytes = ", bs, "XX = ", dst)

	fmt.Println("")

	fmt.Println("txRawD.RefBlockNum =  ", txRawD.RefBlockNum)
	fmt.Println("txRawD.Expiration =  ", txRawD.Expiration)
	fmt.Println("txRawD.Timestamp =   ", txRawD.Timestamp)
	fmt.Println("tx.Signature[0]     = ", hex.EncodeToString(tx.Signature[0]))
	fmt.Println("tx.Ret          =     ", tx.Ret)

	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ Print Test ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ End")
}

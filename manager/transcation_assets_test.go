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

package manager

import (
	"github.com/blocktree/OpenWallet/log"
	"github.com/blocktree/OpenWallet/openwallet"
	"testing"
)

func testCreateTransactionStep(walletID, accountID, to, amount, feeRate string) (*openwallet.RawTransaction, error) {

	err := tm.RefreshAssetsAccountBalance(testApp, accountID)
	if err != nil {
		log.Error("RefreshAssetsAccountBalance failed, unexpected error:", err)
		return nil, err
	}

	rawTx, err := tm.CreateTransaction(testApp, walletID, accountID, amount, to, feeRate, "")

	if err != nil {
		log.Error("CreateTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTx, nil
}

func testSignTransactionStep(rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	_, err := tm.SignTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, "12345678", rawTx)
	if err != nil {
		log.Error("SignTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Info("rawTx:", rawTx)
	return rawTx, nil
}

func testVerifyTransactionStep(rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	//log.Info("rawTx.Signatures:", rawTx.Signatures)

	_, err := tm.VerifyTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("VerifyTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Info("rawTx:", rawTx)
	return rawTx, nil
}

func testSubmitTransactionStep(rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	tx, err := tm.SubmitTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("SubmitTransaction failed, unexpected error:", err)
		return nil, err
	}
	log.Info("wxID:", tx.WxID)
	log.Info("txID:", rawTx.TxID)

	return rawTx, nil
}

func TestTransfer_ETH(t *testing.T) {

	walletID := "WGWfGtgo8uLxVKTWoWkzuoPELMtehb9Vda"
	accountID := "6SVDyi7dgcJHG7V2DP2ZJUvzbRyr2PjTqFxym8pEEBrv"
	to := "d35f9Ea14D063af9B3567064FAB567275b09f03D"

	rawTx, err := testCreateTransactionStep(walletID, accountID, to, "0.0003", "")
	if err != nil {
		return
	}

	_, err = testSignTransactionStep(rawTx)
	if err != nil {
		return
	}

	_, err = testVerifyTransactionStep(rawTx)
	if err != nil {
		return
	}

	_, err = testSubmitTransactionStep(rawTx)
	if err != nil {
		return
	}

}

func TestTransfer_LTC(t *testing.T) {

	walletID := "WMTUzB3LWaSKNKEQw9Sn73FjkEoYGHEp4B"
	accountID := "EbUsW3YaHQ61eNt3f4hDXJAFh9LGmLZWH1VTTSnQmhnL"
	to := "n1Prn7ZbZtd5CTN8Yrj4K9c3gD4u8tjFQzX"

	rawTx, err := testCreateTransactionStep(walletID, accountID, to, "1.1", "0.001")
	if err != nil {
		return
	}

	_, err = testSignTransactionStep(rawTx)
	if err != nil {
		return
	}

	_, err = testVerifyTransactionStep(rawTx)
	if err != nil {
		return
	}

	_, err = testSubmitTransactionStep(rawTx)
	if err != nil {
		return
	}

}

func TestTransfer_NAS(t *testing.T) {

	walletID := "VzQTLspxvbXSmfRGcN6LJVB8otYhJwAGWc"
	accountID := "BjLtC1YN4sWQKzYHtNPdvx3D8yVfXmbyeCQTMHv4JUGG"
	to := "n1Prn7ZbZtd5CTN8Yrj4K9c3gD4u8tjFQzX"

	rawTx, err := testCreateTransactionStep(walletID, accountID, to, "0.00005", "")
	if err != nil {
		return
	}

	_, err = testSignTransactionStep(rawTx)
	if err != nil {
		return
	}

	_, err = testVerifyTransactionStep(rawTx)
	if err != nil {
		return
	}

	_, err = testSubmitTransactionStep(rawTx)
	if err != nil {
		return
	}

}
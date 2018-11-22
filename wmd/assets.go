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

package wmd

import (
	"github.com/blocktree/OpenWallet/assets"
	"github.com/blocktree/OpenWallet/assets/bitcoin"
	"github.com/blocktree/OpenWallet/assets/bitcoincash"
	"github.com/blocktree/OpenWallet/assets/bytom"
	"github.com/blocktree/OpenWallet/assets/cardano"
	"github.com/blocktree/OpenWallet/assets/decred"
	"github.com/blocktree/OpenWallet/assets/ethereum"
	"github.com/blocktree/OpenWallet/assets/hypercash"
	"github.com/blocktree/OpenWallet/assets/icon"
	"github.com/blocktree/OpenWallet/assets/iota"
	"github.com/blocktree/OpenWallet/assets/litecoin"
	"github.com/blocktree/OpenWallet/assets/nebulasio"
	"github.com/blocktree/OpenWallet/assets/qtum"
	"github.com/blocktree/OpenWallet/assets/sia"
	"github.com/blocktree/OpenWallet/assets/stc2345"
	"github.com/blocktree/OpenWallet/assets/tezos"
	"github.com/blocktree/OpenWallet/assets/tron"
	"github.com/blocktree/OpenWallet/log"
)

//WalletManagerInterface 钱包管理器
type WalletManagerInterface interface {
	//初始化配置流程
	InitConfigFlow() error
	//查看配置信息
	ShowConfig() error
	//创建钱包流程
	CreateWalletFlow() error
	//创建地址流程
	CreateAddressFlow() error
	//汇总钱包流程
	SummaryFollow() error
	//备份钱包流程
	BackupWalletFlow() error
	//查看钱包列表，显示信息
	GetWalletList() error
	//发送交易
	TransferFlow() error
	//恢复钱包
	RestoreWalletFlow() error
}

// 节点管理接口
type NodeManagerInterface interface {
	// GetNodeStatus 节点状态
	GetNodeStatus(string) error
	// StartNodeFlow 创建节点
	CreateNodeFlow(string) error
	// StartNodeFlow 开启节点
	StartNodeFlow(string) error
	// StopNodeFlow 关闭节点
	StopNodeFlow(string) error
	// RestartNodeFlow 重启节点
	RestartNodeFlow(string) error
	// RemoveNodeFlow 移除节点
	RemoveNodeFlow(string) error
	// LogsNodeFlow 日志
	LogsNodeFlow(string) error

	// //LoginNode 登陆节点
	// LoginNode() error
	// //ShowNodeInfo 显示节点信息
	// ShowNodeInfo() error
}

// 配置管理接口
//type ConfigManagerInterfac interface {
//	//SetConfigFlow 初始化配置流程
//	SetConfigFlow(subCmd string) error
//	//ShowConfigInfo 查看配置信息
//	ShowConfigInfo(subCmd string) error
//}

//注册钱包管理工具
func init() {
	//注册钱包管理工具
	log.Notice("Wallet Manager Driver Load Successfully.")
	assets.RegAssets(cardano.Symbol, &cardano.WalletManager{})
	assets.RegAssets(bytom.Symbol, &bytom.WalletManager{})
	//RegWMD(strings.ToLower(bopo.Symbol), &bopo.WalletManager{})
	assets.RegAssets(bitcoincash.Symbol, &bitcoincash.WalletManager{})
	assets.RegAssets(sia.Symbol, &sia.WalletManager{})
	assets.RegAssets(ethereum.Symbol, ethereum.NewWalletManager())
	assets.RegAssets(stc2345.Symbol, stc2345.NewWalletManager())
	assets.RegAssets(bitcoin.Symbol, bitcoin.NewWalletManager())
	assets.RegAssets(hypercash.Symbol, hypercash.NewWalletManager())
	assets.RegAssets(iota.Symbol, &iota.WalletManager{})
	assets.RegAssets(tezos.Symbol, tezos.NewWalletManager())
	assets.RegAssets(litecoin.Symbol, litecoin.NewWalletManager())
	assets.RegAssets(qtum.Symbol, qtum.NewWalletManager())
	assets.RegAssets(decred.Symbol, decred.NewWalletManager())
	assets.RegAssets(tron.Symbol, tron.NewWalletManager())
	assets.RegAssets(nebulasio.Symbol, nebulasio.NewWalletManager())
	assets.RegAssets(icon.Symbol, icon.NewWalletManager())
}

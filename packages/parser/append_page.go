// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package parser

import (
	"fmt"

	"github.com/GACHAIN/go-gachain-mvp/packages/utils"
	"strings"
)

// AppendPageInit initialize AppendPage transaction
func (p *Parser) AppendPageInit() error {

	fields := []map[string]string{{"global": "int64"}, {"name": "string"}, {"value": "string"}, {"sign": "bytes"}}
	err := p.GetTxMaps(fields)
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}

// AppendPageFront checks conditions of AppendPage transaction
func (p *Parser) AppendPageFront() error {

	err := p.generalCheck(`edit_page`)
	if err != nil {
		fmt.Printf("<<< generalCheck %s\n", err)
		//return p.ErrInfo(err)
	}

	// Check InputData
	/*verifyData := map[string]string{"name": "string", "value": "string", "menu": "string", "conditions": "string"}
	err = p.CheckInputData(verifyData)
	if err != nil {
		return p.ErrInfo(err)
	}*/

	/*
		Check conditions
		...
	*/

	// must be supplemented
	forSign := fmt.Sprintf("%s,%s,%d,%d,%s,%s,%s", p.TxMap["type"], p.TxMap["time"], p.TxCitizenID, p.TxStateID, p.TxMap["global"], p.TxMap["name"], p.TxMap["value"])
	CheckSignResult, err := utils.CheckSign(p.PublicKeys, forSign, p.TxMap["sign"], false)
	if err != nil {
		fmt.Printf("<<< CheckSign %s\n", err)
		//return p.ErrInfo(err)
	}
	if !CheckSignResult {
		fmt.Printf("<<< CheckSignResult %s\n", err)
		//return p.ErrInfo("incorrect sign")
	}
	if err = p.AccessChange(`pages`, p.TxMaps.String["name"]); err != nil {
		if p.AccessRights(`changing_page`, false) != nil {
			fmt.Printf("<<< AccessRights %s\n", err)
			//return err
		}
	}
	return nil
}

// AppendPage proceeds AppendPage transaction
func (p *Parser) AppendPage() error {

	prefix := p.TxStateIDStr
	if p.TxMaps.Int64["global"] == 1 {
		prefix = "global"
	}
	log.Debug("value page", p.TxMaps.String["value"])
	page, err := p.Single(`SELECT value FROM "`+prefix+`_pages" WHERE name = ?`, p.TxMaps.String["name"]).String()
	if err != nil {
		return p.ErrInfo(err)
	}
	new := strings.Replace(page, "PageEnd:", p.TxMaps.String["value"], -1) + "\r\nPageEnd:"
	_, err = p.selectiveLoggingAndUpd([]string{"value"}, []interface{}{new}, prefix+"_pages", []string{"name"}, []string{p.TxMaps.String["name"]}, true)
	if err != nil {
		return p.ErrInfo(err)
	}

	return nil
}

// AppendPageRollback rollbacks AppendPage transaction
func (p *Parser) AppendPageRollback() error {
	return p.autoRollback()
}

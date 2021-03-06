/*
 * Copyright (C) 2015 - 2017 Wiky Lyu
 *
 * This program is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 * See the GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.";
 */

package void

import (
	. "skywalker/agent/base"
)

type VoidServerAgent struct {
	BaseAgent
}

func (*VoidServerAgent) Name() string {
	return "void"
}

func (*VoidServerAgent) OnInit(string, map[string]interface{}) error {
	return nil
}

func (*VoidServerAgent) OnStart() error {
	return nil
}

func (*VoidServerAgent) GetRemoteAddress(host string, port int) (string, int) {
	return "", 0
}

func (*VoidServerAgent) OnConnectResult(int, string, int) (interface{}, interface{}, error) {
	return nil, nil, nil
}

func (*VoidServerAgent) ReadFromServer([]byte) (interface{}, interface{}, error) {
	return nil, nil, nil
}

func (*VoidServerAgent) ReadFromCA([]byte) (interface{}, interface{}, error) {
	return nil, nil, nil
}

func (*VoidServerAgent) OnClose(bool) {}

func (*VoidServerAgent) GetInfo() []map[string]string {
	return []map[string]string{
		map[string]string{
			"void": "",
		},
	}
}

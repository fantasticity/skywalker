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
package walker

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"net"
	. "skywalker/agent/base"
	"skywalker/cipher"
	"strings"
)

const (
	ERROR_DATA_ERROR     = 1
	ERROR_PORT_INVALID   = 2
	ERROR_RESULT_FAILURE = 3
)

const (
	RESULT_SUCCESS = "success"
)

func pack(data []byte) []byte {
	buf := bytes.Buffer{}
	binary.Write(&buf, binary.LittleEndian, int32(len(data)))
	return append(buf.Bytes(), data...)
}

func unpack(data []byte) int {
	var size int32
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &size)
	return int(size)
}

/* 随机生成IV */
func randomKey(ilen int) []byte {
	if ilen <= 0 {
		return nil
	}
	iv := make([]byte, ilen)
	rand.Read(iv)
	return iv
}

func packRequest(addr string, port uint16, method string) ([]byte, cipher.Encrypter) {
	ip := net.ParseIP(addr)
	atype := AType_DOMAINNAME
	if ip != nil {
		if len(ip) == 4 {
			atype = AType_IPV4
		} else {
			atype = AType_IPV6
		}
	}

	info := cipher.GetCipherInfo(strings.ToLower(method))
	key := randomKey(info.KeySize)
	iv := randomKey(info.IvSize)
	req := Request{
		Version: 0x01,
		Atype:   atype,
		Addr:    addr,
		Port:    int32(port),
		Key:     key,
		Iv:      iv,
	}
	data, _ := proto.Marshal(&req)
	return pack(data), info.EncrypterFunc(key, iv)
}

func unpackResponse(data []byte, method string) (*Response, cipher.Decrypter, []byte, error) {
	if len(data) <= 4 {
		return nil, nil, nil, Error(ERROR_DATA_ERROR, "unpack error")
	}
	size := unpack(data)
	if len(data) < size+4 {
		return nil, nil, nil, Error(ERROR_DATA_ERROR, "unpack error")
	}
	rep := &Response{}
	if err := proto.Unmarshal(data[4:4+size], rep); err != nil {
		return nil, nil, nil, err
	}
	info := cipher.GetCipherInfo(strings.ToLower(method))
	key := randomKey(info.KeySize)
	iv := randomKey(info.IvSize)
	return rep, info.DecrypterFunc(key, iv), data[4+size:], nil
}

/*
 * Copyright (C) 2015 Wiky L
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

package protocol


/* 连接远程服务器的结果 */
const (
    CONNECT_OK = "CONNECT_OK"
    CONNECT_UNKNOWN_HOST = "CONNECT_UNKNOWN_HOST"
    CONNECT_UNREACHABLE = "CONNECT_UNREACHABLE"
    CONNECT_UNKNOWN_ERROR = "CONNECT_UNKNOWN_ERROR"
)


type InboundProtocol interface {
    /* 返回协议名 */
    Name() string
    /* 
     * 读取配置，初始化协议
     * 初始化成功，返回true
     * 初始化失败，返回false
     */
    Start(interface{}) bool

    /* 连接结果，只有入站协议才会调用该方法 */
    ConnectResult(string) (interface{}, interface{}, error)

    /*
     * 读取数据
     * 返回的第一个值为转发数据，第二个值为响应数据，第三个值表示出错
     * 数据可以是[]byte也可以是[][]byte。[][]byte回被看做多个[]byte
     * 出错关闭链接
     * 对于入口协议，第一个有效的数据必须指明远程服务器地址
     */
    Read([]byte) (interface{}, interface{}, error)

    /* 关闭链接，释放资源，收尾工作 */
    Close()
}

type OutboundProtocol interface {
    /* 返回协议名 */
    Name() string
    /* 
     * 读取配置，初始化协议
     * 初始化成功，返回true
     * 初始化失败，返回false
     */
    Start(interface{}) bool

    /* 
     * 获取远程地址，参数是入站协议传递过来的远程服务器地址
     * 出战协议可以使用该地址也可以覆盖，使用自己定义的地址
     */
    GetRemoteAddress(string, string) (string, string)

    /*
     * 读取数据
     * 返回的第一个值为转发数据，第二个值为响应数据，第三个值表示出错
     * 数据可以是[]byte也可以是[][]byte。[][]byte回被看做多个[]byte
     * 出错关闭链接
     * 对于入口协议，第一个有效的数据必须指明远程服务器地址
     */
    Read([]byte) (interface{}, interface{}, error)

    /* 关闭链接，释放资源，收尾工作 */
    Close()
}

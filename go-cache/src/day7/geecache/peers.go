/*
-*- encoding: utf-8 -*-
@File    :   peers.go
@Time    :   2022/10/28 12:37:53
@Author  :   victor2022
@Version :   1.0
@Desc    :   interface for distributed cache to pick data
*/
package geecache

import pb "geecache/geecachepb"

/*
@Time    :   2022/10/28 12:43:23
@Author  :   victor2022
@Desc    :   通过key找到对应的PeerGetter

	PeerGetter中可以实现从本地获取，也可以实现从其他节点获取数据的过程
*/
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

/*
@Time    :   2022/10/28 12:42:56
@Author  :   victor2022
@Desc    :   其中的Get方法可以实现从对应组中获取key对应的值
*/
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}

var _ PeerPicker = (*HttpPool)(nil)

// Copyright (C) 2014-2018 Goodrain Co., Ltd.
// RAINBOND, Application Management Platform

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package appruntimesync

import (
	"context"
	"fmt"
	"net"

	"github.com/goodrain/rainbond/pkg/util"

	"github.com/Sirupsen/logrus"
	client "github.com/coreos/etcd/clientv3"
	"github.com/goodrain/rainbond/cmd/worker/option"
	"github.com/goodrain/rainbond/pkg/appruntimesync/pb"
	"github.com/goodrain/rainbond/pkg/appruntimesync/server"
	discover "github.com/goodrain/rainbond/pkg/discover.v2"
	"github.com/goodrain/rainbond/pkg/util/etcd/etcdlock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//AppRuntimeSync app runtime sync modle
// handle app status and event
type AppRuntimeSync struct {
	conf      option.Config
	server    *grpc.Server
	etcdCli   *client.Client
	ctx       context.Context
	cancel    context.CancelFunc
	srss      *server.AppRuntimeSyncServer
	keepalive *discover.KeepAlive
	master    etcdlock.MasterInterface
	hostIP    string
}

//Start start if have master right
//start grpc server
func (a *AppRuntimeSync) Start(errchan chan error) {
	if a.hostIP == "" {
		ip, err := util.LocalIP()
		if err != nil {
			logrus.Errorf("get ip failed,details %s", err.Error())
			errchan <- err
		}
		a.hostIP = ip.String()
	}
	go a.start(errchan)
}
func (a *AppRuntimeSync) start(errchan chan error) {
	for {
		master, err := etcdlock.CreateMasterLock(a.conf.EtcdEndPoints, "/rainbond/workermaster", fmt.Sprintf("%s:%d", a.hostIP, 6535), 10)
		if err != nil {
			errchan <- err
			return
		}
		a.master = master
		master.Start()
	loop:
		for {
			select {
			case event := <-master.EventsChan():
				if event.Type == etcdlock.MasterAdded {
					a.srss.Start()
					go a.startAppRuntimeSync()
					a.registServer()
				}
				if event.Type == etcdlock.MasterDeleted {
					errchan <- fmt.Errorf("worker node %s exit ", fmt.Sprintf("%s:%d", a.hostIP, 6535))
					return
				}
				if event.Type == etcdlock.MasterError {
					a.master.Stop()
					if a.keepalive != nil {
						a.keepalive.Stop()
						a.keepalive = nil
					}
					break loop
				}
			}
		}
		select {
		case <-a.ctx.Done():
			return
		default:
		}
	}
}

//Stop stop app runtime sync server
func (a *AppRuntimeSync) Stop() error {
	a.srss.Stop()
	if a.master != nil {
		a.master.Stop()
	}
	if a.keepalive != nil {
		a.keepalive.Stop()
	}
	return nil
}

//registServer
//regist sync server to etcd
func (a *AppRuntimeSync) registServer() error {
	if a.keepalive == nil {
		keepalive, err := discover.CreateKeepAlive(a.conf.EtcdEndPoints, "app_sync_runtime_server", "", a.conf.HostIP, 6535)
		if err != nil {
			return fmt.Errorf("create app sync server keepalive error,%s", err.Error())
		}
		a.keepalive = keepalive
	}
	return a.keepalive.Start()
}

//CreateAppRuntimeSync create app runtime sync model
func CreateAppRuntimeSync(conf option.Config) *AppRuntimeSync {
	ars := &AppRuntimeSync{
		conf:   conf,
		server: grpc.NewServer(),
		srss:   server.NewAppRuntimeSyncServer(conf),
		hostIP: conf.HostIP,
	}
	pb.RegisterAppRuntimeSyncServer(ars.server, ars.srss)
	// Register reflection service on gRPC server.
	reflection.Register(ars.server)
	return ars
}

//StartAppRuntimeSync start grpc server and regist to etcd
func (a *AppRuntimeSync) startAppRuntimeSync() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 6535))
	if err != nil {
		logrus.Errorf("failed to listen: %v", err)
		return err
	}
	return a.server.Serve(lis)
}

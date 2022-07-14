from mininet.topo import Topo

from mininet.net import Mininet

from mininet.node import CPULimitedHost

from mininet.link import TCLink

from mininet.util import dumpNodeConnections

from mininet.log import setLogLevel

from mininet.node import RemoteController

from mininet.cli import CLI

from threading import Thread

import threading

import time

import os

import random

from subprocess import Popen

import socket

 

CoreSwitchList = []

AggSwitchList = []

EdgeSwitchList = []

HostList = []

 

def topology(k):

  net = Mininet(host=CPULimitedHost, link=TCLink, controller = RemoteController)

  c1 = net.addController('c1',controller=RemoteController,ip='127.0.0.1',port = 6633)  

  POD = k

  pod = POD

  end = pod/2

  iCoreLayerSwitch = (k/2)**2

  iAggLayerSwitch = k*(k/2)

  iEdgeLayerSwitch = k*(k/2)

  iHost = iEdgeLayerSwitch * (k/2)

  SCount = 0

  for x in range(1, pod*(pod/2)+1):

    PREFIX = "s"

    EdgeSwitchList.append(net.addSwitch(PREFIX + str(x)))

    SCount = SCount+1    

    print "ESwitch[",SCount,"]"

 

  for x in range(SCount+1,SCount+pod*(pod/2)+1):

    PREFIX = "s"

    AggSwitchList.append(net.addSwitch(PREFIX + str(x)))

    SCount = SCount+1 

    print "ASwitch[",SCount,"]"

 

  for x in range(SCount+1,SCount+((pod/2)**2)+1):

    PREFIX = "s"

    CoreSwitchList.append(net.addSwitch(PREFIX + str(x)))

    SCount = SCount+1 

    print "CSwitch[",SCount,"]"

 

  f1 = open('/home/ubuntu/pyretic/pyretic/tutorial/f1.txt', 'w')

  count = 0

  digit2 = 0

  digit3 = 0

  for a in range(0,pod):

    for b in range(0,pod/2):

      for c in range(2,2+(pod/2)):

        count = count+1

        digit2 = count/100

        digit3 = count/10000 

        PREFIX = "h"

        #print "digit2:",digit2

        #print "digit3:",digit3

        #print "count:",count

        print "host ip:","10."+str(a)+"."+str(b)+"."+str(c)

        print "host mac:","00:00:00:"+str(digit3%100).zfill(2)+":"+str(digit2%100).zfill(2)+":"+str(count%100).zfill(2)

        f1.write(PREFIX + str(count) + " " + "00:00:00:"+str(digit3%100).zfill(2)+":"+str(digit2%100).zfill(2)+":"+str(count%100).zfill(2)+"\n")

        HostList.append(net.addHost(PREFIX + str(count),ip="10."+str(a)+"."+str(b)+"."+str(c),mac="00:00:00:"+str(digit3%100).zfill(2)+":"+str(digit2%100).zfill(2)+":"+str(count%100).zfill(2)))

  f1.close() 

  f2=open('/home/ubuntu/pyretic/pyretic/tutorial/f2.txt', 'w')

  for x in range(0, iEdgeLayerSwitch):

    for y in range(0,end):

      net.addLink(EdgeSwitchList[x], HostList[end*x+y],bw=10)

      f2.write(str(HostList[end*x+y]) + " " + str(EdgeSwitchList[x])[1] + " " + str(y+1) +"\n")

  f2.close()

 

  print "iAggLayerSwitch=",iAggLayerSwitch

  for x in range(0, iAggLayerSwitch):

    for y in range(0,end):

      net.addLink(AggSwitchList[x], EdgeSwitchList[end*(x/end)+y], bw=10)

 

 

  for x in range(0, iAggLayerSwitch, end):

    for y in range(0,end):

      for z in range(0,end):

        net.addLink(CoreSwitchList[y*end+z], AggSwitchList[x+y], bw=10)

 

  print "*** Starting network"

  net.build()

  c1.start()

  for sw in EdgeSwitchList:

    sw.start([c1])

  for sw in AggSwitchList:

    sw.start([c1])

  for sw in CoreSwitchList:

    sw.start([c1])

 

  print "Dumpling host connections"

  dumpNodeConnections(net.hosts)

 

  #use arp -s to add static mapping of MAC to IP in each host

  print "len(HostList):",len(HostList)

  for x in HostList:

    for y in HostList:

      if x!=y:

        y.cmd('arp -s '+x.IP()+' '+x.MAC())

 

  net.hosts[0].cmd("join /home/ubuntu/pyretic/pyretic/tutorial/f1.txt /home/ubuntu/pyretic/pyretic/tutorial/f2.txt > /home/ubuntu/pyretic/pyretic/tutorial/f3.txt")

  CLI(net)

  net.stop()

 

if __name__ == '__main__':

  #setLogLevel( 'info' )

  topology(4) 
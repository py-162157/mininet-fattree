#!/usr/bin/python
from mininet.topo import Topo
from mininet.net import Mininet
from mininet.node import RemoteController, CPULimitedHost
from mininet.link import TCLink
from mininet.util import dumpNodeConnections


class FatTreeTopo(Topo):
    k = 4
    host_list = []
    switch_list = []
    
    core_switch_object = []
    aggr_switch_object = []
    edge_switch_object = []
    host_object = []
    
    def __init__( self ):
        "Create custom topo."
        # Initialize topology
        Topo.__init__( self )
        # add switches and hosts
        for i in range(self.k):
            # add core switches
            self.core_switch_object.append(self.addSwitch('10_%d_%d_%d' % (self.k, i/2+1, i%2+1)))
            for j in range(self.k):
                #add hosts
                self.host_object.append(self.addHost('10_%d_%d_%d' % (i, j/2, j%2+2), ip='10.%d.%d.%d' % (i, j/2, j%2+2)))
                #add edge switches
                if j < self.k:
                    self.edge_switch_object.append(self.addSwitch('10_%d_%d_1' % (i, j)))
                #add aggr switches
                else:
                    self.aggr_switch_object.append(self.addSwitch('10_%d_%d_1' % (i, j)))
 
        self.host_list.extend(self.host_object)
        self.switch_list.extend(self.core_switch_object)  
        self.switch_list.extend(self.aggr_switch_object)  
        self.switch_list.extend(self.edge_switch_object)  
 
        #add links 
        for switch in self.switch_list:
            a, b, c, d = switch.split("_")
            # core switches: connect to aggr switches
            if int(b) == 4:
                if int(c) == 1:
                    for i in range(4):      
                        self.addLink(switch, "10_%d_2_1" % i)
                elif int(c) == 2:
                    for i in range(4):      
                        self.addLink(switch, "10_%d_3_1" % i)
            # aggr switches: connect to edge switches
                elif int(c) in (2, 3):
                    self.addLink(switch, "10_%d_0_1" % int(b))
                    self.addLink(switch, "10_%d_1_1" % int(b))
            # edge switches: connect to hosts
                elif int(d) == 1:
                    self.addLink(switch, "10_%d_%d_2" % (int(b), int(c)))
                    self.addLink(switch, "10_%d_%d_3" % (int(b), int(c)))

topos = {'mytopo': (lambda: FatTreeTopo())}
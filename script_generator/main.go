package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
)

const (
	ScriptFile4 = "/home/pengyang/mn-scripts/fat-tree-4.py"
)

type element Node

type Edge struct {
	start  element
	end    element
	intf1  int
	intf2  int
	weight uint
}

type Node struct {
	name   string
	weight uint
}

type FatTree struct {
	cores []Node
	pods  []Pod
}

type Pod struct {
	aggregations []Node
	grounds      []Ground
}

type Ground struct {
	access Node
	hosts  []Node
}

type MyEmunet struct {
	nodes      map[string]element
	edges      map[string]Edge
	intf_index map[string]int
	topo_type  string
}

func new_emunet(topo_type string) MyEmunet {
	nodes := make(map[string]element)
	edges := make(map[string]Edge)
	intf_index := make(map[string]int)

	return MyEmunet{nodes, edges, intf_index, topo_type}
}

func (n *MyEmunet) addnode(name string, weight ...int) error {
	if (len(weight) != 1) && (len(weight) != 0) {
		return errors.New("weight bit must be 1 or no weight")
	} else {
		if _, ok := n.nodes[name]; ok {
			return errors.New("node" + name + "exist")
		}
		w := 1
		if len(weight) == 1 {
			w = weight[0]
		}

		n.nodes[name] = element(Node{
			name:   name,
			weight: uint(w),
		})

		n.intf_index[name] = 0
	}

	return nil
}

func (n *MyEmunet) addedge(node1 string, node2 string, weight ...int) error {
	_, ok1 := n.nodes[node1]
	_, ok2 := n.nodes[node2]
	if !(ok1 && ok2) {
		return errors.New(node1 + "and" + node2 + "are not both exist")
	}

	if (len(weight) != 1) && (len(weight) != 0) {
		return errors.New("weight bit must be 1 or no weight")
	} else {
		if _, ok := n.nodes[node1+"-"+node2]; ok {
			return errors.New("edge" + node1 + "-" + node2 + "exist")
		}

		w := 1
		if len(weight) == 1 {
			w = weight[0]
		}

		n.edges[node1+"-"+node2] = Edge{
			start: element(Node{
				name:   node1,
				weight: n.nodes[node1].weight,
			}),
			end: element(Node{
				name:   node2,
				weight: n.nodes[node2].weight,
			}),
			intf1:  n.intf_index[node1],
			intf2:  n.intf_index[node2],
			weight: uint(w),
		}

		n.edges[node2+"-"+node1] = Edge{
			start: element(Node{
				name:   node2,
				weight: n.nodes[node2].weight,
			}),
			end: element(Node{
				name:   node1,
				weight: n.nodes[node1].weight,
			}),
			intf1:  n.intf_index[node2],
			intf2:  n.intf_index[node1],
			weight: uint(w),
		}

		n.intf_index[node1]++
		n.intf_index[node2]++
	}

	return nil
}

func makerange(min, max int) []uint {
	a := make([]uint, max-min)
	for i := range a {
		a[i] = uint(min) + uint(i)
	}
	return a
}

// var f *os.File
// var err error
// if f, err = os.OpenFile(ScriptFile4, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666); err != nil {
// 	fmt.Println("文件打开失败")
// }

// nodename := "s" + strconv.Itoa(switch_count)
// mynet.addnode(nodename, n)
// io.WriteString(f, "        "+nodename+"= self.addHost('"+nodename+"')\n")

func Generate_Fat_Tree_Topo(arg string, random bool) {
	n, _ := strconv.Atoi(arg)
	var fat_tree FatTree

	var f *os.File
	var err error
	if f, err = os.OpenFile(ScriptFile4, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666); err != nil {
		fmt.Println("文件打开失败")
	}

	mynet := new_emunet("fat-tree")
	switch_count := 1
	var i int
	for i = 1; i <= n*n/4; i++ {
		node := Node{
			name:   "s" + strconv.Itoa(switch_count),
			weight: uint(n),
		}
		fat_tree.cores = append(fat_tree.cores, node)
		nodename := "s" + strconv.Itoa(switch_count)
		mynet.addnode(nodename, n)
		io.WriteString(f, "        "+nodename+" = self.addSwitch('"+nodename+"')\n")
		switch_count++
	}
	for i = 1; i <= n; i++ {
		var pod Pod
		for range makerange(1, n/2+1) {
			node := Node{
				name:   "s" + strconv.Itoa(switch_count),
				weight: uint(n),
			}
			pod.aggregations = append(pod.aggregations, node)

			nodename := "s" + strconv.Itoa(switch_count)
			mynet.addnode(nodename, n)
			io.WriteString(f, "        "+nodename+" = self.addSwitch('"+nodename+"')\n")

			switch_count++

			var ground Ground
			node = Node{
				name:   "s" + strconv.Itoa(switch_count),
				weight: uint(n),
			}
			ground.access = node

			nodename = "s" + strconv.Itoa(switch_count)
			mynet.addnode(nodename, n)
			io.WriteString(f, "        "+nodename+" = self.addSwitch('"+nodename+"')\n")

			switch_count++

			host_count := 1
			for range makerange(1, n/2+1) {
				node = Node{
					name:   "h" + strconv.Itoa(host_count) + ground.access.name,
					weight: 1,
				}
				ground.hosts = append(ground.hosts, node)

				nodename = "h" + strconv.Itoa(host_count) + ground.access.name
				mynet.addnode(nodename, 1)
				io.WriteString(f, "        "+nodename+" = self.addHost('"+nodename+"')\n")

				host_count++

			}
			pod.grounds = append(pod.grounds, ground)
		}
		fat_tree.pods = append(fat_tree.pods, pod)
	}

	io.WriteString(f, "\n")

	for i, core := range fat_tree.cores {
		aggregation_count := i * 2 / n
		for _, pod := range fat_tree.pods {
			r := rand.Intn(5) + 1
			weight := uint(1)
			if random {
				weight = uint(r)
			}
			mynet.addedge(core.name, pod.aggregations[aggregation_count].name, int(weight))
			io.WriteString(f, "        self.addLink('"+core.name+"', '"+pod.aggregations[aggregation_count].name+"', bw=1000)\n")
		}
	}

	for _, pod := range fat_tree.pods {
		for _, aggregation := range pod.aggregations {
			for _, ground := range pod.grounds {
				r := rand.Intn(5) + 1
				weight := uint(1)
				if random {
					weight = uint(r)
				}
				mynet.addedge(aggregation.name, ground.access.name, int(weight))
				io.WriteString(f, "        self.addLink('"+aggregation.name+"', '"+ground.access.name+"', bw=1000)\n")
			}
		}
	}

	for _, pod := range fat_tree.pods {
		for _, ground := range pod.grounds {
			for _, host := range ground.hosts {
				r := rand.Intn(5) + 1
				weight := uint(1)
				if random {
					weight = uint(r)
				}
				mynet.addedge(ground.access.name, host.name, int(weight))
				io.WriteString(f, "        self.addLink('"+ground.access.name+"', '"+host.name+"', bw=1000)\n")
			}
		}
	}
}

func main() {
	Generate_Fat_Tree_Topo("6", false)
}

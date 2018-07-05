package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Node struct {
	Renderer    string       `json:"renderer"`
	Name        string       `json:"name"`
	Class       string       `json:"class,omitempty"`
	Updated     string       `json:"updated,omitempty"`
	MaxVolume   string       `json:"maxVolume,omitempty"`
	Nodes       []Node       `json:"nodes,omitempty"`
	Connections []Connection `json:"connections,omitempty"`
}

type Connection struct {
	Source  string   `json:"source"`
	Target  string   `json:"target"`
	Metrics Metric   `json:"metrics"`
	Class   string   `json:"class"`
	Notices []string `json:"notices"`
}

type Metric struct {
	Normal  float32 `json:"normal,omitempty"`
	Warning float32 `json:"warning,omitempty"`
	Danger  float32 `json:"danger,omitempty"`
}

var rAddNode = regexp.MustCompile(`^\+ ?(.*)`)
var rRmNode = regexp.MustCompile(`^- ?(.*)`)
var rAddConn = regexp.MustCompile(`^(.*) ?-> ?(.*): ?(.*)`)
var rRmConn = regexp.MustCompile(`^(.*) ?-> ?(.*)`)
var rStatus = regexp.MustCompile(`^s[tatus]?`)
var rJsonStatus = regexp.MustCompile(`^j[son]?`)
var rClear = regexp.MustCompile(`^c[lear]?`)

var mutex = &sync.Mutex{}

func (ns *Node) addNode(name string) string {
	mutex.Lock()
	defer mutex.Unlock()
	ns.Nodes = append(ns.Nodes, *NewNode(name, "region", "normal"))
	return fmt.Sprintf("Node '%s' added", name)
}

func (ns *Node) removeNode(wanted string) string {
	mutex.Lock()
	defer mutex.Unlock()
	for i, n := range ns.Nodes {
		if n.Name == wanted {
			ns.Nodes = append(ns.Nodes[:i], ns.Nodes[i+1:]...)
			return fmt.Sprintf("Node '%s' removed", wanted)
		}
	}
	return fmt.Sprintf("Node '%s' does not exist", wanted)
}

func (ns *Node) addConnection(source, target string, value float32) string {
	mutex.Lock()
	defer mutex.Unlock()

	for i, conn := range ns.Connections {
		if conn.Target == target && conn.Source == source {
			ns.Connections[i].Metrics.Normal = value
			return fmt.Sprintf("Connection '%s->%s:%.0f' updated", source, target, value)
		}
	}
	ns.Connections = append(ns.Connections, *NewConnection(source, target, value))
	return fmt.Sprintf("Connection '%s->%s:%.0f' added", source, target, value)
}

func (ns *Node) handleInput(in string) string {

	switch {
	case rAddNode.MatchString(in):
		match := rAddNode.FindStringSubmatch(in)
		name := match[1]
		return ns.addNode(name)

	case rRmNode.MatchString(in):
		match := rRmNode.FindStringSubmatch(in)
		name := match[1]
		return ns.removeNode(name)

	case rAddConn.MatchString(in):
		match := rAddConn.FindStringSubmatch(in)
		fmt.Println(match)
		value, err := strconv.ParseFloat(match[3], 32)
		if err != nil {
			return fmt.Sprintf("'%v' is not a valid integer/float", match[3])
		}
		source, target := match[1], match[2]
		return ns.addConnection(source, target, float32(value))

	case rRmConn.MatchString(in):
		return "Connection removed"

	case rStatus.MatchString(in):
		return ns.TextStatus()

	case rJsonStatus.MatchString(in):
		return ns.JsonStatus()

	case rClear.MatchString(in):
		return strings.Repeat("\n", 100)

	default:
		return "Invalid operation"
	}

	return "error"

}

func (n *Node) JsonStatus() string {
	res, err := json.MarshalIndent(n, "", "   ")
	if err != nil {
		log.Fatal(err)
	}
	return string(res)
}

func (n *Node) TextStatus() string {
	var s string
	for _, n := range n.Nodes {
		s += n.Name + "\n"
	}
	return s
}

func NewNode(name, renderer, class string) *Node {
	return &Node{
		Renderer:    renderer,
		Name:        name,
		Class:       class,
		Nodes:       []Node{},
		Connections: []Connection{},
	}
}

func NewConnection(source, target string, value float32) *Connection {
	return &Connection{
		Source: source,
		Target: target,
		Metrics: Metric{
			Normal: value,
			Danger: 92.37},
		Class:   "normal",
		Notices: []string{},
	}
}

func serveJson(rw http.ResponseWriter, req *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Write([]byte(ns.JsonStatus()))
}

var ns = NewNode("edge", "global", "")

func main() {

	for _, val := range []string{"+A", "+B", "+C", "A->C:100", "A->B:100"} {
		fmt.Println(ns.handleInput(val))
	}

	//ns.Connections = append(ns.Connections, *NewConnection("A", "B", 100.0))

	reader := bufio.NewReader(os.Stdin)
	http.HandleFunc("/", serveJson)
	go func() {
		http.ListenAndServe(":8081", nil)
	}()

	for {
		fmt.Print(">")
		in, _ := reader.ReadString('\n')
		if in == "q\n" || in == "quit\n" {
			break
		}

		fmt.Println(ns.handleInput(in))
	}
}

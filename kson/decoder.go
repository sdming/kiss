// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package kson

import (
	"bytes"
)

const (
	stateNone = iota
	stateList
	stateHash
	stateListItem
	stateHashItem
)

func stateName(state int) string {
	switch state {
	case stateNone:
		return "none"
	case stateList:
		return "list"
	case stateHash:
		return "hash"
	case stateListItem:
		return "listitem"
	case stateHashItem:
		return "hashitem"
	}
	panic("error state " + stateName(state))
}

type nodestack struct {
	parent *nodestack
	name   []byte
	state  int
	node   *Node
}

func newNodestack() *nodestack {
	return &nodestack{node: &Node{}}
}

func (stack *nodestack) enter(state int) *nodestack {
	//fmt.Println("nodestack enter", stateName(state))

	var node *Node
	switch state {
	case stateList:
		node = &Node{Type: NodeList, List: make([]*Node, 0, capacity)}
	case stateHash:
		node = &Node{Type: NodeHash, Hash: make(map[string]*Node, capacity)}
	case stateHashItem:
		node = &Node{}
	case stateListItem:
		node = &Node{}
	}
	return &nodestack{parent: stack, state: state, node: node}
}

func (stack *nodestack) exit(state int) *nodestack {
	//fmt.Println("nodestack exit", stateName(state), "parent", stateName(stack.parent.state))

	if stack.state != state {
		panic("node stack exit state does not match " + stateName(state))
	}

	if (state == stateListItem || state == stateHashItem) && stack.node.Type == NodeNone {
		stack.node.Type = NodeLiteral
		//stack.node.Literal = []byte{}
	}

	if state == stateListItem && stack.parent.state == stateList {
		stack.parent.node.List = append(stack.parent.node.List, stack.node)
	} else if state == stateHashItem && stack.parent.state == stateHash {
		stack.parent.node.Hash[string(stack.name)] = stack.node
	} else if stack.parent.state == stateNone {
		stack.parent.node = stack.node
	} else if stack.parent.state == stateListItem || stack.parent.state == stateHashItem {
		stack.parent.node = stack.node
	} else {

	}

	return stack.parent
}

type decoder struct {
	data   []byte
	off    int
	length int
	//buff  bytes.Buffer

	states   [16]int
	stateOff uint
}

func (d *decoder) enterState(state int) int {
	d.stateOff++
	d.states[d.stateOff] = state
	return state
}

func (d *decoder) state() int {
	return d.states[d.stateOff]
}

func (d *decoder) exitState(state int) int {
	if d.state() != state {
		panic("decoder state does not match " + stateName(state))
	}

	if d.stateOff > 0 {
		d.stateOff--
		return d.state()
	}
	panic("decoder exit state error " + stateName(state))
	return stateNone
}

func (d *decoder) readByte() byte {
	if d.off < d.length {
		c := d.data[d.off]
		d.off += 1
		return c
	}
	return eof
}

func (d *decoder) readLine() (end int) {
	end, _ = d.readBytes('\n')
	return
}

func (d *decoder) readBytes(delim byte) (end int, ok bool) {
	i := bytes.IndexByte(d.data[d.off:], delim)
	if i < 0 {
		d.off = d.length
		return d.length, false
	}
	end = d.off + i
	d.off = end + 1
	return end, true
}

func (d *decoder) readBytesofLine(delim byte) (end int, ok bool) {
	i := bytes.IndexAny(d.data[d.off:], string([]byte{delim, '\n'}))
	if i < 0 {
		d.off = d.length
		return d.off - 1, false
	}
	end = d.off + i
	ok = d.data[end] == delim
	d.off = end + 1
	return
}

func (d *decoder) endofLine() bool {
	off := d.off

	for l := d.length; off < l; off++ {
		switch d.data[off] {
		case '\n':
			d.off = off
			return true
		case ' ', '\t', '\r':
			continue
		default:
			d.off = off + 1
			return false
		}
	}
	d.off = off
	return true
}

type FormatError struct {
	Message string
	// column   int
	// row      int
	// position int
}

func (e *FormatError) Error() string {
	return e.Message
}

func isContainerDelim(c byte) bool {
	switch c {
	case '[', ']', '{', '}':
		return true
	}
	return false
}

func isSpace(c byte) bool {
	switch c {
	case ' ', '\t', '\r', '\n':
		return true
	}
	return false
}

func Parse(data []byte) (node *Node, err error) {
	if data == nil || len(data) == 0 {
		err = &FormatError{"data to parse is emty"}
		return
	}

	dec := &decoder{data: data, length: len(data)}
	state, stack := dec.state(), newNodestack()

	for {
		off := dec.off
		c := dec.readByte()
		//fmt.Printf("read state=[%s], off=[%v]; c=[%v],  %c \n", stateName(state), off, c, c)

		if c == eof || c == '\n' { // || c == '\r'
			if state == stateListItem || state == stateHashItem { // 
				state, stack = dec.exitState(state), stack.exit(state)
			}
		}
		if c == eof {
			break
		} else if isSpace(c) {
			continue
		}

		if state == stateHash && c != '}' {
			i, ok := dec.readBytesofLine(':')
			if !ok || off == i {
				err = &FormatError{"hash format error"}
				return
			}
			//fmt.Printf("name:[%s]\n", string(bytes.TrimSpace(dec.data[off:i])))
			state, stack = dec.enterState(stateHashItem), stack.enter(stateHashItem)
			stack.name = bytes.TrimSpace(dec.data[off:i])
			continue
		} else if state == stateList && c != ']' {
			state, stack = dec.enterState(stateListItem), stack.enter(stateListItem)
		}

		if isContainerDelim(c) {
			if !dec.endofLine() {
				err = &FormatError{string(c) + " is not end of line"}
				return
			}
			switch c {
			case '[':
				state, stack = dec.enterState(stateList), stack.enter(stateList)
			case ']':
				state, stack = dec.exitState(stateList), stack.exit(stateList)
			case '{':
				state, stack = dec.enterState(stateHash), stack.enter(stateHash)
			case '}':
				state, stack = dec.exitState(stateHash), stack.exit(stateHash)
			}
			continue
		}

		var value []byte
		if c == '`' || c == '"' {
			if end, ok := dec.readBytes(c); !ok || !dec.endofLine() {
				err = &FormatError{"quote format error"}
				return
			} else {
				value = data[off+1 : end]
			}
		} else {
			//value = bytes.TrimSpace(data[off:dec.readLine()])
			value = bytes.TrimSpace(data[off:dec.readLine()])
		}

		//fmt.Printf("value:[%s] \n", string(value))
		stack.node.Type = NodeLiteral
		stack.node.Literal = string(value)

		if state == stateHashItem || state == stateListItem {
			state, stack = dec.exitState(state), stack.exit(state)
		}
	}

	if dec.stateOff > 0 {
		err = &FormatError{"format error " + stateName(dec.state())}
		return
	}

	return stack.node, nil
}

// Unmarshal unmarshal data to the value pointed to by v.
func Unmarshal(data []byte, v interface{}) error {
	node, err := Parse(data)
	if err != nil {
		return err
	}
	return node.Value(v)
}

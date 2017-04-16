// Copyright 2013 Julien Schmidt. All rights reserved.
// Copyright (c) 2015-2016, 招牌疯子
// Copyright (c) 2017, Kirill Danshin
// Use of this source code is governed by a BSD-style license that can be found
// in the 3rd-Party License/fasthttprouter file.

package gramework

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Slash constant used to minimize string allocations
const Slash = "/"

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func countParams(path string) uint8 {
	var n uint
	for i := zero; i < len(path); i++ {
		if path[i] != ':' && path[i] != '*' {
			continue
		}
		n++
	}
	if n >= 255 {
		return 255
	}
	return uint8(n)
}

type nodeType uint8

const (
	static nodeType = iota // default
	root
	param
	catchAll
)

type node struct {
	path      string
	wildChild bool
	nType     nodeType
	maxParams uint8
	indices   string
	children  []*node
	handle    RequestHandler
	priority  uint32
	hits      uint32
	router    *router
}

// increments priority of the given child and reorders if necessary
func (n *node) incrementChildPrio(pos int) int {
	n.children[pos].priority++
	prio := n.children[pos].priority

	// adjust position (move to front)
	newPos := pos
	for newPos > zero && n.children[newPos-one].priority < prio {
		// swap node positions
		tmpN := n.children[newPos-one]
		n.children[newPos-one] = n.children[newPos]
		n.children[newPos] = tmpN

		newPos--
	}

	// build new index char string
	if newPos != pos {
		n.indices = n.indices[:newPos] + // unchanged prefix, might be empty
			n.indices[pos:pos+one] + // the index char we move
			n.indices[newPos:pos] + n.indices[pos+one:] // rest without char at 'pos'
	}

	return newPos
}

// addRoute adds a node with the given handle to the path.
// Not concurrency-safe!
func (n *node) addRoute(path string, handle RequestHandler, r *router) {
	fullPath := path
	n.priority++
	numParams := countParams(path)
	if n.router == nil {
		n.router = r
	}

	// non-empty tree
	if len(n.path) > zero || len(n.children) > zero {
	walk:
		for {
			// Update maxParams of the current node
			if numParams > n.maxParams {
				n.maxParams = numParams
			}

			// Find the longest common prefix.
			// This also implies that the common prefix contains no ':' or '*'
			// since the existing key can't contain those chars.
			i := zero
			max := min(len(path), len(n.path))
			for i < max && path[i] == n.path[i] {
				i++
			}

			// Split edge
			if i < len(n.path) {
				child := node{
					path:      n.path[i:],
					wildChild: n.wildChild,
					nType:     static,
					indices:   n.indices,
					children:  n.children,
					handle:    n.handle,
					priority:  n.priority - one,
					router:    n.router,
				}

				// Update maxParams (max of all children)
				for i := range child.children {
					if child.children[i].maxParams > child.maxParams {
						child.maxParams = child.children[i].maxParams
					}
				}

				n.children = []*node{&child}
				// []byte for proper unicode char conversion, see #65
				n.indices = string([]byte{n.path[i]})
				n.path = path[:i]
				n.handle = nil
				n.wildChild = false
			}

			// Make new node a child of this node
			if i < len(path) {
				path = path[i:]

				if n.wildChild {
					n = n.children[zero]
					n.priority++

					// Update maxParams of the child node
					if numParams > n.maxParams {
						n.maxParams = numParams
					}
					numParams--

					// Check if the wildcard matches
					if len(path) >= len(n.path) && n.path == path[:len(n.path)] &&
						// Check for longer wildcard, e.g. :name and :names
						(len(n.path) >= len(path) || path[len(n.path)] == SlashByte) {
						continue walk
					} else {
						// Wildcard conflict
						pathSeg := strings.SplitN(path, PathSlash, 2)[zero]
						prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
						panic("'" + pathSeg +
							"' in new path '" + fullPath +
							"' conflicts with existing wildcard '" + n.path +
							"' in existing prefix '" + prefix +
							"'")
					}
				}

				c := path[zero]

				// slash after param
				if n.nType == param && c == SlashByte && len(n.children) == one {
					n = n.children[zero]
					n.priority++
					continue walk
				}

				// Check if a child with the next path byte exists
				for i := zero; i < len(n.indices); i++ {
					if c == n.indices[i] {
						i = n.incrementChildPrio(i)
						n = n.children[i]
						continue walk
					}
				}

				// Otherwise insert it
				if c != ':' && c != '*' {
					// []byte for proper unicode char conversion, see #65
					n.indices += string([]byte{c})
					child := &node{
						maxParams: numParams,
						router:    n.router,
					}
					n.children = append(n.children, child)
					n.incrementChildPrio(len(n.indices) - one)
					n = child
				}
				n.insertChild(numParams, path, fullPath, handle)
				return

			} else if i == len(path) { // Make node a (in-path) leaf
				if n.handle != nil {
					panic("a handle is already registered for path '" + fullPath + "'")
				}
				n.handle = handle
			}
			return
		}
	} else { // Empty tree
		n.insertChild(numParams, path, fullPath, handle)
		n.nType = root
	}
}

func (n *node) insertChild(numParams uint8, path, fullPath string, handle RequestHandler) {
	var offset int // already handled bytes of the path

	// find prefix until first wildcard (beginning with ':' or '*')
	for i, max := zero, len(path); numParams > zero; i++ {
		c := path[i]
		if c != ':' && c != '*' {
			continue
		}

		// find wildcard end (either '/' or path end)
		end := i + one
		for end < max && path[end] != SlashByte {
			switch path[end] {
			// the wildcard name must not contain ':' and '*'
			case ':', '*':
				panic("only one wildcard per path segment is allowed, has: '" +
					path[i:] + "' in path '" + fullPath + "'")
			default:
				end++
			}
		}

		// check if this Node existing children which would be
		// unreachable if we insert the wildcard here
		if len(n.children) > zero {
			panic("wildcard route '" + path[i:end] +
				"' conflicts with existing children in path '" + fullPath + "'")
		}

		// check if the wildcard has a name
		if end-i < 2 {
			panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
		}

		if c == ':' { // param
			// split path at the beginning of the wildcard
			if i > zero {
				n.path = path[offset:i]
				offset = i
			}

			child := &node{
				nType:     param,
				maxParams: numParams,
				router:    n.router,
			}
			n.children = []*node{child}
			n.wildChild = true
			n = child
			n.priority++
			numParams--

			// if the path doesn't end with the wildcard, then there
			// will be another non-wildcard subpath starting with '/'
			if end < max {
				n.path = path[offset:end]
				offset = end

				child := &node{
					maxParams: numParams,
					priority:  one,
					router:    n.router,
				}
				n.children = []*node{child}
				n = child
			}

		} else { // catchAll
			if end != max || numParams > one {
				panic("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
			}

			if len(n.path) > zero && n.path[len(n.path)-one] == SlashByte {
				panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
			}

			// currently fixed width one for '/'
			i--
			if path[i] != SlashByte {
				panic("no / before catch-all in path '" + fullPath + "'")
			}

			n.path = path[offset:i]

			// first node: catchAll node with empty path
			child := &node{
				wildChild: true,
				nType:     catchAll,
				maxParams: one,
				router:    n.router,
			}
			n.children = []*node{child}
			n.indices = string(path[i])
			n = child
			n.priority++

			// second node: node holding the variable
			child = &node{
				path:      path[i:],
				nType:     catchAll,
				maxParams: one,
				handle:    handle,
				priority:  one,
				router:    n.router,
			}
			n.children = []*node{child}

			return
		}
	}

	// insert remaining path part and handle to the leaf
	n.path = path[offset:]
	n.handle = handle
}

// GetValue returns the handle registered with the given path (key). The values
// of wildcards are saved to a map.
// If no handle can be found, a TSR (trailing slash redirect) recommendation is
// made if a handle exists with an extra (without the) trailing slash for the
// given path.
func (n *node) GetValue(reqPath string, ctx *Context, method string) (handle RequestHandler, tsr bool) {
	if n.router == nil {
		panic("no router!")
	}
	if n.router.cache == nil {
		panic("no cache!")
	}
	if record, ok := n.router.cache.Get(reqPath, method); ok {
		for name, value := range record.values {
			ctx.SetUserValue(name, value)
		}
		return record.n.handle, record.tsr
	}
	path := reqPath
walk: // outer loop for walking the tree
	for {
		if len(path) > len(n.path) {
			if path[:len(n.path)] == n.path {
				path = path[len(n.path):]
				// If this node does not have a wildcard (param or catchAll)
				// child,  we can just look up the next child node and continue
				// to walk down the tree
				if !n.wildChild {
					c := path[zero]
					for i := zero; i < len(n.indices); i++ {
						if c == n.indices[i] {
							n = n.children[i]
							continue walk
						}
					}

					// Nothing found.
					// We can recommend to redirect to the same URL without a
					// trailing slash if a leaf exists for that path.
					tsr = (path == PathSlash && n.handle != nil)
					return
				}

				// handle wildcard child
				n = n.children[zero]
				switch n.nType {
				case param:
					// find param end (either '/' or path end)
					end := zero
					for end < len(path) && path[end] != SlashByte {
						end++
					}

					// handle calls to Router.allowed method with nil context
					if ctx != nil {
						ctx.SetUserValue(n.path[one:], path[:end])
					}

					// we need to go deeper!
					if end < len(path) {
						if len(n.children) > zero {
							path = path[end:]
							n = n.children[zero]
							continue walk
						}

						// ... but we can't
						tsr = (len(path) == end+one)
						return
					}

					if handle = n.handle; handle != nil {
						n.hits++
						if n.hits > 32 {
							n.router.cache.PutWild(reqPath, n, tsr, map[string]string{n.path[one:]: path[:end]}, method)
						}
						return
					} else if len(n.children) == one {
						// No handle found. Check if a handle for this path + a
						// trailing slash exists for TSR recommendation
						n = n.children[zero]
						tsr = (n.path == PathSlash && n.handle != nil)
					}

					return

				case catchAll:
					if ctx != nil {
						// save param value
						ctx.SetUserValue(n.path[2:], path)
					}
					handle = n.handle
					n.hits++
					if n.hits > 32 {
						n.router.cache.PutWild(reqPath, n, tsr, map[string]string{n.path[2:]: path}, method)
					}
					return

				default:
					panic("invalid node type")
				}
			}
		} else if path == n.path {
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if handle = n.handle; handle != nil {
				n.router.cache.Put(reqPath, n, tsr, method)
				return
			}

			if path == PathSlash && n.wildChild && n.nType != root {
				tsr = true
				return
			}

			// No handle found. Check if a handle for this path + a
			// trailing slash exists for trailing slash recommendation
			for i := zero; i < len(n.indices); i++ {
				if n.indices[i] == SlashByte {
					n = n.children[i]
					tsr = (len(n.path) == one && n.handle != nil) ||
						(n.nType == catchAll && n.children[zero].handle != nil)
					return
				}
			}

			return
		}

		// Nothing found. We can recommend to redirect to the same URL with an
		// extra trailing slash if a leaf exists for that path
		tsr = (path == PathSlash) ||
			(len(n.path) == len(path)+one && n.path[len(path)] == SlashByte &&
				path == n.path[:len(n.path)-one] && n.handle != nil)
		return
	}
}

// FindCaseInsensitivePath makes a case-insensitive lookup of the given path
// and tries to find a handler.
// It can optionally also fix trailing slashes.
// It returns the case-corrected path and a bool indicating whether the lookup
// was successful.
func (n *node) FindCaseInsensitivePath(path string, fixTrailingSlash bool) ([]byte, bool) {
	return n.findCaseInsensitivePathRec(
		path,
		strings.ToLower(path),
		make([]byte, zero, len(path)+one), // preallocate enough memory for new path
		[4]byte{},                         // empty rune buffer
		fixTrailingSlash,
	)
}

// shift bytes in array by n bytes left
func shiftNRuneBytes(rb [4]byte, n int) [4]byte {
	switch n {
	case zero:
		return rb
	case one:
		return [4]byte{rb[one], rb[2], rb[3], zero}
	case 2:
		return [4]byte{rb[2], rb[3]}
	case 3:
		return [4]byte{rb[3]}
	default:
		return [4]byte{}
	}
}

// recursive case-insensitive lookup function used by n.findCaseInsensitivePath
func (n *node) findCaseInsensitivePathRec(path, loPath string, ciPath []byte, rb [4]byte, fixTrailingSlash bool) ([]byte, bool) {
	loNPath := strings.ToLower(n.path)

walk: // outer loop for walking the tree
	for len(loPath) >= len(loNPath) && (len(loNPath) == zero || loPath[one:len(loNPath)] == loNPath[one:]) {
		// add common path to result
		ciPath = append(ciPath, n.path...)

		if path = path[len(n.path):]; len(path) > zero {
			loOld := loPath
			loPath = loPath[len(loNPath):]

			// If this node does not have a wildcard (param or catchAll) child,
			// we can just look up the next child node and continue to walk down
			// the tree
			if !n.wildChild {
				// skip rune bytes already processed
				rb = shiftNRuneBytes(rb, len(loNPath))

				if rb[zero] != zero {
					// old rune not finished
					for i := zero; i < len(n.indices); i++ {
						if n.indices[i] == rb[zero] {
							// continue with child node
							n = n.children[i]
							loNPath = strings.ToLower(n.path)
							continue walk
						}
					}
				} else {
					// process a new rune
					var rv rune

					// find rune start
					// runes are up to 4 byte long,
					// -4 would definitely be another rune
					var off int
					for max := min(len(loNPath), 3); off < max; off++ {
						if i := len(loNPath) - off; utf8.RuneStart(loOld[i]) {
							// read rune from cached lowercase path
							rv, _ = utf8.DecodeRuneInString(loOld[i:])
							break
						}
					}

					// calculate lowercase bytes of current rune
					utf8.EncodeRune(rb[:], rv)
					// skipp already processed bytes
					rb = shiftNRuneBytes(rb, off)

					for i := zero; i < len(n.indices); i++ {
						// lowercase matches
						if n.indices[i] == rb[zero] {
							// must use a recursive approach since both the
							// uppercase byte and the lowercase byte might exist
							// as an index
							if out, found := n.children[i].findCaseInsensitivePathRec(
								path, loPath, ciPath, rb, fixTrailingSlash,
							); found {
								return out, true
							}
							break
						}
					}

					// same for uppercase rune, if it differs
					if up := unicode.ToUpper(rv); up != rv {
						utf8.EncodeRune(rb[:], up)
						rb = shiftNRuneBytes(rb, off)

						for i := zero; i < len(n.indices); i++ {
							// uppercase matches
							if n.indices[i] == rb[zero] {
								// continue with child node
								n = n.children[i]
								loNPath = strings.ToLower(n.path)
								continue walk
							}
						}
					}
				}

				// Nothing found. We can recommend to redirect to the same URL
				// without a trailing slash if a leaf exists for that path
				return ciPath, (fixTrailingSlash && path == PathSlash && n.handle != nil)
			}

			n = n.children[zero]
			switch n.nType {
			case param:
				// find param end (either '/' or path end)
				k := zero
				for k < len(path) && path[k] != SlashByte {
					k++
				}

				// add param value to case insensitive path
				ciPath = append(ciPath, path[:k]...)

				// we need to go deeper!
				if k < len(path) {
					if len(n.children) > zero {
						// continue with child node
						n = n.children[zero]
						loNPath = strings.ToLower(n.path)
						loPath = loPath[k:]
						path = path[k:]
						continue
					}

					// ... but we can't
					if fixTrailingSlash && len(path) == k+one {
						return ciPath, true
					}
					return ciPath, false
				}

				if n.handle != nil {
					return ciPath, true
				} else if fixTrailingSlash && len(n.children) == one {
					// No handle found. Check if a handle for this path + a
					// trailing slash exists
					n = n.children[zero]
					if n.path == PathSlash && n.handle != nil {
						return append(ciPath, SlashByte), true
					}
				}
				return ciPath, false

			case catchAll:
				return append(ciPath, path...), true

			default:
				panic("invalid node type")
			}
		} else {
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if n.handle != nil {
				return ciPath, true
			}

			// No handle found.
			// Try to fix the path by adding a trailing slash
			if fixTrailingSlash {
				for i := zero; i < len(n.indices); i++ {
					if n.indices[i] == SlashByte {
						n = n.children[i]
						if (len(n.path) == one && n.handle != nil) ||
							(n.nType == catchAll && n.children[zero].handle != nil) {
							return append(ciPath, SlashByte), true
						}
						return ciPath, false
					}
				}
			}
			return ciPath, false
		}
	}

	// Nothing found.
	// Try to fix the path by adding / removing a trailing slash
	if fixTrailingSlash {
		if path == PathSlash {
			return ciPath, true
		}
		if len(loPath)+one == len(loNPath) && loNPath[len(loPath)] == SlashByte &&
			loPath[one:] == loNPath[one:len(loPath)] && n.handle != nil {
			return append(ciPath, n.path...), true
		}
	}
	return ciPath, false
}

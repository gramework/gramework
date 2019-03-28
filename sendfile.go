// Copyright 2017-present Kirill Danshin and Gramework contributors
// Copyright 2019-present Highload LTD (UK CN: 11893420)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

// // Sendfile sends count bytes from f to remote a TCP connection.
// // f offset is always relative to the current offset.
// func Sendfile(conn *net.TCPConn, f *os.File, count int64) (n int64, err error) {
// 	lr := &io.LimitedReader{N: count, R: f}
// 	n, err = conn.ReadFrom(lr)
// 	return
// }

// // Sendfile tries to serve a file with sendfile()
// // otherwise fallbacks to slower copyZeroAlloc
// func (ctx *Context) Sendfile(filepath string) {
// 	ctx.Hijack(sendFileFunc(filepath))
// }

// var (
// 	internalServerErrorBytes = []byte("HTTP/1.1 500 Internal Server error\r\nConnection: close\r\n\r\nInternal Server Error")
// )

// func sendFileFunc(filepath string) func(conn net.Conn) {
// 	return func(conn net.Conn) {
// 		f, err := os.Open(filepath)
// 		if err != nil {
// 			conn.Write(internalServerErrorBytes)
// 		}
// 		fi, err := f.Stat()
// 		if err != nil || fi.IsDir() {
// 			conn.Write(internalServerErrorBytes)
// 		}
// 		if tcpConn, ok := conn.(*net.TCPConn); ok {
// 			Sendfile(tcpConn, f, fi.Size())
// 			log.Printf("Wow! So sendfile, much fast!")
// 			return
// 		}

// 		copyZeroAlloc(conn, f)
// 	}
// }

// func copyZeroAlloc(w io.Writer, r io.Reader) (int64, error) {
// 	vbuf := copyBufPool.Get()
// 	buf := vbuf.([]byte)
// 	n, err := io.CopyBuffer(w, r, buf)
// 	copyBufPool.Put(vbuf)
// 	return n, err
// }

// var copyBufPool = sync.Pool{
// 	New: func() interface{} {
// 		return make([]byte, 4096)
// 	},
// }

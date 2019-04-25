package sctp // OMIT

//go:linkname sysSocket net.sysSocket
func sysSocket(family, sotype, proto int) (int, error)

//go:linkname newNetFD net.newFD
func newNetFD(sysfd, family, sotype int, net string) (unsafe.Pointer, error)

//go:linkname initNetFD net.(*netFD).init
func initNetFD(netfd unsafe.Pointer) error

//go:linkname setDeadline net.(*netFD).SetDeadline
func setDeadline(netfd unsafe.Pointer, t time.Time) error

//go:linkname seReadDeadline net.(*netFD).SetReadDeadline
func setReadDeadline(netfd unsafe.Pointer, t time.Time) error

//go:linkname setWriteDeadline net.(*netFD).SetWriteDeadline
func setWriteDeadline(netfd unsafe.Pointer, t time.Time) error

//go:linkname setsockopt syscall.setsockopt
func setsockopt(s int, level int, name int, val unsafe.Pointer, vallen uintptr) (err error)

//go:linkname getsockopt syscall.getsockopt
func getsockopt(s int, level int, name int, val unsafe.Pointer, vallen *_Socklen) (err error)

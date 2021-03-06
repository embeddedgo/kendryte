// DO NOT EDIT THIS FILE. GENERATED BY xgen.

// +build k210

package timer

import (
	"embedded/mmio"
	"unsafe"

	"github.com/embeddedgo/kendryte/p/bus"
	"github.com/embeddedgo/kendryte/p/mmap"
)

type Periph struct {
	CH              [4]RCH
	_               [20]uint32
	INTSTAT_ALL     RINTSTAT_ALL
	EOI_ALL         REOI_ALL
	RAW_INTSTAT_ALL RRAW_INTSTAT_ALL
	COMP_VERSION    RCOMP_VERSION
	LOAD_COUNT2     [4]RLOAD_COUNT2
}

func TIMER0() *Periph { return (*Periph)(unsafe.Pointer(uintptr(mmap.TIMER0_BASE))) }
func TIMER1() *Periph { return (*Periph)(unsafe.Pointer(uintptr(mmap.TIMER1_BASE))) }
func TIMER2() *Periph { return (*Periph)(unsafe.Pointer(uintptr(mmap.TIMER2_BASE))) }

func (p *Periph) BaseAddr() uintptr {
	return uintptr(unsafe.Pointer(p))
}

func (p *Periph) Bus() bus.Bus {
	return bus.APB0
}

type RCH struct {
	LOAD    RLOAD
	CURRENT RCURRENT
	CONTROL RCONTROL
	EOI     REOI
	INTSTAT RINTSTAT
}

type LOAD uint32

type RLOAD struct{ mmio.U32 }

func (r *RLOAD) LoadBits(mask LOAD) LOAD { return LOAD(r.U32.LoadBits(uint32(mask))) }
func (r *RLOAD) StoreBits(mask, b LOAD)  { r.U32.StoreBits(uint32(mask), uint32(b)) }
func (r *RLOAD) SetBits(mask LOAD)       { r.U32.SetBits(uint32(mask)) }
func (r *RLOAD) ClearBits(mask LOAD)     { r.U32.ClearBits(uint32(mask)) }
func (r *RLOAD) Load() LOAD              { return LOAD(r.U32.Load()) }
func (r *RLOAD) Store(b LOAD)            { r.U32.Store(uint32(b)) }

type RMLOAD struct{ mmio.UM32 }

func (rm RMLOAD) Load() LOAD   { return LOAD(rm.UM32.Load()) }
func (rm RMLOAD) Store(b LOAD) { rm.UM32.Store(uint32(b)) }

type CURRENT uint32

type RCURRENT struct{ mmio.U32 }

func (r *RCURRENT) LoadBits(mask CURRENT) CURRENT { return CURRENT(r.U32.LoadBits(uint32(mask))) }
func (r *RCURRENT) StoreBits(mask, b CURRENT)     { r.U32.StoreBits(uint32(mask), uint32(b)) }
func (r *RCURRENT) SetBits(mask CURRENT)          { r.U32.SetBits(uint32(mask)) }
func (r *RCURRENT) ClearBits(mask CURRENT)        { r.U32.ClearBits(uint32(mask)) }
func (r *RCURRENT) Load() CURRENT                 { return CURRENT(r.U32.Load()) }
func (r *RCURRENT) Store(b CURRENT)               { r.U32.Store(uint32(b)) }

type RMCURRENT struct{ mmio.UM32 }

func (rm RMCURRENT) Load() CURRENT   { return CURRENT(rm.UM32.Load()) }
func (rm RMCURRENT) Store(b CURRENT) { rm.UM32.Store(uint32(b)) }

type CONTROL uint32

type RCONTROL struct{ mmio.U32 }

func (r *RCONTROL) LoadBits(mask CONTROL) CONTROL { return CONTROL(r.U32.LoadBits(uint32(mask))) }
func (r *RCONTROL) StoreBits(mask, b CONTROL)     { r.U32.StoreBits(uint32(mask), uint32(b)) }
func (r *RCONTROL) SetBits(mask CONTROL)          { r.U32.SetBits(uint32(mask)) }
func (r *RCONTROL) ClearBits(mask CONTROL)        { r.U32.ClearBits(uint32(mask)) }
func (r *RCONTROL) Load() CONTROL                 { return CONTROL(r.U32.Load()) }
func (r *RCONTROL) Store(b CONTROL)               { r.U32.Store(uint32(b)) }

type RMCONTROL struct{ mmio.UM32 }

func (rm RMCONTROL) Load() CONTROL   { return CONTROL(rm.UM32.Load()) }
func (rm RMCONTROL) Store(b CONTROL) { rm.UM32.Store(uint32(b)) }

func (p *Periph) ENABLE(n int) RMCONTROL {
	return RMCONTROL{mmio.UM32{&p.CH[n].CONTROL.U32, uint32(ENABLE)}}
}

func (p *Periph) MODE(n int) RMCONTROL {
	return RMCONTROL{mmio.UM32{&p.CH[n].CONTROL.U32, uint32(MODE)}}
}

func (p *Periph) INTERRUPT(n int) RMCONTROL {
	return RMCONTROL{mmio.UM32{&p.CH[n].CONTROL.U32, uint32(INTERRUPT)}}
}

func (p *Periph) PWM_ENABLE(n int) RMCONTROL {
	return RMCONTROL{mmio.UM32{&p.CH[n].CONTROL.U32, uint32(PWM_ENABLE)}}
}

type EOI uint32

type REOI struct{ mmio.U32 }

func (r *REOI) LoadBits(mask EOI) EOI { return EOI(r.U32.LoadBits(uint32(mask))) }
func (r *REOI) StoreBits(mask, b EOI) { r.U32.StoreBits(uint32(mask), uint32(b)) }
func (r *REOI) SetBits(mask EOI)      { r.U32.SetBits(uint32(mask)) }
func (r *REOI) ClearBits(mask EOI)    { r.U32.ClearBits(uint32(mask)) }
func (r *REOI) Load() EOI             { return EOI(r.U32.Load()) }
func (r *REOI) Store(b EOI)           { r.U32.Store(uint32(b)) }

type RMEOI struct{ mmio.UM32 }

func (rm RMEOI) Load() EOI   { return EOI(rm.UM32.Load()) }
func (rm RMEOI) Store(b EOI) { rm.UM32.Store(uint32(b)) }

type INTSTAT uint32

type RINTSTAT struct{ mmio.U32 }

func (r *RINTSTAT) LoadBits(mask INTSTAT) INTSTAT { return INTSTAT(r.U32.LoadBits(uint32(mask))) }
func (r *RINTSTAT) StoreBits(mask, b INTSTAT)     { r.U32.StoreBits(uint32(mask), uint32(b)) }
func (r *RINTSTAT) SetBits(mask INTSTAT)          { r.U32.SetBits(uint32(mask)) }
func (r *RINTSTAT) ClearBits(mask INTSTAT)        { r.U32.ClearBits(uint32(mask)) }
func (r *RINTSTAT) Load() INTSTAT                 { return INTSTAT(r.U32.Load()) }
func (r *RINTSTAT) Store(b INTSTAT)               { r.U32.Store(uint32(b)) }

type RMINTSTAT struct{ mmio.UM32 }

func (rm RMINTSTAT) Load() INTSTAT   { return INTSTAT(rm.UM32.Load()) }
func (rm RMINTSTAT) Store(b INTSTAT) { rm.UM32.Store(uint32(b)) }

type INTSTAT_ALL uint32

type RINTSTAT_ALL struct{ mmio.U32 }

func (r *RINTSTAT_ALL) LoadBits(mask INTSTAT_ALL) INTSTAT_ALL {
	return INTSTAT_ALL(r.U32.LoadBits(uint32(mask)))
}
func (r *RINTSTAT_ALL) StoreBits(mask, b INTSTAT_ALL) { r.U32.StoreBits(uint32(mask), uint32(b)) }
func (r *RINTSTAT_ALL) SetBits(mask INTSTAT_ALL)      { r.U32.SetBits(uint32(mask)) }
func (r *RINTSTAT_ALL) ClearBits(mask INTSTAT_ALL)    { r.U32.ClearBits(uint32(mask)) }
func (r *RINTSTAT_ALL) Load() INTSTAT_ALL             { return INTSTAT_ALL(r.U32.Load()) }
func (r *RINTSTAT_ALL) Store(b INTSTAT_ALL)           { r.U32.Store(uint32(b)) }

type RMINTSTAT_ALL struct{ mmio.UM32 }

func (rm RMINTSTAT_ALL) Load() INTSTAT_ALL   { return INTSTAT_ALL(rm.UM32.Load()) }
func (rm RMINTSTAT_ALL) Store(b INTSTAT_ALL) { rm.UM32.Store(uint32(b)) }

type EOI_ALL uint32

type REOI_ALL struct{ mmio.U32 }

func (r *REOI_ALL) LoadBits(mask EOI_ALL) EOI_ALL { return EOI_ALL(r.U32.LoadBits(uint32(mask))) }
func (r *REOI_ALL) StoreBits(mask, b EOI_ALL)     { r.U32.StoreBits(uint32(mask), uint32(b)) }
func (r *REOI_ALL) SetBits(mask EOI_ALL)          { r.U32.SetBits(uint32(mask)) }
func (r *REOI_ALL) ClearBits(mask EOI_ALL)        { r.U32.ClearBits(uint32(mask)) }
func (r *REOI_ALL) Load() EOI_ALL                 { return EOI_ALL(r.U32.Load()) }
func (r *REOI_ALL) Store(b EOI_ALL)               { r.U32.Store(uint32(b)) }

type RMEOI_ALL struct{ mmio.UM32 }

func (rm RMEOI_ALL) Load() EOI_ALL   { return EOI_ALL(rm.UM32.Load()) }
func (rm RMEOI_ALL) Store(b EOI_ALL) { rm.UM32.Store(uint32(b)) }

type RAW_INTSTAT_ALL uint32

type RRAW_INTSTAT_ALL struct{ mmio.U32 }

func (r *RRAW_INTSTAT_ALL) LoadBits(mask RAW_INTSTAT_ALL) RAW_INTSTAT_ALL {
	return RAW_INTSTAT_ALL(r.U32.LoadBits(uint32(mask)))
}
func (r *RRAW_INTSTAT_ALL) StoreBits(mask, b RAW_INTSTAT_ALL) {
	r.U32.StoreBits(uint32(mask), uint32(b))
}
func (r *RRAW_INTSTAT_ALL) SetBits(mask RAW_INTSTAT_ALL)   { r.U32.SetBits(uint32(mask)) }
func (r *RRAW_INTSTAT_ALL) ClearBits(mask RAW_INTSTAT_ALL) { r.U32.ClearBits(uint32(mask)) }
func (r *RRAW_INTSTAT_ALL) Load() RAW_INTSTAT_ALL          { return RAW_INTSTAT_ALL(r.U32.Load()) }
func (r *RRAW_INTSTAT_ALL) Store(b RAW_INTSTAT_ALL)        { r.U32.Store(uint32(b)) }

type RMRAW_INTSTAT_ALL struct{ mmio.UM32 }

func (rm RMRAW_INTSTAT_ALL) Load() RAW_INTSTAT_ALL   { return RAW_INTSTAT_ALL(rm.UM32.Load()) }
func (rm RMRAW_INTSTAT_ALL) Store(b RAW_INTSTAT_ALL) { rm.UM32.Store(uint32(b)) }

type COMP_VERSION uint32

type RCOMP_VERSION struct{ mmio.U32 }

func (r *RCOMP_VERSION) LoadBits(mask COMP_VERSION) COMP_VERSION {
	return COMP_VERSION(r.U32.LoadBits(uint32(mask)))
}
func (r *RCOMP_VERSION) StoreBits(mask, b COMP_VERSION) { r.U32.StoreBits(uint32(mask), uint32(b)) }
func (r *RCOMP_VERSION) SetBits(mask COMP_VERSION)      { r.U32.SetBits(uint32(mask)) }
func (r *RCOMP_VERSION) ClearBits(mask COMP_VERSION)    { r.U32.ClearBits(uint32(mask)) }
func (r *RCOMP_VERSION) Load() COMP_VERSION             { return COMP_VERSION(r.U32.Load()) }
func (r *RCOMP_VERSION) Store(b COMP_VERSION)           { r.U32.Store(uint32(b)) }

type RMCOMP_VERSION struct{ mmio.UM32 }

func (rm RMCOMP_VERSION) Load() COMP_VERSION   { return COMP_VERSION(rm.UM32.Load()) }
func (rm RMCOMP_VERSION) Store(b COMP_VERSION) { rm.UM32.Store(uint32(b)) }

type LOAD_COUNT2 uint32

type RLOAD_COUNT2 struct{ mmio.U32 }

func (r *RLOAD_COUNT2) LoadBits(mask LOAD_COUNT2) LOAD_COUNT2 {
	return LOAD_COUNT2(r.U32.LoadBits(uint32(mask)))
}
func (r *RLOAD_COUNT2) StoreBits(mask, b LOAD_COUNT2) { r.U32.StoreBits(uint32(mask), uint32(b)) }
func (r *RLOAD_COUNT2) SetBits(mask LOAD_COUNT2)      { r.U32.SetBits(uint32(mask)) }
func (r *RLOAD_COUNT2) ClearBits(mask LOAD_COUNT2)    { r.U32.ClearBits(uint32(mask)) }
func (r *RLOAD_COUNT2) Load() LOAD_COUNT2             { return LOAD_COUNT2(r.U32.Load()) }
func (r *RLOAD_COUNT2) Store(b LOAD_COUNT2)           { r.U32.Store(uint32(b)) }

type RMLOAD_COUNT2 struct{ mmio.UM32 }

func (rm RMLOAD_COUNT2) Load() LOAD_COUNT2   { return LOAD_COUNT2(rm.UM32.Load()) }
func (rm RMLOAD_COUNT2) Store(b LOAD_COUNT2) { rm.UM32.Store(uint32(b)) }

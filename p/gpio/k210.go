// DO NOT EDIT THIS FILE. GENERATED BY svdxgen.

// +build k210

// Package gpio provides access to the registers of the GPIO peripheral.
//
// Instances:
//  GPIO  GPIO_BASE  APB0  GPIO  General Purpose Input/Output Interface
// Registers:
//  0x000 32  DATA_OUTPUT           Data (output) registers
//  0x004 32  DIRECTION             Data direction registers
//  0x008 32  SOURCE                Data source registers
//  0x030 32  INTERRUPT_ENABLE      Interrupt enable/disable registers
//  0x034 32  INTERRUPT_MASK        Interrupt mask registers
//  0x038 32  INTERRUPT_LEVEL       Interrupt level registers
//  0x03C 32  INTERRUPT_POLARITY    Interrupt polarity registers
//  0x040 32  INTERRUPT_STATUS      Interrupt status registers
//  0x044 32  INTERRUPT_STATUS_RAW  Raw interrupt status registers
//  0x048 32  INTERRUPT_DEBOUNCE    Interrupt debounce registers
//  0x04C 32  INTERRUPT_CLEAR       Registers for clearing interrupts
//  0x050 32  DATA_INPUT            External port (data input) registers
//  0x060 32  SYNC_LEVEL            Sync level registers
//  0x064 32  ID_CODE               ID code
//  0x068 32  INTERRUPT_BOTHEDGE    Interrupt both edge type
// Import:
//  github.com/embeddedgo/kendryte/p/bus
//  github.com/embeddedgo/kendryte/p/mmap
package gpio

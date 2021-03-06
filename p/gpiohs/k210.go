// DO NOT EDIT THIS FILE. GENERATED BY svdxgen.

// +build k210

// Package gpiohs provides access to the registers of the GPIOHS peripheral.
//
// Instances:
//  GPIOHS  GPIOHS_BASE  -  GPIOHS0,GPIOHS1,GPIOHS2,GPIOHS3,GPIOHS4,GPIOHS5,GPIOHS6,GPIOHS7,GPIOHS8,GPIOHS9,GPIOHS10,GPIOHS11,GPIOHS12,GPIOHS13,GPIOHS14,GPIOHS15,GPIOHS16,GPIOHS17,GPIOHS18,GPIOHS19,GPIOHS20,GPIOHS21,GPIOHS22,GPIOHS23,GPIOHS24,GPIOHS25,GPIOHS26,GPIOHS27,GPIOHS28,GPIOHS29,GPIOHS30,GPIOHS31  High-speed GPIO
// Registers:
//  0x000 32  INPUT_VAL   Input Value Register
//  0x004 32  INPUT_EN    Pin Input Enable Register
//  0x008 32  OUTPUT_EN   Pin Output Enable Register
//  0x00C 32  OUTPUT_VAL  Output Value Register
//  0x010 32  PULLUP_EN   Internal Pull-Up Enable Register
//  0x014 32  DRIVE       Drive Strength Register
//  0x018 32  RISE_IE     Rise Interrupt Enable Register
//  0x01C 32  RISE_IP     Rise Interrupt Pending Register
//  0x020 32  FALL_IE     Fall Interrupt Enable Register
//  0x024 32  FALL_IP     Fall Interrupt Pending Register
//  0x028 32  HIGH_IE     High Interrupt Enable Register
//  0x02C 32  HIGH_IP     High Interrupt Pending Register
//  0x030 32  LOW_IE      Low Interrupt Enable Register
//  0x034 32  LOW_IP      Low Interrupt Pending Register
//  0x038 32  IOF_EN      HW I/O Function Enable Register
//  0x03C 32  IOF_SEL     HW I/O Function Select Register
//  0x040 32  OUTPUT_XOR  Output XOR (invert) Register
// Import:
//  github.com/embeddedgo/kendryte/p/mmap
package gpiohs

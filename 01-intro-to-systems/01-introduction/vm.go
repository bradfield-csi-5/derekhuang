package vm

const (
	Load  = 0x01
	Store = 0x02
	Add   = 0x03
	Sub   = 0x04
	Halt  = 0xff
)

// Stretch goals
const (
	Addi = 0x05
	Subi = 0x06
	Jump = 0x07
	Beqz = 0x08
)

// Given a 256 byte array of "memory", run the stored program
// to completion, modifying the data in place to reflect the result
//
// The memory format is:
//
// 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f ... ff
// __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ ... __
// ^==DATA===============^ ^==INSTRUCTIONS==============^
func compute(memory []byte) {

	registers := [3]byte{8, 0, 0} // PC, R1 and R2

	// Keep looping, like a physical computer's clock
	for {
		pc := registers[0]
		op := memory[pc]

		// decode and execute
		switch op {
		case Load:
			reg := memory[pc+1]
			addr := memory[pc+2]
			registers[reg] = memory[addr]
			pc += 3
		case Store:
			reg := memory[pc+1]
			addr := memory[pc+2]
			memory[addr] = registers[reg]
			pc += 3
		case Add:
			reg1 := memory[pc+1]
			reg2 := memory[pc+2]
			registers[reg1] += registers[reg2]
			pc += 3
		case Sub:
			reg1 := memory[pc+1]
			reg2 := memory[pc+2]
			registers[reg1] -= registers[reg2]
			pc += 3
		case Addi:
			reg := memory[pc+1]
			registers[reg] += memory[pc+2]
			pc += 3
		case Subi:
			reg := memory[pc+1]
			registers[reg] -= memory[pc+2]
			pc += 3
		case Jump:
			pc = memory[pc+1]
		case Beqz:
			reg := memory[pc+1]
			if registers[reg] == 0 {
				pc += memory[pc+2]
			}
			pc += 3
		case Halt:
			return
		}
		registers[0] = pc
	}
}

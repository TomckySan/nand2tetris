    @switch_state
    M=-1
    D=0
    @SWITCH_STATE_LOOP
    0;JMP

(MAIN_LOOP)
    @KBD
    D=M
    @SWITCH_STATE_LOOP
    D;JEQ // D == 0 is no key input
    D=-1

(SWITCH_STATE_LOOP)
    @temp_state
    M=D

    @switch_state
    D=D-M // -1 or key ascii
    @MAIN_LOOP
    D;JEQ

    @temp_state
    D=M
    @switch_state
    M=D // -1 or 0

    @SCREEN
    D=A
    @8192
    D=D+A
    @i
    M=D

(SCREEN_LOOP)
    @i
    D=M-1
    M=D
    @MAIN_LOOP
    D;JLT
    @switch_state
    D=M
    @i
    A=M
    M=D
    @SCREEN_LOOP
    0;JMP

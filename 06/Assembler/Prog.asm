//    

// 1から100までの整数の和

    @i
    M=1        // i=1
    @sum
    M=0
(LOOP)
    @i
    D=M
    @100
    D=D-A
    @END
    D;JGT
    @i
    D=M // D=i
    @sum
    M=D+M
    @i
    M=M+1
    @LOOP
    0;JMP
(END)
    @END
    0;JMP

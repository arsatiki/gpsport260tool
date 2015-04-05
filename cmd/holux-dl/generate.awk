function clean(line) {
        split(line, components, "--")
        line = components[1]
        match(line, /^ */)
        return substr(line, RLENGTH + 1, length(line) - RLENGTH)
}

function concat(a, b) {
        return (length(a) > 0? a " ": "") clean(b)
}

BEGIN {
        if (OUT == "")
                OUT = "sql_const.go";
}

NF == 0 { next; }
$1 == "--" && NF == 2 { mode = "const"; n = $2; next; }
$1 == "--" && NF == 3 { 
        mode = "array"
        n = $2
        k = $3
        vars[n] = vars[n] FS k
        next
}

mode == "const" {
        consts[n] = concat(consts[n], $0)
}

mode == "array" {
        varkeys[n,k] = concat(varkeys[n,k], $0)
}

END {
        print "package main\n" > OUT
        print "// Generated from setup.sql" >> OUT
        print "const (" >> OUT
        for (c in consts) {
                printf "\t%s = `%s`", c, consts[c] >> OUT
        }
        print "\n)\n" >> OUT

        print "var (" >> OUT
        for (n in vars) {
                printf "\t%s = map[string]string{\n", n >> OUT
                split(vars[n], keys)
                
                for (i in keys) {
                        k = keys[i]
                        printf "\t\t\"%s\": `%s`,\n", k, varkeys[n,k] >> OUT
                }
                
                print "\t}" >> OUT
        }
        print ")" >> OUT
                
}

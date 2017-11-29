package utils

// Decodes cmd line args and returns the hash
// eg -> ./main -v -e production => {"v":"", "e":"production"}
func DecodeCmdLineArgs(args []string) (res map[string]string) {
	res = map[string]string{}
	const dash byte = 0x2d
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg[0] == dash {
			key := arg[1:]
			if i < len(args)-1 {
				if val := args[i+1]; val[0] != dash {
					res[key] = val
					i++
					continue
				}
			}
			res[key] = ""
		}
	}
	return res
}

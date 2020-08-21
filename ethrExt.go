package ethr

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// agreement http tcp udp https icmp
func EthrRun(ip string,agreement string,bcpl string) {
	//
	// If version is not set via ldflags, then default to UNKNOWN
	//
	if gVersion == "" {
		gVersion = "[VERSION: UNKNOWN]"
	}
	//
	// Set GOMAXPROCS to 1024 as running large number of goroutines in a loop
	// to send network traffic results in timer starvation, as well as unfair
	// processing time across goroutines resulting in starvation of many TCP
	// connections. Using a higher number of threads via GOMAXPROCS solves this
	// problem.
	//
	runtime.GOMAXPROCS(1024)

	//flag.Usage = func() { ethrUsage(gVersion) }
	isServer := false
	clientDest :=  ip
	testTypePtr := bcpl // b bandwidth c connections/s p packets/s l 延迟
	thCount := 1
	bufLenStr := ""
	protocol :=  agreement  // 协议
	outputFile := defaultLogFileName
	debug := false
	noOutput := false
	duration :=  5*time.Second // 测试时间
	showUI := false
	rttCount := 1000
	portStr := ""
	modeStr := ""
	use4 := false
	use6 := false
	gap :=  time.Duration(0)
	reverse := false
	ncs := false
	ic := false

	//
	// TODO: Handle the case if there are incorrect arguments
	// fmt.Println("Number of incorrect arguments: " + strconv.Itoa(flag.NArg()))
	//

	//
	// Only used in client mode, to control whether to display per connection
	// statistics or not.
	//
	gNoConnectionStats = ncs

	//
	// Only used in client mode to ignore HTTPS cert errors.
	//
	gIgnoreCert = ic

	if debug {
		loggingLevel = LogLevelDebug
	}

	xMode := false
	switch modeStr {
	case "":
	case "x":
		xMode = true
	default:
		printUsageError("Invalid value for execution mode (-m).")
	}
	mode := ethrModeInv

	if isServer {
		if clientDest != "" {
			printUsageError("Invalid arguments, \"-c\" cannot be used with \"-s\".")
		}
		if xMode {
			mode = ethrModeExtServer
		} else {
			mode = ethrModeServer
		}
	} else if clientDest != "" {
		if xMode {
			mode = ethrModeExtClient
		} else {
			mode = ethrModeClient
		}
	} else {
		printUsageError("Invalid arguments, use either \"-s\" or \"-c\".")
	}

	if reverse && mode != ethrModeClient {
		printUsageError("Invalid arguments, \"-r\" can only be used in client mode.")
	}

	if use4 && !use6 {
		ipVer = ethrIPv4
	} else if use6 && !use4 {
		ipVer = ethrIPv6
	}

	//Default latency test to 1KB if length is not specified
	switch bufLenStr {
	case "":
		bufLenStr = getDefaultBufferLenStr(testTypePtr)
	}

	bufLen := unitToNumber(bufLenStr)
	if bufLen == 0 {
		printUsageError(fmt.Sprintf("Invalid length specified: %s" + bufLenStr))
	}

	if rttCount <= 0 {
		printUsageError(fmt.Sprintf("Invalid RTT count for latency test: %d", rttCount))
	}

	var testType EthrTestType
	switch testTypePtr {
	case "":
		switch mode {
		case ethrModeServer:
			testType = All
		case ethrModeExtServer:
			testType = All
		case ethrModeClient:
			testType = Bandwidth
		case ethrModeExtClient:
			testType = ConnLatency
		}
	case "b":
		testType = Bandwidth
	case "c":
		testType = Cps
	case "p":
		testType = Pps
	case "l":
		testType = Latency
	case "cl":
		testType = ConnLatency
	default:
		printUsageError(fmt.Sprintf("Invalid value \"%s\" specified for parameter \"-t\".\n"+
			"Valid parameters and values are:\n", testTypePtr))
	}

	p := strings.ToUpper(protocol)
	proto := TCP
	switch p {
	case "TCP":
		proto = TCP
	case "UDP":
		proto = UDP
	case "HTTP":
		proto = HTTP
	case "HTTPS":
		proto = HTTPS
	case "ICMP":
		proto = ICMP
	default:
		printUsageError(fmt.Sprintf("Invalid value \"%s\" specified for parameter \"-p\".\n"+
			"Valid parameters and values are:\n", protocol))
	}

	if thCount <= 0 {
		thCount = runtime.NumCPU()
	}

	//
	// For Pkt/s, we always override the buffer size to be just 1 byte.
	// TODO: Evaluate in future, if we need to support > 1 byte packets for
	//       Pkt/s testing.
	//
	if testType == Pps {
		bufLen = 1
	}

	testParam := EthrTestParam{EthrTestID{EthrProtocol(proto), testType},
		uint32(thCount),
		uint32(bufLen),
		uint32(rttCount),
		reverse}
	validateTestParam(mode, testParam)

	generatePortNumbers(portStr)

	logFileName := outputFile
	if !noOutput {
		if logFileName == defaultLogFileName {
			switch mode {
			case ethrModeServer:
				logFileName = "ethrs.log"
			case ethrModeExtServer:
				logFileName = "ethrxs.log"
			case ethrModeClient:
				logFileName = "ethrc.log"
			case ethrModeExtClient:
				logFileName = "ethrxc.log"
			}
		}
		logInit(logFileName)
	}

	clientParam := ethrClientParam{duration, gap}
	serverParam := ethrServerParam{showUI}

	switch mode {
	case ethrModeServer:
		runServer(testParam, serverParam)
	case ethrModeExtServer:
		runXServer(testParam, serverParam)
	case ethrModeClient:
		runClient(testParam, clientParam, clientDest)
	case ethrModeExtClient:
		runXClient(testParam, clientParam, clientDest)
	}
}
func EthrRunServer() {
	//
	// If version is not set via ldflags, then default to UNKNOWN
	//
	if gVersion == "" {
		gVersion = "[VERSION: UNKNOWN]"
	}
	//
	// Set GOMAXPROCS to 1024 as running large number of goroutines in a loop
	// to send network traffic results in timer starvation, as well as unfair
	// processing time across goroutines resulting in starvation of many TCP
	// connections. Using a higher number of threads via GOMAXPROCS solves this
	// problem.
	//
	runtime.GOMAXPROCS(1024)

	//flag.Usage = func() { ethrUsage(gVersion) }
	isServer := true
	clientDest :=  ""
	testTypePtr := ""
	thCount := 1
	bufLenStr := ""
	protocol :=  "tcp"  // 协议
	outputFile := defaultLogFileName
	debug := false
	noOutput := false
	duration :=  5*time.Second // 测试时间
	showUI := false
	rttCount := 1000
	portStr := ""
	modeStr := ""
	use4 := false
	use6 := false
	gap :=  time.Duration(0)
	reverse := false
	ncs := false
	ic := false

	//
	// TODO: Handle the case if there are incorrect arguments
	// fmt.Println("Number of incorrect arguments: " + strconv.Itoa(flag.NArg()))
	//

	//
	// Only used in client mode, to control whether to display per connection
	// statistics or not.
	//
	gNoConnectionStats = ncs

	//
	// Only used in client mode to ignore HTTPS cert errors.
	//
	gIgnoreCert = ic

	if debug {
		loggingLevel = LogLevelDebug
	}

	xMode := false
	switch modeStr {
	case "":
	case "x":
		xMode = true
	default:
		printUsageError("Invalid value for execution mode (-m).")
	}
	mode := ethrModeInv

	if isServer {
		if clientDest != "" {
			printUsageError("Invalid arguments, \"-c\" cannot be used with \"-s\".")
		}
		if xMode {
			mode = ethrModeExtServer
		} else {
			mode = ethrModeServer
		}
	} else if clientDest != "" {
		if xMode {
			mode = ethrModeExtClient
		} else {
			mode = ethrModeClient
		}
	} else {
		printUsageError("Invalid arguments, use either \"-s\" or \"-c\".")
	}

	if reverse && mode != ethrModeClient {
		printUsageError("Invalid arguments, \"-r\" can only be used in client mode.")
	}

	if use4 && !use6 {
		ipVer = ethrIPv4
	} else if use6 && !use4 {
		ipVer = ethrIPv6
	}

	//Default latency test to 1KB if length is not specified
	switch bufLenStr {
	case "":
		bufLenStr = getDefaultBufferLenStr(testTypePtr)
	}

	bufLen := unitToNumber(bufLenStr)
	if bufLen == 0 {
		printUsageError(fmt.Sprintf("Invalid length specified: %s" + bufLenStr))
	}

	if rttCount <= 0 {
		printUsageError(fmt.Sprintf("Invalid RTT count for latency test: %d", rttCount))
	}

	var testType EthrTestType
	switch testTypePtr {
	case "":
		switch mode {
		case ethrModeServer:
			testType = All
		case ethrModeExtServer:
			testType = All
		case ethrModeClient:
			testType = Bandwidth
		case ethrModeExtClient:
			testType = ConnLatency
		}
	case "b":
		testType = Bandwidth
	case "c":
		testType = Cps
	case "p":
		testType = Pps
	case "l":
		testType = Latency
	case "cl":
		testType = ConnLatency
	default:
		printUsageError(fmt.Sprintf("Invalid value \"%s\" specified for parameter \"-t\".\n"+
			"Valid parameters and values are:\n", testTypePtr))
	}

	p := strings.ToUpper(protocol)
	proto := TCP
	switch p {
	case "TCP":
		proto = TCP
	case "UDP":
		proto = UDP
	case "HTTP":
		proto = HTTP
	case "HTTPS":
		proto = HTTPS
	case "ICMP":
		proto = ICMP
	default:
		printUsageError(fmt.Sprintf("Invalid value \"%s\" specified for parameter \"-p\".\n"+
			"Valid parameters and values are:\n", protocol))
	}

	if thCount <= 0 {
		thCount = runtime.NumCPU()
	}

	//
	// For Pkt/s, we always override the buffer size to be just 1 byte.
	// TODO: Evaluate in future, if we need to support > 1 byte packets for
	//       Pkt/s testing.
	//
	if testType == Pps {
		bufLen = 1
	}

	testParam := EthrTestParam{EthrTestID{EthrProtocol(proto), testType},
		uint32(thCount),
		uint32(bufLen),
		uint32(rttCount),
		reverse}
	validateTestParam(mode, testParam)

	generatePortNumbers(portStr)

	logFileName := outputFile
	if !noOutput {
		if logFileName == defaultLogFileName {
			switch mode {
			case ethrModeServer:
				logFileName = "ethrs.log"
			case ethrModeExtServer:
				logFileName = "ethrxs.log"
			case ethrModeClient:
				logFileName = "ethrc.log"
			case ethrModeExtClient:
				logFileName = "ethrxc.log"
			}
		}
		logInit(logFileName)
	}

	clientParam := ethrClientParam{duration, gap}
	serverParam := ethrServerParam{showUI}

	switch mode {
	case ethrModeServer:
		runServer(testParam, serverParam)
	case ethrModeExtServer:
		runXServer(testParam, serverParam)
	case ethrModeClient:
		runClient(testParam, clientParam, clientDest)
	case ethrModeExtClient:
		runXClient(testParam, clientParam, clientDest)
	}
}
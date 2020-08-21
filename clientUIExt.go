//-----------------------------------------------------------------------------
// Copyright (C) Microsoft. All rights reserved.
// Licensed under the MIT license.
// See LICENSE.txt file in the project root for full license information.
//-----------------------------------------------------------------------------
package ethr

import (
	"sync/atomic"
)

type BindWithStruct struct {
	Protocol string `json:"protocol"`
	Bits     string `json:"bits"`
}

//封装带宽数据
type cData struct {
	BandwidthArr []BindWithStruct
	IsDone bool
	IsRunning bool
}
func(c *cData) clearData(){
	c.BandwidthArr = make([]BindWithStruct,0)
	gInterval=0

}
func(c *cData) Init(){
	c.IsDone = false
	c.clearData()
	c.IsRunning = false
}
func getHttpTestResult(test *ethrTest, value uint64, seconds uint64) BindWithStruct {
	var res BindWithStruct
	res.Protocol = protoToString(test.testParam.TestID.Protocol)
	if test.testParam.TestID.Type == Bandwidth && (test.testParam.TestID.Protocol == TCP ||
		test.testParam.TestID.Protocol == UDP) {
		if gInterval == 0 {
			ui.printMsg("- - - - - - - - - - - - - - - - - - - - - - -")
			ui.printMsg("[ ID]   Protocol    Interval      Bits/s")
		}
		cvalue := uint64(0)
		ccount := 0
		test.connListDo(func(ec *ethrConn) {
			val := atomic.SwapUint64(&ec.data, 0)
			val /= seconds
			if !gNoConnectionStats {
				ui.printMsg("[%3d]     %-5s    %03d-%03d sec   %7s", ec.fd,
					protoToString(test.testParam.TestID.Protocol),
					gInterval, gInterval+1, bytesToRate(val))
				res.Bits = bytesToRate(val)
			}
			cvalue += val
			ccount++
		})
		if ccount > 1 || gNoConnectionStats {
			ui.printMsg("[SUM]     %-5s    %03d-%03d sec   %7s",
				protoToString(test.testParam.TestID.Protocol),
				gInterval, gInterval+1, bytesToRate(cvalue))
			res.Bits = bytesToRate(cvalue)
			if !gNoConnectionStats {
				ui.printMsg("- - - - - - - - - - - - - - - - - - - - - - -")
			}
		}
		logResults([]string{test.session.remoteAddr, protoToString(test.testParam.TestID.Protocol),
			bytesToRate(cvalue), "", "", ""})
	} else if test.testParam.TestID.Type == Cps {
		if gInterval == 0 {
			ui.printMsg("- - - - - - - - - - - - - - - - - - - - - - -")
			ui.printMsg("Protocol    Interval      Conn/s")
		}
		ui.printMsg("  %-5s    %03d-%03d sec   %7s",
			protoToString(test.testParam.TestID.Protocol),
			gInterval, gInterval+1, cpsToString(value))
		logResults([]string{test.session.remoteAddr, protoToString(test.testParam.TestID.Protocol),
			"", cpsToString(value), "", ""})

		res.Bits = bytesToRate(value)

	} else if test.testParam.TestID.Type == Pps {
		if gInterval == 0 {
			ui.printMsg("- - - - - - - - - - - - - - - - - - - - - - -")
			ui.printMsg("Protocol    Interval      Pkts/s")
		}
		ui.printMsg("  %-5s    %03d-%03d sec   %7s",
			protoToString(test.testParam.TestID.Protocol),
			gInterval, gInterval+1, ppsToString(value))
		logResults([]string{test.session.remoteAddr, protoToString(test.testParam.TestID.Protocol),
			"", "", ppsToString(value), ""})
		res.Bits = bytesToRate(value)

	} else if test.testParam.TestID.Type == Bandwidth &&
		(test.testParam.TestID.Protocol == HTTP || test.testParam.TestID.Protocol == HTTPS) {
		if gInterval == 0 {
			ui.printMsg("- - - - - - - - - - - - - - - - - - - - - - -")
			ui.printMsg("Protocol    Interval      Bits/s")
		}
		ui.printMsg("  %-5s    %03d-%03d sec   %7s", protoToString(test.testParam.TestID.Protocol), gInterval, gInterval+1, bytesToRate(value))
		logResults([]string{test.session.remoteAddr, protoToString(test.testParam.TestID.Protocol),
			bytesToRate(value), "", "", ""})

		res.Bits = bytesToRate(value)

	}

	gInterval++
	return res
}

// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stfnmllr/go-resp3/client"
)

const (
	primaryDB = iota
	secondDB
	thirdDB
	fourthDB
	maxDB
)

type testFct func(conn client.Conn, ctx *testCTX, t *testing.T)

type fct struct {
	name     string
	f        testFct
	parallel bool
}

var helloRedis = "Hello Redis"

var fcts = []fct{
	{client.CmdDo, testDo, true},
	// Connection
	{client.CmdAuth, testAuth, true},
	{client.CmdEcho, testEcho, true},
	{client.CmdPing, testPing, true},
	{client.CmdSwapdb, testSwapdb, false},
	// Server
	{client.CmdAclCat, testAclCat, true},
	{client.CmdAclDeluser, testAclDeluser, true},
	{client.CmdAclGenpass, testAclGenpass, true},
	{client.CmdAclGetuser, testAclGetuser, true},
	{client.CmdAclHelp, testAclHelp, true},
	{client.CmdAclList, testAclList, true},
	{client.CmdAclLogCount, testAclLogCount, false},
	{client.CmdAclLogReset, testAclLogReset, false},
	{client.CmdAclUsers, testAclUsers, true},
	{client.CmdAclSetuser, testAclSetuser, true},
	{client.CmdAclWhoami, testAclWhoami, true},
	{"ACL", testACL, false},
	{client.CmdBgrewriteaof, testBgrewriteaof, true},
	{client.CmdBgsave, testBgsave, true},
	{client.CmdClientGetname, testClientGetname, true},
	{client.CmdClientId, testClientId, true},
	{client.CmdClientKill, testClientKill, true},
	{client.CmdClientList, testClientList, true},
	{client.CmdClientPause, testClientPause, true},
	{client.CmdClientReply, testClientReply, true},
	{client.CmdClientTracking, testClientTracking, false},
	{client.CmdClientUnblock, testClientUnblock, true},
	{client.CmdCommand, testCommand, true},
	{client.CmdCommandCount, testCommandCount, true},
	{client.CmdCommandGetkeys, testCommandGetkeys, true},
	{client.CmdCommandInfo, testCommandInfo, true},
	{client.CmdConfigGet, testConfigGet, true},
	{client.CmdConfigResetstat, testConfigResetstat, true},
	{client.CmdConfigSet, testConfigSet, true},
	{client.CmdDbsize, testDbsize, true},
	{client.CmdDebugObject, testDebugObject, true},
	{client.CmdInfo, testInfo, true},
	{client.CmdLastsave, testLastsave, true},
	{"Latency", testLatency, false},
	{client.CmdLolwut, testLolwut, true},
	{client.CmdMemoryDoctor, testMemoryDoctor, true},
	{client.CmdMemoryHelp, testMemoryHelp, true},
	{client.CmdMemoryMallocStats, testMemoryMallocStats, true},
	{client.CmdMemoryPurge, testMemoryPurge, true},
	{client.CmdMemoryStats, testMemoryStats, true},
	{client.CmdMemoryUsage, testMemoryUsage, true},
	{client.CmdModuleList, testModuleList, true},
	{client.CmdRole, testRole, true},
	{client.CmdSave, testSave, true},
	{"Slowlog", testSlowlog, false},
	{client.CmdTime, testTime, true},
	// Strings
	{client.CmdAppend, testAppend, true},
	{client.CmdBitcount, testBitcount, true},
	{client.CmdBitfield, testBitfield, true},
	{client.CmdBitopAnd, testBitopAnd, true},
	{client.CmdBitopOr, testBitopOr, true},
	{client.CmdBitopXor, testBitopXor, true},
	{client.CmdBitopNot, testBitopNot, true},
	{client.CmdBitpos, testBitpos, true},
	{client.CmdDecr, testDecr, true},
	{client.CmdDecrby, testDecrby, true},
	{client.CmdGet, testGet, true},
	{client.CmdGetbit, testGetbit, true},
	{client.CmdGetrange, testGetrange, true},
	{client.CmdGetset, testGetset, true},
	{client.CmdIncr, testIncr, true},
	{client.CmdIncrby, testIncrby, true},
	{client.CmdIncrbyfloat, testIncrbyfloat, true},
	{client.CmdMget, testMget, true},
	{client.CmdMset, testMset, true},
	{client.CmdMsetNx, testMsetNx, true},
	{client.CmdSet, testSet, true},
	{client.CmdSetbit, testSetbit, true},
	{client.CmdSetrange, testSetrange, true},
	{client.CmdStralgoLcsStrings, testStralgoLcsStrings, true},
	{client.CmdStralgoLcsLenStrings, testStralgoLcsLenStrings, true},
	{client.CmdStralgoLcsIdxStrings, testStralgoLcsIdxStrings, true},
	{client.CmdStralgoLcsKeys, testStralgoLcsKeys, true},
	{client.CmdStralgoLcsLenKeys, testStralgoLcsLenKeys, true},
	{client.CmdStralgoLcsIdxKeys, testStralgoLcsIdxKeys, true},
	{client.CmdStrlen, testStrlen, true},
	// Keys
	{client.CmdDel, testDel, true},
	{client.CmdDump, testDump, true},
	{client.CmdExists, testExists, true},
	{client.CmdExpire, testExpire, true},
	{client.CmdExpireat, testExpireat, true},
	{client.CmdKeys, testKeys, true},
	{client.CmdMove, testMove, true},
	{client.CmdObjectEncoding, testObjectEncoding, true},
	{client.CmdObjectFreq, testObjectFreq, false},
	{client.CmdObjectHelp, testObjectHelp, true},
	{client.CmdObjectIdletime, testObjectIdletime, false},
	{client.CmdObjectRefcount, testObjectRefcount, true},
	{client.CmdPersist, testPersist, true},
	{client.CmdPexpire, testPexpire, true},
	{client.CmdPexpireat, testPexpireat, true},
	{client.CmdRandomkey, testRandomkey, true},
	{client.CmdRename, testRename, true},
	{client.CmdRenameNx, testRenameNx, true},
	{client.CmdRestore, testRestore, true},
	{client.CmdScan, testScan, true},
	{client.CmdSort, testSort, true},
	{client.CmdTouch, testTouch, true},
	{client.CmdType, testType, true},
	{client.CmdUnlink, testUnlink, true},
	{client.CmdWait, testWait, true},
	// Lists
	{client.CmdBlpop, testBlpop, true},
	{client.CmdBrpop, testBrpop, true},
	{client.CmdBrpoplpush, testBrpoplpush, true},
	{client.CmdLindex, testLindex, true},
	{client.CmdLinsert, testLinsert, true},
	{client.CmdLlen, testLlen, true},
	{client.CmdLpop, testLpop, true},
	{client.CmdLpos, testLpos, true},
	{client.CmdLpush, testLpush, true},
	{client.CmdLpushx, testLpushx, true},
	{client.CmdLrange, testLrange, true},
	{client.CmdLrem, testLrem, true},
	{client.CmdLset, testLset, true},
	{client.CmdLtrim, testLtrim, true},
	{client.CmdRpop, testRpop, true},
	{client.CmdRpoplpush, testRpoplpush, true},
	{client.CmdRpush, testRpush, true},
	{client.CmdRpushx, testRpushx, true},
	// Hashes
	{client.CmdHdel, testHdel, true},
	{client.CmdHexists, testHexists, true},
	{client.CmdHget, testHget, true},
	{client.CmdHgetall, testHgetall, true},
	{client.CmdHincrby, testHincrby, true},
	{client.CmdHincrbyfloat, testHincrbyfloat, true},
	{client.CmdHkeys, testHkeys, true},
	{client.CmdHlen, testHlen, true},
	{client.CmdHmget, testHmget, true},
	{client.CmdHset, testHset, true},
	{client.CmdHsetNx, testHsetNx, true},
	{client.CmdHstrlen, testHstrlen, true},
	{client.CmdHvals, testHvals, true},
	{client.CmdHscan, testHscan, true},
	// pubsub
	{"Pubsub", testPubsub, true},
	// Sets
	{client.CmdSadd, testSadd, true},
	{client.CmdScard, testScard, true},
	{client.CmdSdiff, testSdiff, true},
	{client.CmdSdiffstore, testSdiffstore, true},
	{client.CmdSinter, testSinter, true},
	{client.CmdSinterstore, testSinterstore, true},
	{client.CmdSismember, testSismember, true},
	{client.CmdSmembers, testSmembers, true},
	{client.CmdSmove, testSmove, true},
	{client.CmdSpop, testSpop, true},
	{client.CmdSrandmember, testSrandmember, true},
	{client.CmdSrem, testSrem, true},
	{client.CmdSunion, testSunion, true},
	{client.CmdSunionstore, testSunionstore, true},
	{client.CmdSscan, testSscan, true},
	// Sorted Sets
	{client.CmdBzpopmax, testBzpopmax, true},
	{client.CmdBzpopmin, testBzpopmin, true},
	{client.CmdZadd, testZadd, true},
	{client.CmdZcard, testZcard, true},
	{client.CmdZcount, testZcount, true},
	{client.CmdZincrby, testZincrby, true},
	{client.CmdZinterstore, testZinterstore, true},
	{client.CmdZlexcount, testZlexcount, true},
	{client.CmdZpopmax, testZpopmax, true},
	{client.CmdZpopmin, testZpopmin, true},
	{client.CmdZrange, testZrange, true},
	{client.CmdZrangebylex, testZrangebylex, true},
	{client.CmdZrangebyscore, testZrangebyscore, true},
	{client.CmdZrank, testZrank, true},
	{client.CmdZrem, testZrem, true},
	{client.CmdZremrangebylex, testZremrangebylex, true},
	{client.CmdZremrangebyrank, testZremrangebyrank, true},
	{client.CmdZremrangebyscore, testZremrangebyscore, true},
	{client.CmdZrevrange, testZrevrange, true},
	{client.CmdZrevrangebylex, testZrevrangebylex, true},
	{client.CmdZrevrangebyscore, testZrevrangebyscore, true},
	{client.CmdZrevrank, testZrevrank, true},
	{client.CmdZscan, testZscan, true},
	{client.CmdZscore, testZscore, true},
	{client.CmdZunionstore, testZunionstore, true},
	// Geo
	{client.CmdGeoadd, testGeoadd, true},
	{client.CmdGeodist, testGeodist, true},
	{client.CmdGeohash, testGeohash, true},
	{client.CmdGeopos, testGeopos, true},
	{client.CmdGeoradius, testGeoradius, true},
	{client.CmdGeoradiusbymember, testGeoradiusbymember, true},
	// HyperLogLog
	{client.CmdPfadd, testPfadd, true},
	{client.CmdPfcount, testPfcount, true},
	{client.CmdPfmerge, testPfmerge, true},
	// Scripting
	{client.CmdEval, testEval, true},
	{client.CmdEvalsha, testEvalsha, true},
	{client.CmdScriptExists, testScriptExists, true},
	{client.CmdScriptLoad, testScriptLoad, true},
	// Streams
	{client.CmdXadd, testXadd, true},
	{client.CmdXdel, testXdel, true},
	{client.CmdXgroupCreate, testXgroupCreate, true},
	{client.CmdXgroupSetid, testXgroupSetid, true},
	{client.CmdXgroupDestroy, testXgroupDestroy, true},
	{client.CmdXgroupHelp, testXgroupHelp, true},
	{client.CmdXinfoStream, testXinfoStream, true},
	{client.CmdXinfoHelp, testXinfoHelp, true},
	{client.CmdXlen, testXlen, true},
	{client.CmdXrange, testXrange, true},
	{client.CmdXread, testXread, true},
	{client.CmdXrevrange, testXrevrange, true},
	{client.CmdXtrim, testXtrim, true},
	// Transaction
	{"Transaction", testTransaction, false},
}

type testCTX struct {
	dialer client.Dialer
	keys   []interface{}
	users  []string
}

func newTestCTX(dialer client.Dialer) *testCTX {
	return &testCTX{
		dialer: dialer,
		keys:   make([]interface{}, 0),
		users:  make([]string, 0),
	}
}

func (ctx *testCTX) newKey(s string) string {
	key := client.RandomKey(s)
	ctx.keys = append(ctx.keys, key)
	return key
}

func (ctx *testCTX) newUser(s string) string {
	user := client.RandomKey(s)
	ctx.users = append(ctx.users, user)
	return user
}

func (ctx *testCTX) cleanup(conn client.Conn) {
	for i := int64(0); i < maxDB; i++ {
		conn.Select(i)
		conn.Del(ctx.keys)
		conn.AclDeluser(ctx.users)
	}
	conn.Select(primaryDB)
}

const callerSkip = 1

func assertNil(t *testing.T, v interface{}) {
	_, _, line, _ := runtime.Caller(callerSkip)
	if v != nil {
		t.Fatalf("caller line %d: got %v - expected <nil>", line, v)
	}
}

func assertTrue(t *testing.T, b bool) {
	_, _, line, _ := runtime.Caller(callerSkip)
	if !b {
		t.Fatalf("caller line %d: got <false> - expected <true>", line)
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	_, _, line, _ := runtime.Caller(callerSkip)
	if v == nil {
		t.Fatalf("caller line %d: expected non <nil>", line)
	}
}

func assertEqual(t *testing.T, v1, v2 interface{}) {
	_, _, line, _ := runtime.Caller(callerSkip)
	equal := false

	val1, val2 := reflect.ValueOf(v1), reflect.ValueOf(v2)

	switch val1.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch val2.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			equal = (val1.Int() == val2.Int())
		}

	case reflect.Float32, reflect.Float64:
		switch val2.Kind() {
		case reflect.Float32, reflect.Float64:
			equal = (val1.Float() == val2.Float())
		}

	default:
		equal = reflect.DeepEqual(v1, v2)
	}

	if !equal {
		t.Fatalf("caller line %[1]d: got %[2]T %[2]v - expected %[3]T %[3]v", line, v1, v2)
	}
}

// Redis configuration helper functions
const (
	maxmemoryPolicy         = "maxmemory-policy"
	latencyMonitorThreshold = "latency-monitor-threshold"
)

func setConfig(conn client.Conn, key, value string, t *testing.T) string {
	m, err := conn.ConfigGet(key).ToStringStringMap()
	assertNil(t, err)
	oldValue, ok := m[key]
	assertTrue(t, ok)
	ok, err = conn.ConfigSet(key, value).ToBool()
	assertNil(t, err)
	assertTrue(t, ok)
	return oldValue
}

// Do
func testDo(conn client.Conn, ctx *testCTX, t *testing.T) {
	s, err := conn.Do("PING").ToString()
	assertNil(t, err)
	assertEqual(t, s, "PONG")
	s, err = conn.Do("PING", &helloRedis).ToString()
	assertNil(t, err)
	assertEqual(t, s, helloRedis)
}

// Connection
func testAuth(conn client.Conn, ctx *testCTX, t *testing.T) {
	ok, err := conn.Auth(nil, "doNotKnow").ToBool()
	assertNotNil(t, err)
	assertEqual(t, ok, false)
}

func testEcho(conn client.Conn, ctx *testCTX, t *testing.T) {
	s, err := conn.Echo(helloRedis).ToString()
	assertNil(t, err)
	assertEqual(t, s, helloRedis)
}

func testPing(conn client.Conn, ctx *testCTX, t *testing.T) {
	s, err := conn.Ping(nil).ToString()
	assertNil(t, err)
	assertEqual(t, s, "PONG")
	s, err = conn.Ping(&helloRedis).ToString()
	assertNil(t, err)
	assertEqual(t, s, helloRedis)
}

func testSwapdb(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")

	defer conn.Select(primaryDB)

	err := conn.Select(thirdDB).Err()
	assertNil(t, err)
	err = conn.Set(myKey, "bar").Err()
	assertNil(t, err)
	err = conn.Select(fourthDB).Err()
	assertNil(t, err)
	err = conn.Swapdb(thirdDB, fourthDB).Err()
	assertNil(t, err)
	s, err := conn.Get(myKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "bar")
}

// Server

func testAclCat(conn client.Conn, ctx *testCTX, t *testing.T) {
	cats, err := conn.AclCat(nil).ToStringSlice()
	assertNil(t, err)
	for _, cat := range cats {
		err := conn.AclCat(&cat).Err()
		assertNil(t, err)
	}
}

func testAclDeluser(conn client.Conn, ctx *testCTX, t *testing.T) {
	u1, u2, u3 := ctx.newUser("u1"), ctx.newUser("u2"), ctx.newUser("u3")
	err := conn.AclSetuser(u1, nil).Err()
	assertNil(t, err)
	err = conn.AclSetuser(u2, nil).Err()
	assertNil(t, err)
	i, err := conn.AclDeluser([]string{u1, u2, u3}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
}

func testAclGenpass(conn client.Conn, ctx *testCTX, t *testing.T) {
	s, err := conn.AclGenpass(nil).ToString()
	assertNil(t, err)
	assertEqual(t, len(s), 64)
	s, err = conn.AclGenpass(client.Int64Ptr(32)).ToString()
	assertNil(t, err)
	assertEqual(t, len(s), 8)
	s, err = conn.AclGenpass(client.Int64Ptr(5)).ToString()
	assertNil(t, err)
	assertEqual(t, len(s), 2)
}

func testAclGetuser(conn client.Conn, ctx *testCTX, t *testing.T) {
	myuser := ctx.newUser("myuser")
	b, err := conn.AclSetuser(myuser, nil).ToBool()
	assertNil(t, err)
	assertTrue(t, b)
	_, err = conn.AclGetuser(myuser).ToStringMap()
	assertNil(t, err)
	i, err := conn.AclDeluser([]string{myuser}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
}

func testAclHelp(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.AclHelp().ToStringSlice()
	assertNil(t, err)
}

func testAclList(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.AclList().ToStringSlice()
	assertNil(t, err)
}

func testAclLogCount(conn client.Conn, ctx *testCTX, t *testing.T) {
	user := ctx.newKey("user")
	password := ctx.newKey("password")

	err := conn.Auth(&user, password).Err() // set invalid user / password
	assertNotNil(t, err)
	s, err := conn.AclLogCount(client.Int64Ptr(1)).ToIntfSlice()
	assertNil(t, err)
	assertEqual(t, len(s), 1)
}

func testAclLogReset(conn client.Conn, ctx *testCTX, t *testing.T) {
	user := ctx.newKey("user")
	password := ctx.newKey("password")

	err := conn.Auth(&user, password).Err() // set invalid user / password
	assertNotNil(t, err)
	err = conn.AclLogReset().Err()
	assertNil(t, err)
	s, err := conn.AclLogCount(nil).ToIntfSlice()
	assertNil(t, err)
	assertEqual(t, len(s), 0)
}

func testAclSetuser(conn client.Conn, ctx *testCTX, t *testing.T) {
	myuser := ctx.newUser("myuser")
	b, err := conn.AclSetuser(myuser, nil).ToBool()
	assertNil(t, err)
	assertTrue(t, b)
	b, err = conn.AclSetuser(myuser, []string{"on", ">p1pp0", "~cached:*", "+get"}).ToBool()
	assertNil(t, err)
	assertTrue(t, b)
	i, err := conn.AclDeluser([]string{myuser}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
}

func testAclUsers(conn client.Conn, ctx *testCTX, t *testing.T) {
	err := conn.AclUsers().Err()
	assertNil(t, err)
}

func testAclWhoami(conn client.Conn, ctx *testCTX, t *testing.T) {
	s, err := conn.AclWhoami().ToString()
	assertNil(t, err)
	assertEqual(t, s, "default")
}

// ACL test (cannot be executed in parallel - change user with restrictions - needs exclusive connection
func testACL(conn client.Conn, ctx *testCTX, t *testing.T) {
	conn, err := ctx.dialer.Dial("")
	if err != nil {
		t.Fatal(err)
	}

	myuser := ctx.newUser("myuser")

	b, err := conn.AclSetuser(myuser, []string{"on", ">p1pp0", "~cached:*", "+get"}).ToBool()
	assertNil(t, err)
	assertTrue(t, b)

	b, err = conn.Auth(&myuser, "p1pp0").ToBool()
	assertNil(t, err)
	assertTrue(t, b)

	err = conn.Get("foo").Err()
	assertNotNil(t, err)

	err = conn.Get("cached:1234").Err()
	assertNil(t, err)

	err = conn.Set("cached:1234", "zap").Err()
	assertNotNil(t, err)
}

func testBgrewriteaof(conn client.Conn, ctx *testCTX, t *testing.T) {
	conn.Bgrewriteaof().ToString()
	// might return with error when rewrite was started already
	//assertNil(t, err)
}

func testBgsave(conn client.Conn, ctx *testCTX, t *testing.T) {
	conn.Bgsave(false).ToString()
	conn.Bgsave(true).ToString()
	// might return with error when save is executed in parallel
	// assertNil(t, err)
	// so, just fire the two options
}

func testClientGetname(conn client.Conn, ctx *testCTX, t *testing.T) {
	myName := ctx.newKey("myName")
	ok, err := conn.ClientSetname(myName).ToBool()
	assertNil(t, err)
	assertTrue(t, ok)
	s, err := conn.ClientGetname().ToString()
	assertNil(t, err)
	assertEqual(t, s, myName)
}

func testClientId(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.ClientId().ToInt64()
	assertNil(t, err)
}

func testClientKill(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")

	idCh := make(chan int64, 0)
	done := make(chan (struct{}), 0)

	// Change Key in different connection.
	go func() {
		defer close(done)
		conn, err := client.Dial("")
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		// send id
		id, err := conn.ClientId().ToInt64()
		if err != nil {
			t.Error(err)
			return
		}
		idCh <- id
		// blocking operation.
		err = conn.Brpop([]interface{}{myKey}, 0).Err()
		assertNotNil(t, err)
	}()

	id := <-idCh
	// kill
	i, err := conn.ClientKill(client.Int64Ptr(id), nil, nil, true).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	// wait for go-routine
	<-done
}

func testClientList(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.ClientList(nil).ToString()
	assertNil(t, err)
}

func testClientPause(conn client.Conn, ctx *testCTX, t *testing.T) {
	ok, err := conn.ClientPause(0).ToBool()
	assertNil(t, err)
	assertTrue(t, ok)
}

func testClientReply(conn client.Conn, ctx *testCTX, t *testing.T) {
	// TODO client reply off / client reply skip
	ok, err := conn.ClientReply(client.ReplyModeOn).ToBool()
	assertNil(t, err)
	assertTrue(t, ok)
}

// Client tracking - client side caching
func testClientTracking(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")

	invalidated := make(chan struct{}, 0)

	dialer := ctx.dialer
	dialer.InvalidateCallback = func(keys []string) {
		t.Logf("cache invalidate key %v", keys)

		for _, key := range keys {
			if key == myKey {
				close(invalidated)
				break
			}
		}
	}

	conn, err := dialer.Dial("")
	assertNil(t, err)
	defer conn.Close()

	redir, err := conn.ClientGetredir().ToInt64()
	assertNil(t, err)
	assertEqual(t, redir, -1) // client tracking not active.

	assertNil(t, conn.ClientTracking(true, nil, nil, false, true, false, false).Err()) // client tracking optin.

	redir, err = conn.ClientGetredir().ToInt64()
	assertNil(t, err)
	assertEqual(t, redir, 0) // client caching active but no redirection.

	assertNil(t, conn.Set(myKey, helloRedis).Err())
	// optin -> set caching on for myKey.
	assertNil(t, conn.ClientCaching(true).Err())
	assertNil(t, conn.Get(myKey).Err())

	done := make(chan struct{}, 0)

	// Change Key in different connection.
	go func() {
		conn, err := client.Dial("")
		assertNil(t, err)
		defer conn.Close()
		// Update key.
		assertNil(t, conn.Set(myKey, "Update myKey").Err())
		close(done)
	}()

	<-done        // wait for other connection to change myKey.
	<-invalidated // wait for invalidation notification in callback.

	assertNil(t, conn.ClientTracking(false, nil, nil, false, false, false, false).Err())
}

func testClientUnblock(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")

	idCh := make(chan int64, 0)
	done := make(chan (struct{}), 0)

	// Change Key in different connection.
	go func() {
		defer close(done)
		conn, err := client.Dial("")
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		// send id
		id, err := conn.ClientId().ToInt64()
		if err != nil {
			t.Error(err)
			return
		}
		idCh <- id
		// blocking operation.
		err = conn.Brpop([]interface{}{myKey}, 0).Err()
		if err != nil {
			t.Error(err)
			return
		}
	}()

	id := <-idCh
	// unblock
	for {
		b, err := conn.ClientUnblock(id, nil).ToBool()
		assertNil(t, err)
		if b {
			break
		}
	}
	// wait for go-routine
	<-done
}

func testCommand(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.Command().ToIntfSlice()
	assertNil(t, err)
}

func testCommandCount(conn client.Conn, ctx *testCTX, t *testing.T) {
	slice, err := conn.Command().ToIntfSlice()
	assertNil(t, err)
	i, err := conn.CommandCount().ToInt64()
	assertNil(t, err)
	assertEqual(t, i, len(slice))
}

func testCommandGetkeys(conn client.Conn, ctx *testCTX, t *testing.T) {
	slice, err := conn.CommandGetkeys([]interface{}{"mset", "a", "b", "c", "d", "e", "f"}).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"a", "c", "e"})
	slice, err = conn.CommandGetkeys([]interface{}{"eval", "not consulted", 3, "key1", "key2", "key3", "arg1", "arg2", "arg3", "argN"}).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"key1", "key2", "key3"})
	slice, err = conn.CommandGetkeys([]interface{}{"sort", "mylist", "alpha", "store", "outlist"}).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"mylist", "outlist"})
}

func testCommandInfo(conn client.Conn, ctx *testCTX, t *testing.T) {
	slice, err := conn.CommandInfo([]string{"get", "set", "eval"}).ToIntfSlice()
	assertNil(t, err)
	assertEqual(t, len(slice), 3)
	slice, err = conn.CommandInfo([]string{"foo", "evalsha", "config", "bar"}).ToIntfSlice()
	assertNil(t, err)
	assertEqual(t, len(slice), 4)
	assertNil(t, slice[0])
	assertNil(t, slice[3])
}

func testConfigGet(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.ConfigGet("*").ToStringStringMap()
	assertNil(t, err)
}

func testConfigResetstat(conn client.Conn, ctx *testCTX, t *testing.T) {
	ok, err := conn.ConfigResetstat().ToBool()
	assertNil(t, err)
	assertTrue(t, ok)
}

func testConfigSet(conn client.Conn, ctx *testCTX, t *testing.T) {
	m, err := conn.ConfigGet("loglevel").ToStringStringMap()
	assertNil(t, err)
	for k, v := range m {
		ok, err := conn.ConfigSet(k, v).ToBool()
		assertNil(t, err)
		assertTrue(t, ok)
	}
}

func testDbsize(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.Dbsize().ToInt64()
	assertNil(t, err)
}

func testDebugObject(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "foobar").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	_, err = conn.DebugObject(myKey).ToString()
	assertNil(t, err)
}

func testInfo(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.Info(nil).ToString()
	assertNil(t, err)
}

func testLastsave(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.Lastsave().ToInt64()
	assertNil(t, err)
}

func testLatencyDoctor(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.LatencyDoctor().ToString()
	assertNil(t, err)
}

func testLatency(conn client.Conn, ctx *testCTX, t *testing.T) {
	saveThreshold := setConfig(conn, latencyMonitorThreshold, "100", t)
	defer setConfig(conn, latencyMonitorThreshold, saveThreshold, t)

	_, err := conn.LatencyHelp().ToStringSlice()
	assertNil(t, err)

	err = conn.LatencyReset(nil).Err()
	assertNil(t, err)

	err = conn.Do("debug", "sleep", ".1").Err()
	assertNil(t, err)
	err = conn.Do("debug", "sleep", ".2").Err()
	assertNil(t, err)
	err = conn.Do("debug", "sleep", ".3").Err()
	assertNil(t, err)
	err = conn.Do("debug", "sleep", ".5").Err()
	assertNil(t, err)
	err = conn.Do("debug", "sleep", ".4").Err()
	assertNil(t, err)

	s, err := conn.LatencyGraph("command").ToString()
	assertNil(t, err)
	t.Logf("\nLatency Graph:\n%s", s)

	s, err = conn.LatencyDoctor().ToString()
	assertNil(t, err)
	assertNotNil(t, s)
	//t.Logf("\nLatency Doctor:\n%s", s)

	slice, err := conn.LatencyLatest().ToIntfSlice()
	assertNil(t, err)
	assertNotNil(t, slice)

	slice2, err := conn.LatencyHistory("command").ToIntfSlice2()
	assertNil(t, err)
	assertNotNil(t, slice2)
}

func testLolwut(conn client.Conn, ctx *testCTX, t *testing.T) {
	for i := 5; i <= 6; i++ {
		s, err := conn.Lolwut(client.Int64Ptr(int64(i))).ToString()
		assertNil(t, err)
		t.Logf("\n%s", s)
	}
}

func testMemoryDoctor(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.MemoryDoctor().ToString()
	assertNil(t, err)
}

func testMemoryHelp(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.MemoryHelp().ToStringSlice()
	assertNil(t, err)
}

func testMemoryMallocStats(conn client.Conn, ctx *testCTX, t *testing.T) {
	err := conn.MemoryMallocStats().Err()
	assertNil(t, err)
}

func testMemoryPurge(conn client.Conn, ctx *testCTX, t *testing.T) {
	err := conn.MemoryPurge().Err()
	assertNil(t, err)
}

func testMemoryStats(conn client.Conn, ctx *testCTX, t *testing.T) {
	m, err := conn.MemoryStats().ToMap()
	assertNil(t, err)
	assertNotNil(t, m)
}

func testMemoryUsage(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	err := conn.Set(myKey, "bar").Err()
	assertNil(t, err)
	_, err = conn.MemoryUsage(myKey, nil).ToInt64()
	assertNil(t, err)
}

func testModuleList(conn client.Conn, ctx *testCTX, t *testing.T) {
	err := conn.ModuleList().Err()
	assertNil(t, err)
}

func testRole(conn client.Conn, ctx *testCTX, t *testing.T) {
	slice, err := conn.Role().ToIntfSlice()
	assertNil(t, err)
	assertEqual(t, slice[0], "master")
}

func testSave(conn client.Conn, ctx *testCTX, t *testing.T) {
	conn.Save().ToString()
	// might return with error when save is executed in parallel
	//assertNil(t, err)
}

func testSlowlog(conn client.Conn, ctx *testCTX, t *testing.T) {
	l, err := conn.SlowlogLen().ToInt64()
	assertNil(t, err)
	//slice, err := conn.SlowlogGet(nil).ToSlice() // Redis 6 beta: slowlog get without number of entries returns max. 10 entries
	slice, err := conn.SlowlogGet(client.Int64Ptr(l)).ToIntfSlice()
	assertNil(t, err)
	assertEqual(t, l, len(slice))
	ok, err := conn.SlowlogReset().ToBool()
	assertNil(t, err)
	assertTrue(t, ok)
}

func testTime(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.Time().ToInt64Slice()
	assertNil(t, err)
}

// Strings
func testAppend(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	assertNil(t, conn.Del([]interface{}{myKey}).Err())
	i, err := conn.Append(myKey, "Hello").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 5)
	i, err = conn.Append(myKey, " World").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 11)
	s, err := conn.Get(myKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "Hello World")
}

func testBitcount(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "foobar").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Bitcount(myKey, nil).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 26)
	i, err = conn.Bitcount(myKey, &client.StartEnd{}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 4)
	i, err = conn.Bitcount(myKey, &client.StartEnd{Start: 1, End: 1}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 6)
}

func testBitfield(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	r, err := conn.Bitfield(myKey, []interface{}{client.TypeOffsetIncrement{Type: "i5", Offset: 100, Increment: 1}, client.TypeOffset{Type: "u4", Offset: 0}}).ToInt64Slice()
	assertNil(t, err)
	assertEqual(t, len(r), 2)
	assertEqual(t, r[0], 1)
	assertEqual(t, r[1], 0)

	myKey = ctx.newKey("myKey")
	r, err = conn.Bitfield(myKey, []interface{}{client.TypeOffsetIncrement{Type: "u2", Offset: 100, Increment: 1}, client.OverflowSat, client.TypeOffsetIncrement{Type: "u2", Offset: 102, Increment: 1}}).ToInt64Slice()
	assertNil(t, err)
	assertEqual(t, len(r), 2)
	assertEqual(t, r[0], 1)
	assertEqual(t, r[1], 1)
	r, err = conn.Bitfield(myKey, []interface{}{client.TypeOffsetIncrement{Type: "u2", Offset: 100, Increment: 1}, client.OverflowSat, client.TypeOffsetIncrement{Type: "u2", Offset: 102, Increment: 1}}).ToInt64Slice()
	assertNil(t, err)
	assertEqual(t, len(r), 2)
	assertEqual(t, r[0], 2)
	assertEqual(t, r[1], 2)
	r, err = conn.Bitfield(myKey, []interface{}{client.TypeOffsetIncrement{Type: "u2", Offset: 100, Increment: 1}, client.OverflowSat, client.TypeOffsetIncrement{Type: "u2", Offset: 102, Increment: 1}}).ToInt64Slice()
	assertNil(t, err)
	assertEqual(t, len(r), 2)
	assertEqual(t, r[0], 3)
	assertEqual(t, r[1], 3)
	r, err = conn.Bitfield(myKey, []interface{}{client.TypeOffsetIncrement{Type: "u2", Offset: 100, Increment: 1}, client.OverflowSat, client.TypeOffsetIncrement{Type: "u2", Offset: 102, Increment: 1}}).ToInt64Slice()
	assertNil(t, err)
	assertEqual(t, len(r), 2)
	assertEqual(t, r[0], 0)
	assertEqual(t, r[1], 3)
}

func testBitopAnd(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2, dest := ctx.newKey("key1"), ctx.newKey("key2"), ctx.newKey("dest")
	ok, err := conn.Set(key1, "foobar").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	ok, err = conn.Set(key2, "abcdef").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	err = conn.BitopAnd(dest, nil).Err() // check invalid command (no keys)
	assertNotNil(t, err)
	i, err := conn.BitopAnd(dest, []interface{}{key1, key2}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 6)
	s, err := conn.Get(dest).ToString()
	assertNil(t, err)
	assertEqual(t, s, "`bc`ab")
}

func testBitopOr(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2, dest := ctx.newKey("key1"), ctx.newKey("key2"), ctx.newKey("dest")
	ok, err := conn.Set(key1, "foobar").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	ok, err = conn.Set(key2, "abcdef").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	err = conn.BitopOr(dest, nil).Err() // check invalid command (no keys)
	assertNotNil(t, err)
	i, err := conn.BitopOr(dest, []interface{}{key1, key2}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 6)
	s, err := conn.Get(dest).ToString()
	assertNil(t, err)
	assertEqual(t, s, "goofev")
}

func testBitopXor(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2, dest := ctx.newKey("key1"), ctx.newKey("key2"), ctx.newKey("dest")
	ok, err := conn.Set(key1, "foobar").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	ok, err = conn.Set(key2, "abcdef").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	err = conn.BitopXor(dest, nil).Err() // check invalid command (no keys)
	assertNotNil(t, err)
	i, err := conn.BitopXor(dest, []interface{}{key1, key2}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 6)
	s, err := conn.Get(dest).ToString()
	assertNil(t, err)
	assertEqual(t, s, string([]byte{'\a', '\r', '\x0c', '\x06', '\x04', '\x14'}))
}

func testBitopNot(conn client.Conn, ctx *testCTX, t *testing.T) {
	key, dest := ctx.newKey("key"), ctx.newKey("dest")
	ok, err := conn.Set(key, "foobar").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.BitopNot(dest, key).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 6)
	s, err := conn.Get(dest).ToString()
	assertNil(t, err)
	assertEqual(t, s, string([]byte{'\x99', '\x90', '\x90', '\x9d', '\x9e', '\x8d'}))
}

func testBitpos(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "\xff\xf0\x00").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Bitpos(myKey, 0, nil, nil).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 12)

	ok, err = conn.Set(myKey, "\x00\xff\xf0").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err = conn.Bitpos(myKey, 1, client.Int64Ptr(0), nil).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 8)
	i, err = conn.Bitpos(myKey, 1, client.Int64Ptr(2), nil).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 16)

	ok, err = conn.Set(myKey, "\x00\x00\x00").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err = conn.Bitpos(myKey, 1, nil, nil).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, -1)
}

func testDecr(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "10").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Decr(myKey).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 9)
	ok, err = conn.Set(myKey, "234293482390480948029348230948").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	err = conn.Decr(myKey).Err()
	assertNotNil(t, err)
}

func testDecrby(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "10").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Decrby(myKey, 3).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 7)
}

func testGet(conn client.Conn, ctx *testCTX, t *testing.T) {
	notExisting, myKey := ctx.newKey("notExisting"), ctx.newKey("myKey")
	b, err := conn.Get(notExisting).IsNull()
	assertNil(t, err)
	assertTrue(t, b)
	ok, err := conn.Set(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	s, err := conn.Get(myKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "Hello")
}

func testGetbit(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	i, err := conn.Setbit(myKey, 7, 1).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
	i, err = conn.Getbit(myKey, 0).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
	i, err = conn.Getbit(myKey, 7).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Getbit(myKey, 100).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
}

func testGetrange(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "This is a string").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	s, err := conn.Getrange(myKey, 0, 3).ToString()
	assertNil(t, err)
	assertEqual(t, s, "This")
	s, err = conn.Getrange(myKey, -3, -1).ToString()
	assertNil(t, err)
	assertEqual(t, s, "ing")
	s, err = conn.Getrange(myKey, 0, -1).ToString()
	assertNil(t, err)
	assertEqual(t, s, "This is a string")
	s, err = conn.Getrange(myKey, 10, 100).ToString()
	assertNil(t, err)
	assertEqual(t, s, "string")
}

func testGetset(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	s, err := conn.Getset(myKey, "World").ToString()
	assertNil(t, err)
	assertEqual(t, s, "Hello")
	s, err = conn.Get(myKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "World")
}

func testIncr(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "10").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Incr(myKey).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 11)
	s, err := conn.Get(myKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "11")
}

func testIncrby(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "10").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Incrby(myKey, 5).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 15)
}

func testIncrbyfloat(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "10.50").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	f, err := conn.Incrbyfloat(myKey, 0.1).ToFloat64()
	assertNil(t, err)
	assertEqual(t, f, 10.6)
	f, err = conn.Incrbyfloat(myKey, -5).ToFloat64()
	assertNil(t, err)
	assertEqual(t, f, 5.6)
	ok, err = conn.Set(myKey, "5.0e3").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	f, err = conn.Incrbyfloat(myKey, 2.0e2).ToFloat64()
	assertNil(t, err)
	assertEqual(t, f, 5200.0)
}

func testMget(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2, nonExisting := ctx.newKey("key1"), ctx.newKey("key2"), ctx.newKey("nonExisting")
	ok, err := conn.Set(key1, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	ok, err = conn.Set(key2, "World").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	slice, err := conn.Mget([]interface{}{key1, key2, nonExisting}).ToIntfSlice()
	assertNil(t, err)
	assertEqual(t, slice, []interface{}{"Hello", "World", nil})
}

func testMset(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("key1"), ctx.newKey("key2")
	ok, err := conn.Mset([]client.KeyValue{{Key: key1, Value: "Hello"}, {Key: key2, Value: "World"}}).ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	s, err := conn.Get(key1).ToString()
	assertNil(t, err)
	assertEqual(t, s, "Hello")
	s, err = conn.Get(key2).ToString()
	assertNil(t, err)
	assertEqual(t, s, "World")
}

func testMsetNx(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2, key3 := ctx.newKey("key1"), ctx.newKey("key2"), ctx.newKey("key3")
	i, err := conn.MsetNx([]client.KeyValue{{Key: key1, Value: "Hello"}, {Key: key2, Value: "there"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.MsetNx([]client.KeyValue{{Key: key2, Value: "new"}, {Key: key3, Value: "World"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
	slice, err := conn.Mget([]interface{}{key1, key2, key3}).ToIntfSlice()
	assertNil(t, err)
	assertEqual(t, slice, []interface{}{"Hello", "there", nil})
}

func testSet(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")

	// SetNx
	err := conn.Del([]interface{}{myKey}).Err()
	assertNil(t, err)
	ok, err := conn.SetNx(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	b, err := conn.SetNx(myKey, "World").IsNull()
	assertNil(t, err)
	assertTrue(t, b)
	s, err := conn.Get(myKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "Hello")

	// SetXx
	err = conn.Del([]interface{}{myKey}).Err()
	assertNil(t, err)
	b, err = conn.SetXx(myKey, "Hello").IsNull()
	assertNil(t, err)
	assertTrue(t, b)
	ok, err = conn.SetNx(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	ok, err = conn.SetXx(myKey, "World").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	s, err = conn.Get(myKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "World")

	// SetEx
	err = conn.Del([]interface{}{myKey}).Err()
	assertNil(t, err)
	ok, err = conn.SetEx(myKey, "Hello", 10).ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.TTL(myKey).ToInt64()
	assertNil(t, err)
	t.Logf("Key expires in %d seconds", i)
	s, err = conn.Get(myKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "Hello")

	// SetExNx
	err = conn.Del([]interface{}{myKey}).Err()
	assertNil(t, err)
	ok, err = conn.SetExNx(myKey, "Hello", 10).ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	b, err = conn.SetExNx(myKey, "World", 10).IsNull()
	assertNil(t, err)
	assertTrue(t, b)
	i, err = conn.TTL(myKey).ToInt64()
	assertNil(t, err)
	t.Logf("Key expires in %d seconds", i)
	s, err = conn.Get(myKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "Hello")

	// SetExXx
	err = conn.Del([]interface{}{myKey}).Err()
	assertNil(t, err)
	b, err = conn.SetExXx(myKey, "Hello", 10).IsNull()
	assertNil(t, err)
	assertTrue(t, b)
	ok, err = conn.SetExNx(myKey, "Hello", 10).ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	ok, err = conn.SetExXx(myKey, "World", 10).ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err = conn.TTL(myKey).ToInt64()
	assertNil(t, err)
	t.Logf("Key expires in %d seconds", i)
	s, err = conn.Get(myKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "World")

	// SetPx
	err = conn.Del([]interface{}{myKey}).Err()
	assertNil(t, err)
	ok, err = conn.SetPx(myKey, "Hello", 1000).ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err = conn.PTTL(myKey).ToInt64()
	assertNil(t, err)
	t.Logf("Key expires in %d milliseconds", i)
	s, err = conn.Get(myKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "Hello")

	// SetPxNx
	err = conn.Del([]interface{}{myKey}).Err()
	assertNil(t, err)
	ok, err = conn.SetPxNx(myKey, "Hello", 1000).ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	b, err = conn.SetPxNx(myKey, "World", 1000).IsNull()
	assertNil(t, err)
	assertTrue(t, b)
	i, err = conn.PTTL(myKey).ToInt64()
	assertNil(t, err)
	t.Logf("Key expires in %d milliseconds", i)
	s, err = conn.Get(myKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "Hello")

	// SetPxXx
	err = conn.Del([]interface{}{myKey}).Err()
	assertNil(t, err)
	b, err = conn.SetPxXx(myKey, "Hello", 1000).IsNull()
	assertNil(t, err)
	assertTrue(t, b)
	ok, err = conn.SetPxNx(myKey, "Hello", 1000).ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	ok, err = conn.SetPxXx(myKey, "World", 1000).ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err = conn.PTTL(myKey).ToInt64()
	assertNil(t, err)
	t.Logf("Key expires in %d milliseconds", i)
	s, err = conn.Get(myKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "World")
}

func testSetbit(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	i, err := conn.Setbit(myKey, 7, 1).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
	i, err = conn.Setbit(myKey, 7, 0).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	s, err := conn.Get(myKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, string([]byte{'\x00'}))
}

func testSetrange(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("key1"), ctx.newKey("key2")
	ok, err := conn.Set(key1, "Hello World").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Setrange(key1, 6, "Redis").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 11)
	s, err := conn.Get(key1).ToString()
	assertNil(t, err)
	assertEqual(t, s, "Hello Redis")

	i, err = conn.Setrange(key2, 6, "Redis").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 11)
	s, err = conn.Get(key2).ToString()
	assertNil(t, err)
	assertEqual(t, s, string([]byte{'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', 'R', 'e', 'd', 'i', 's'}))
}

// Lcs
const (
	rna1   = "CACCTTCCCAGGTAACAAACCAACCAACTTTCGATCTCTTGTAGATCTGTTCTCTAAACGAACTTTAAAATCTGTGTGGCTGTCACTCGGCTGCATGCTTAGTGCACTCACGCAGTATAATTAATAACTAATTACTGTCGTTGACAGGACACGAGTAACTCGTCTATCTTCTGCAGGCTGCTTACGGTTTCGTCCGTGTTGCAGCCGATCATCAGCACATCTAGGTTTCGTCCGGGTGTG"
	rna2   = "ATTAAAGGTTTATACCTTCCCAGGTAACAAACCAACCAACTTTCGATCTCTTGTAGATCTGTTCTCTAAACGAACTTTAAAATCTGTGTGGCTGTCACTCGGCTGCATGCTTAGTGCACTCACGCAGTATAATTAATAACTAATTACTGTCGTTGACAGGACACGAGTAACTCGTCTATCTTCTGCAGGCTGCTTACGGTTTCGTCCGTGTTGCAGCCGATCATCAGCACATCTAGGTTT"
	rnalcs = "ACCTTCCCAGGTAACAAACCAACCAACTTTCGATCTCTTGTAGATCTGTTCTCTAAACGAACTTTAAAATCTGTGTGGCTGTCACTCGGCTGCATGCTTAGTGCACTCACGCAGTATAATTAATAACTAATTACTGTCGTTGACAGGACACGAGTAACTCGTCTATCTTCTGCAGGCTGCTTACGGTTTCGTCCGTGTTGCAGCCGATCATCAGCACATCTAGGTTT"
)

func testStralgoLcsStrings(conn client.Conn, ctx *testCTX, t *testing.T) {
	s, err := conn.StralgoLcsStrings(rna1, rna2).ToString()
	assertNil(t, err)
	assertEqual(t, s, rnalcs)
}

func testStralgoLcsLenStrings(conn client.Conn, ctx *testCTX, t *testing.T) {
	i, err := conn.StralgoLcsLenStrings(rna1, rna2).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, len(rnalcs))
}

func testStralgoLcsIdxStrings(conn client.Conn, ctx *testCTX, t *testing.T) {
	m, err := conn.StralgoLcsIdxStrings(rna1, rna2, false, nil).ToStringValueMap()
	assertNil(t, err)

	slice, err := m["matches"].ToIntfSlice3()
	assertNil(t, err)

	assertEqual(t, slice, [][][]interface{}{
		{{int64(238), int64(238)}, {int64(239), int64(239)}},
		{{int64(236), int64(236)}, {int64(238), int64(238)}},
		{{int64(229), int64(230)}, {int64(236), int64(237)}},
		{{int64(224), int64(224)}, {int64(235), int64(235)}},
		{{int64(1), int64(222)}, {int64(13), int64(234)}},
	})
}

func testStralgoLcsKeys(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("key1"), ctx.newKey("key2")
	assertNil(t, conn.Set(key1, rna1).Err())
	assertNil(t, conn.Set(key2, rna2).Err())
	s, err := conn.StralgoLcsKeys(key1, key2).ToString()
	assertNil(t, err)
	assertEqual(t, s, rnalcs)
}

func testStralgoLcsLenKeys(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("key1"), ctx.newKey("key2")
	assertNil(t, conn.Set(key1, rna1).Err())
	assertNil(t, conn.Set(key2, rna2).Err())
	i, err := conn.StralgoLcsLenKeys(key1, key2).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, len(rnalcs))
}

func testStralgoLcsIdxKeys(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("key1"), ctx.newKey("key2")
	assertNil(t, conn.Set(key1, rna1).Err())
	assertNil(t, conn.Set(key2, rna2).Err())

	m, err := conn.StralgoLcsIdxKeys(key1, key2, false, nil).ToStringValueMap()
	assertNil(t, err)

	slice, err := m["matches"].ToIntfSlice3()
	assertNil(t, err)

	assertEqual(t, slice, [][][]interface{}{
		{{int64(238), int64(238)}, {int64(239), int64(239)}},
		{{int64(236), int64(236)}, {int64(238), int64(238)}},
		{{int64(229), int64(230)}, {int64(236), int64(237)}},
		{{int64(224), int64(224)}, {int64(235), int64(235)}},
		{{int64(1), int64(222)}, {int64(13), int64(234)}},
	})

	l, err := m["len"].ToInt64()
	assertNil(t, err)
	assertEqual(t, l, 227)

	m, err = conn.StralgoLcsIdxKeys(key1, key2, true, nil).ToStringValueMap()
	assertNil(t, err)

	tree, err := m["matches"].ToTree()
	assertNil(t, err)

	assertEqual(t, tree, []interface{}{
		[]interface{}{[]interface{}{int64(238), int64(238)}, []interface{}{int64(239), int64(239)}, int64(1)},
		[]interface{}{[]interface{}{int64(236), int64(236)}, []interface{}{int64(238), int64(238)}, int64(1)},
		[]interface{}{[]interface{}{int64(229), int64(230)}, []interface{}{int64(236), int64(237)}, int64(2)},
		[]interface{}{[]interface{}{int64(224), int64(224)}, []interface{}{int64(235), int64(235)}, int64(1)},
		[]interface{}{[]interface{}{int64(1), int64(222)}, []interface{}{int64(13), int64(234)}, int64(222)},
	})

	l, err = m["len"].ToInt64()
	assertNil(t, err)
	assertEqual(t, l, 227)

	m, err = conn.StralgoLcsIdxKeys(key1, key2, true, client.Int64Ptr(5)).ToStringValueMap()
	assertNil(t, err)

	tree, err = m["matches"].ToTree()
	assertNil(t, err)

	assertEqual(t, tree, []interface{}{
		[]interface{}{[]interface{}{int64(1), int64(222)}, []interface{}{int64(13), int64(234)}, int64(222)},
	})

	l, err = m["len"].ToInt64()
	assertNil(t, err)
	assertEqual(t, l, 227)
}

func testStrlen(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey, nonExisting := ctx.newKey("myKey"), ctx.newKey("nonExisting")
	ok, err := conn.Set(myKey, "Hello World").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Strlen(myKey).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 11)
	i, err = conn.Strlen(nonExisting).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
}

// Keys
func testDel(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2, key3 := ctx.newKey("key1"), ctx.newKey("key2"), ctx.newKey("key3")
	ok, err := conn.Set(key1, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	ok, err = conn.Set(key2, "World").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Del([]interface{}{key1, key2, key3}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
}

func testDump(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "10").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	s, err := conn.Dump(myKey).ToString()
	assertNil(t, err)
	t.Logf("%x", s)
}

func testExists(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2, nosuchkey := ctx.newKey("key1"), ctx.newKey("key2"), ctx.newKey("nosuchkey")
	ok, err := conn.Set(key1, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Exists([]interface{}{key1}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Exists([]interface{}{nosuchkey}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
	ok, err = conn.Set(key2, "World").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err = conn.Exists([]interface{}{key1, key2, nosuchkey}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
}

func testExpire(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Expire(myKey, 10).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.TTL(myKey).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 10)
	ok, err = conn.Set(myKey, "Hello World").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err = conn.TTL(myKey).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, -1)
}

func testExpireat(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Expireat(myKey, time.Now().Unix()).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Exists([]interface{}{myKey}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
}

func testKeys(conn client.Conn, ctx *testCTX, t *testing.T) {
	firstname, lastname, age := ctx.newKey("firstname"), ctx.newKey("lastname"), ctx.newKey("age")
	ok, err := conn.Mset([]client.KeyValue{{Key: firstname, Value: "Jack"}, {Key: lastname, Value: "Stuntman"}, {Key: age, Value: 35}}).ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	slice, err := conn.Keys("*name*").ToStringSlice()
	assertNil(t, err)
	// test if keys are in result
	m := make(map[string]bool, len(slice))
	for _, v := range slice {
		m[v] = true
	}
	_, ok = m[firstname]
	assertTrue(t, ok)
	_, ok = m[lastname]
	assertTrue(t, ok)
	slice, err = conn.Keys(age[:len(age)-2] + "??").ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{age})
}

func testMove(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Move(myKey, secondDB).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
}

func testObjectEncoding(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	_, err = conn.ObjectEncoding(myKey).ToString()
	assertNil(t, err)
}

func testObjectFreq(conn client.Conn, ctx *testCTX, t *testing.T) {
	// set maxmemory policy (otherwise ObjectFreq might report error)
	savePolicy := setConfig(conn, maxmemoryPolicy, "volatile-lfu", t)
	defer setConfig(conn, maxmemoryPolicy, savePolicy, t)

	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	_, err = conn.ObjectFreq(myKey).ToInt64()
	assertNil(t, err)
}

func testObjectHelp(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.ObjectHelp().ToStringSlice()
	assertNil(t, err)
}

func testObjectIdletime(conn client.Conn, ctx *testCTX, t *testing.T) {
	// set maxmemory policy (otherwise ObjectIdletime might report error)
	savePolicy := setConfig(conn, maxmemoryPolicy, "noeviction", t)
	defer setConfig(conn, maxmemoryPolicy, savePolicy, t)

	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	_, err = conn.ObjectIdletime(myKey).ToInt64()
	assertNil(t, err)
}

func testObjectRefcount(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	_, err = conn.ObjectRefcount(myKey).ToInt64()
	assertNil(t, err)
}

func testPersist(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Expire(myKey, 10).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.TTL(myKey).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 10)
	i, err = conn.Persist(myKey).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.TTL(myKey).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, -1)
}

func testPexpire(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Pexpire(myKey, 1499).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.TTL(myKey).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.PTTL(myKey).ToInt64()
	assertNil(t, err)
	t.Logf("Key expires in %d milliseconds", i)
}

func testPexpireat(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Pexpireat(myKey, (time.Now().Unix()-1)*1000+999).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.TTL(myKey).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, -2)
	i, err = conn.PTTL(myKey).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, -2)
}

func testRandomkey(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	s, err := conn.Randomkey().ToString()
	assertNil(t, err)
	t.Logf("Randomkey %s", s)
}

func testRename(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey, myOtherKey := ctx.newKey("myKey"), ctx.newKey("myOtherKey")
	ok, err := conn.Set(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	ok, err = conn.Rename(myKey, myOtherKey).ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	s, err := conn.Get(myOtherKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "Hello")
}

func testRenameNx(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey, myOtherKey := ctx.newKey("myKey"), ctx.newKey("myOtherKey")
	ok, err := conn.Set(myKey, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	ok, err = conn.Set(myOtherKey, "World").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.RenameNx(myKey, myOtherKey).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
	s, err := conn.Get(myOtherKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "World")
}

func testRestore(conn client.Conn, ctx *testCTX, t *testing.T) {
	myKey := ctx.newKey("myKey")
	ok, err := conn.Set(myKey, "10").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	s, err := conn.Dump(myKey).ToString()
	assertNil(t, err)
	err = conn.Del([]interface{}{myKey}).Err()
	assertNil(t, err)
	ok, err = conn.Restore(myKey, 0, s, false, false, nil, nil).ToBool()
	assertNil(t, err)
	assertTrue(t, ok)
	s, err = conn.Get(myKey).ToString()
	assertNil(t, err)
	assertEqual(t, s, "10")
}

func testScan(conn client.Conn, ctx *testCTX, t *testing.T) {
	slice, err := conn.Scan(0, nil, nil, nil).ToSlice()
	assertNil(t, err)
	cursor, err := slice[0].ToInt64()
	assertNil(t, err)
	for cursor != 0 {
		slice, err = conn.Scan(cursor, nil, nil, nil).ToSlice()
		assertNil(t, err)
		cursor, err = slice[0].ToInt64()
		assertNil(t, err)
	}
}

func testSort(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList := ctx.newKey("myList")
	err := conn.Rpush(myList, []interface{}{"c", "b", "a"}).Err()
	assertNil(t, err)
	slice, err := conn.Sort(myList, nil, nil, nil, nil, true, nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"a", "b", "c"})
}

func testTouch(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("key1"), ctx.newKey("key2")
	ok, err := conn.Set(key1, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	ok, err = conn.Set(key2, "World").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Touch([]interface{}{key1, key2}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
}

func testType(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2, key3 := ctx.newKey("key1"), ctx.newKey("key2"), ctx.newKey("key3")
	ok, err := conn.Set(key1, "value").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Lpush(key2, []interface{}{"value"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Sadd(key3, []interface{}{"value"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	s, err := conn.Type(key1).ToString()
	assertNil(t, err)
	assertEqual(t, s, "string")
	s, err = conn.Type(key2).ToString()
	assertNil(t, err)
	assertEqual(t, s, "list")
	s, err = conn.Type(key3).ToString()
	assertNil(t, err)
	assertEqual(t, s, "set")
}

func testUnlink(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("key1"), ctx.newKey("key2")
	ok, err := conn.Set(key1, "Hello").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	ok, err = conn.Set(key2, "World").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	i, err := conn.Unlink([]interface{}{key1, key2}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
}

func testWait(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.Wait(1, 1).ToInt64() // timeout - if no replicas, wait without timeout will wait forever
	assertNil(t, err)
}

// Lists
func testBlpop(conn client.Conn, ctx *testCTX, t *testing.T) {
	list1, list2 := ctx.newKey("list1"), ctx.newKey("list2")
	i, err := conn.Rpush(list1, []interface{}{"a", "b", "c"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	slice, err := conn.Blpop([]interface{}{list1, list2}, 0).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{list1, "a"})
}

func testBrpop(conn client.Conn, ctx *testCTX, t *testing.T) {
	list1, list2 := ctx.newKey("list1"), ctx.newKey("list2")
	i, err := conn.Rpush(list1, []interface{}{"a", "b", "c"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	slice, err := conn.Brpop([]interface{}{list1, list2}, 0).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{list1, "c"})
}

func testBrpoplpush(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList, myOtherList := ctx.newKey("myList"), ctx.newKey("myOtherList")
	i, err := conn.Rpush(myList, []interface{}{"one", "two", "three"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	s, err := conn.Brpoplpush(myList, myOtherList, 0).ToString()
	assertNil(t, err)
	assertEqual(t, s, "three")
	slice, err := conn.Lrange(myList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"one", "two"})
	slice, err = conn.Lrange(myOtherList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"three"})
}

func testLindex(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList := ctx.newKey("myList")
	i, err := conn.Lpush(myList, []interface{}{"World", "Hello"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	s, err := conn.Lindex(myList, 0).ToString()
	assertNil(t, err)
	assertEqual(t, s, "Hello")
	s, err = conn.Lindex(myList, -1).ToString()
	assertNil(t, err)
	assertEqual(t, s, "World")
	b, err := conn.Lindex(myList, 3).IsNull()
	assertNil(t, err)
	assertTrue(t, b)
}

func testLinsert(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList := ctx.newKey("myList")
	i, err := conn.Rpush(myList, []interface{}{"Hello", "World"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	i, err = conn.Linsert(myList, true, "World", "There").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	slice, err := conn.Lrange(myList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"Hello", "There", "World"})
}

func testLlen(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList := ctx.newKey("myList")
	i, err := conn.Lpush(myList, []interface{}{"World", "Hello"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	i, err = conn.Llen(myList).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
}

func testLpop(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList := ctx.newKey("myList")
	i, err := conn.Rpush(myList, []interface{}{"one", "two", "three"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	s, err := conn.Lpop(myList).ToString()
	assertNil(t, err)
	assertEqual(t, s, "one")
	slice, err := conn.Lrange(myList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"two", "three"})
}

func testLpos(conn client.Conn, ctx *testCTX, t *testing.T) {
	version := conn.ConnInfo().RedisVersion
	if version.Compare(client.ParseVersion(client.CmdLposVersion)) == -1 {
		t.Logf("lpos not available in redis version %s - expected %s", version, client.CmdLposVersion)
		return
	}
	myList := ctx.newKey("myList")
	i, err := conn.Rpush(myList, []interface{}{"a", "b", "c", "d", 1, 2, 3, 4, 3, 3, 3}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 11)
	i, err = conn.Lpos(myList, 3, nil, nil, nil).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 6)
	slice, err := conn.Lpos(myList, 3, client.Int64Ptr(2), client.Int64Ptr(0), nil).ToInt64Slice()
	assertNil(t, err)
	assertEqual(t, slice, []int64{8, 9, 10})
}

func testLpush(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList := ctx.newKey("myList")
	i, err := conn.Lpush(myList, []interface{}{"World", "Hello"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	slice, err := conn.Lrange(myList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"Hello", "World"})
}

func testLpushx(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList, myOtherList := ctx.newKey("myList"), ctx.newKey("myOtherList")
	i, err := conn.Lpush(myList, []interface{}{"World"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Lpushx(myList, []interface{}{"Hello"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	i, err = conn.Lpushx(myOtherList, []interface{}{"Hello"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
	slice, err := conn.Lrange(myList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"Hello", "World"})
	slice, err = conn.Lrange(myOtherList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{})
}

func testLrange(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList := ctx.newKey("myList")
	i, err := conn.Rpush(myList, []interface{}{"one", "two", "three"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	slice, err := conn.Lrange(myList, 0, 0).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"one"})
	slice, err = conn.Lrange(myList, -3, 2).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"one", "two", "three"})
	slice, err = conn.Lrange(myList, -100, 100).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"one", "two", "three"})
	slice, err = conn.Lrange(myList, 5, 10).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{})
}

func testLrem(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList := ctx.newKey("myList")
	i, err := conn.Rpush(myList, []interface{}{"hello", "hello", "foo", "hello"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 4)
	i, err = conn.Lrem(myList, -2, "hello").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	slice, err := conn.Lrange(myList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"hello", "foo"})
}

func testLset(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList := ctx.newKey("myList")
	i, err := conn.Rpush(myList, []interface{}{"one", "two", "three"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	ok, err := conn.Lset(myList, 0, "four").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	ok, err = conn.Lset(myList, -2, "five").ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	slice, err := conn.Lrange(myList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"four", "five", "three"})
}

func testLtrim(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList := ctx.newKey("myList")
	i, err := conn.Rpush(myList, []interface{}{"one", "two", "three"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	ok, err := conn.Ltrim(myList, 1, -1).ToBool()
	assertNil(t, err)
	assertEqual(t, ok, true)
	slice, err := conn.Lrange(myList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"two", "three"})
}

func testRpop(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList := ctx.newKey("myList")
	i, err := conn.Rpush(myList, []interface{}{"one", "two", "three"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	s, err := conn.Rpop(myList).ToString()
	assertNil(t, err)
	assertEqual(t, s, "three")
	slice, err := conn.Lrange(myList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"one", "two"})
}

func testRpoplpush(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList, myOtherList := ctx.newKey("myList"), ctx.newKey("myOtherList")
	i, err := conn.Rpush(myList, []interface{}{"one", "two", "three"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	s, err := conn.Rpoplpush(myList, myOtherList).ToString()
	assertNil(t, err)
	assertEqual(t, s, "three")
	slice, err := conn.Lrange(myList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"one", "two"})
	slice, err = conn.Lrange(myOtherList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"three"})
}

func testRpush(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList := ctx.newKey("myList")
	i, err := conn.Rpush(myList, []interface{}{"Hello", "World"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	slice, err := conn.Lrange(myList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"Hello", "World"})
}

func testRpushx(conn client.Conn, ctx *testCTX, t *testing.T) {
	myList, myOtherList := ctx.newKey("myList"), ctx.newKey("myOtherList")
	i, err := conn.Rpush(myList, []interface{}{"Hello"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Rpushx(myList, []interface{}{"World"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	i, err = conn.Lpushx(myOtherList, []interface{}{"World"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
	slice, err := conn.Lrange(myList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"Hello", "World"})
	slice, err = conn.Lrange(myOtherList, 0, -1).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{})
}

// Hashes
func testHdel(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myHash")
	assertNil(t, conn.Hset(key, []client.FieldValue{{Field: "field1", Value: "foo"}}).Err())
	i, err := conn.Hdel(key, []interface{}{"field1"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Hdel(key, []interface{}{"field2"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
}

func testHexists(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myHash")
	assertNil(t, conn.Hset(key, []client.FieldValue{{Field: "field1", Value: "foo"}}).Err())
	i, err := conn.Hexists(key, "field1").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Hexists(key, "field2").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
}

func testHget(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myHash")
	assertNil(t, conn.Hset(key, []client.FieldValue{{Field: "field1", Value: "foo"}}).Err())
	s, err := conn.Hget(key, "field1").ToString()
	assertNil(t, err)
	assertEqual(t, s, "foo")
	b, err := conn.Hget(key, "field2").IsNull()
	assertNil(t, err)
	assertTrue(t, b)
}

func testHgetall(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myHash")
	assertNil(t, conn.Hset(key, []client.FieldValue{{Field: "field1", Value: "Hello"}, {Field: "field2", Value: "World"}}).Err())
	m, err := conn.Hgetall(key).ToStringStringMap()
	assertNil(t, err)
	assertEqual(t, m, map[string]string{"field1": "Hello", "field2": "World"})
}

func testHincrby(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myHash")
	i, err := conn.Hset(key, []client.FieldValue{{Field: "field", Value: 5}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Hincrby(key, "field", 1).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 6)
	i, err = conn.Hincrby(key, "field", -1).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 5)
	i, err = conn.Hincrby(key, "field", -10).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, -5)
}

func testHincrbyfloat(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myHash")
	i, err := conn.Hset(key, []client.FieldValue{{Field: "field", Value: 10.50}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	f, err := conn.Hincrbyfloat(key, "field", 0.1).ToFloat64()
	assertNil(t, err)
	assertEqual(t, f, 10.6)
	f, err = conn.Hincrbyfloat(key, "field", -5).ToFloat64()
	assertNil(t, err)
	assertEqual(t, f, 5.6)
	i, err = conn.Hset(key, []client.FieldValue{{Field: "field", Value: "5.0e3"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
	f, err = conn.Hincrbyfloat(key, "field", 2.0e2).ToFloat64()
	assertNil(t, err)
	assertEqual(t, f, 5200.0)
}

func testHkeys(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myHash")
	assertNil(t, conn.Hset(key, []client.FieldValue{{Field: "field1", Value: "Hello"}, {Field: "field2", Value: "World"}}).Err())
	slice, err := conn.Hkeys(key).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"field1", "field2"})
}

func testHlen(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myHash")
	assertNil(t, conn.Hset(key, []client.FieldValue{{Field: "field1", Value: "Hello"}, {Field: "field2", Value: "World"}}).Err())
	i, err := conn.Hlen(key).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
}

func testHmget(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myHash")
	assertNil(t, conn.Hset(key, []client.FieldValue{{Field: "field1", Value: "Hello"}, {Field: "field2", Value: "World"}}).Err())
	a2, err := conn.Hmget(key, []interface{}{"field1", "field2", "nofield"}).ToIntfSlice()
	assertNil(t, err)
	assertEqual(t, a2, []interface{}{"Hello", "World", nil})
}

func testHset(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myHash")
	assertNil(t, conn.Hset(key, []client.FieldValue{{Field: "field1", Value: "Hello"}}).Err())
	s, err := conn.Hget(key, "field1").ToString()
	assertNil(t, err)
	assertEqual(t, s, "Hello")
}

func testHsetNx(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myHash")
	i, err := conn.HsetNx(key, "field3", "Hello").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.HsetNx(key, "field3", "World").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
	s, err := conn.Hget(key, "field3").ToString()
	assertNil(t, err)
	assertEqual(t, s, "Hello")
}

func testHstrlen(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myHash")
	i, err := conn.Hset(key, []client.FieldValue{{Field: "f1", Value: "HelloWorld"}, {Field: "f2", Value: 99}, {Field: "f3", Value: -256}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Hstrlen(key, "f1").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 10)
	i, err = conn.Hstrlen(key, "f2").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	i, err = conn.Hstrlen(key, "f3").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 4)
}

func testHvals(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myHash")
	assertNil(t, conn.Hset(key, []client.FieldValue{{Field: "field1", Value: "Hello"}, {Field: "field2", Value: "World"}}).Err())
	slice, err := conn.Hvals(key).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"Hello", "World"})
}

func testHscan(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myHash")
	i, err := conn.Hset(key, []client.FieldValue{{Field: "field1", Value: "Hello"}, {Field: "field2", Value: "World"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	slice, err := conn.Hscan(key, 0, nil, nil).ToSlice()
	assertNil(t, err)
	cursor, err := slice[0].ToInt64()
	assertNil(t, err)
	for cursor != 0 {
		slice, err = conn.Hscan(key, cursor, nil, nil).ToSlice()
		assertNil(t, err)
		cursor, err = slice[0].ToInt64()
		assertNil(t, err)
	}
}

// Pubsub
func msgCallback(ch chan<- string) client.MsgCallback {
	return func(pattern, channel, msg string) {
		ch <- msg
	}
}

func testPubsub(conn client.Conn, ctx *testCTX, t *testing.T) {
	// provoking error - subscribe with no channel
	err := conn.Subscribe([]string{}, nil).Err()
	assertNotNil(t, err)

	channel := client.RandomKey("")
	// unsubscribe from not subscribed channel
	assertNil(t, conn.Unsubscribe([]string{channel}).Err())

	// check if all out of band messages are consumed
	ch := make(chan string, 1)

	assertNil(t, conn.Subscribe([]string{"chan1", "chan2", "chan1", channel}, msgCallback(ch)).Err())
	noOfClients, err := conn.Publish(channel, "mymessage").ToInt64()
	assertNil(t, err)
	assertEqual(t, noOfClients, 1)
	r := <-ch
	assertEqual(t, r, "mymessage")

	// try to punsubscribe channels subscribed with subscribe
	assertNil(t, conn.Punsubscribe([]string{"chan1"}).Err())

	// try to subscribe non-pattern channel with psubscribe
	assertNil(t, conn.Psubscribe([]string{"chan3*"}, nil).Err())
	channels, err := conn.PubsubChannels(nil).ToStringSlice()
	assertNil(t, err)

	t.Logf("channels: %v", channels)

	channels, err = conn.PubsubChannels(client.StringPtr("*")).ToStringSlice()
	assertNil(t, err)

	t.Logf("channels: %v", channels)

	slice, err := conn.PubsubNumsub([]string{channel}).ToStringSlice()
	assertNil(t, err)
	t.Logf("number of subscribers: %v", slice)

	i, err := conn.PubsubNumpat().ToInt64()
	assertNil(t, err)
	t.Logf("number of patterns: %d", i)

}

// Sets
func testSadd(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("mySet")
	i, err := conn.Sadd(key, []interface{}{"Hello"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Sadd(key, []interface{}{"World"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Sadd(key, []interface{}{"World"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
}

func testScard(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("mySet")
	i, err := conn.Sadd(key, []interface{}{"Hello", "World"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	i, err = conn.Scard(key).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
}

func testSdiff(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("key1"), ctx.newKey("key2")
	i, err := conn.Sadd(key1, []interface{}{"a", "b", "c"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Sadd(key2, []interface{}{"c", "d", "e"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	set, err := conn.Sdiff([]interface{}{key1, key2}).ToStringSet()
	assertNil(t, err)
	assertEqual(t, set, map[string]bool{"b": true, "a": true})
}

func testSdiffstore(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("key1"), ctx.newKey("key2")
	i, err := conn.Sadd(key1, []interface{}{"a", "b", "c"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Sadd(key2, []interface{}{"c", "d", "e"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	key := ctx.newKey("key")
	i, err = conn.Sdiffstore(key, []interface{}{key1, key2}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	set, err := conn.Smembers(key).ToStringSet()
	assertNil(t, err)
	assertEqual(t, set, map[string]bool{"b": true, "a": true})
}

func testSinter(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("key1"), ctx.newKey("key2")
	i, err := conn.Sadd(key1, []interface{}{"a", "b", "c"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Sadd(key2, []interface{}{"c", "d", "e"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	set, err := conn.Sinter([]interface{}{key1, key2}).ToStringSet()
	assertNil(t, err)
	assertEqual(t, set, map[string]bool{"c": true})
}

func testSinterstore(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("key1"), ctx.newKey("key2")
	i, err := conn.Sadd(key1, []interface{}{"a", "b", "c"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Sadd(key2, []interface{}{"c", "d", "e"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	key := ctx.newKey("key")
	i, err = conn.Sinterstore(key, []interface{}{key1, key2}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	set, err := conn.Smembers(key).ToStringSet()
	assertNil(t, err)
	assertEqual(t, set, map[string]bool{"c": true})
}

func testSismember(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("mySet")
	i, err := conn.Sadd(key, []interface{}{"one"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Sismember(key, "one").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Sismember(key, "two").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
}

func testSmembers(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("mySet")
	i, err := conn.Sadd(key, []interface{}{"Hello", "World"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	set, err := conn.Smembers(key).ToStringSet()
	assertNil(t, err)
	assertEqual(t, set, map[string]bool{"Hello": true, "World": true})
}

func testSmove(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("mySet")
	i, err := conn.Sadd(key, []interface{}{"one", "two"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	otherKey := ctx.newKey("myotherset")
	i, err = conn.Sadd(otherKey, []interface{}{"three"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Smove(key, otherKey, "two").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	set, err := conn.Smembers(key).ToStringSet()
	assertNil(t, err)
	assertEqual(t, set, map[string]bool{"one": true})
	set, err = conn.Smembers(otherKey).ToStringSet()
	assertNil(t, err)
	assertEqual(t, set, map[string]bool{"two": true, "three": true})
}

func testSpop(conn client.Conn, ctx *testCTX, t *testing.T) {
	values := map[string]bool{"one": true, "two": true, "three": true}

	key := ctx.newKey("mySet")
	i, err := conn.Sadd(key, []interface{}{"one", "two", "three"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	s, err := conn.Spop(key, nil).ToString() // pop random item
	assertNil(t, err)
	delete(values, s)
	set, err := conn.Smembers(key).ToStringSet()
	assertNil(t, err)
	assertEqual(t, set, values)
	i, err = conn.Sadd(key, []interface{}{"four", "five"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	values["four"], values["five"] = true, true
	set, err = conn.Spop(key, client.Int64Ptr(3)).ToStringSet()
	assertNil(t, err)
	for k := range set {
		delete(values, k)
	}
	set, err = conn.Smembers(key).ToStringSet()
	assertNil(t, err)
	assertEqual(t, set, values)
}

func testSrandmember(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("mySet")
	assertNil(t, conn.Srandmember(key, nil).Err())
	assertNil(t, conn.Srandmember(key, client.Int64Ptr(2)).Err())
	assertNil(t, conn.Srandmember(key, client.Int64Ptr(-5)).Err())
}

func testSrem(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("mySet")
	i, err := conn.Sadd(key, []interface{}{"one", "two", "three"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Srem(key, []interface{}{"one"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Srem(key, []interface{}{"four"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
	set, err := conn.Smembers(key).ToStringSet()
	assertNil(t, err)
	assertEqual(t, set, map[string]bool{"two": true, "three": true})
}

func testSunion(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("key1"), ctx.newKey("key2")
	i, err := conn.Sadd(key1, []interface{}{"a", "b", "c"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Sadd(key2, []interface{}{"c", "d", "e"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	set, err := conn.Sunion([]interface{}{key1, key2}).ToStringSet()
	assertNil(t, err)
	assertEqual(t, set, map[string]bool{"a": true, "b": true, "c": true, "d": true, "e": true})
}

func testSunionstore(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("key1"), ctx.newKey("key2")
	i, err := conn.Sadd(key1, []interface{}{"a", "b", "c"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Sadd(key2, []interface{}{"c", "d", "e"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	key := ctx.newKey("key")
	i, err = conn.Sunionstore(key, []interface{}{key1, key2}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 5)
	set, err := conn.Smembers(key).ToStringSet()
	assertNil(t, err)
	assertEqual(t, set, map[string]bool{"a": true, "b": true, "c": true, "d": true, "e": true})
}

func testSscan(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("mySet")
	i, err := conn.Sadd(key, []interface{}{"Hello", "World"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	slice, err := conn.Sscan(key, 0, nil, nil).ToSlice()
	assertNil(t, err)
	cursor, err := slice[0].ToInt64()
	assertNil(t, err)
	for cursor != 0 {
		slice, err = conn.Sscan(key, cursor, nil, nil).ToSlice()
		assertNil(t, err)
		cursor, err = slice[0].ToInt64()
		assertNil(t, err)
	}
}

// Sorted Sets
func testBzpopmax(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("myZset1"), ctx.newKey("myZset2")
	i, err := conn.Zadd(key1, []client.ScoreMember{{Score: 0, Member: "a"}, {Score: 1, Member: "b"}, {Score: 2, Member: "c"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	slice, err := conn.Bzpopmax([]interface{}{key1, key2}, 0).ToIntfSlice()
	assertNil(t, err)
	assertEqual(t, slice, []interface{}{key1, "c", float64(2)})
}

func testBzpopmin(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("myZset1"), ctx.newKey("myZset2")
	i, err := conn.Zadd(key1, []client.ScoreMember{{Score: 0, Member: "a"}, {Score: 1, Member: "b"}, {Score: 2, Member: "c"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	slice, err := conn.Bzpopmin([]interface{}{key1, key2}, 0).ToIntfSlice()
	assertNil(t, err)
	assertEqual(t, slice, []interface{}{key1, "a", float64(0)})
}

func testZadd(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")

	// Zadd
	err := conn.Del([]interface{}{key}).Err()
	assertNil(t, err)
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "uno"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Zadd(key, []client.ScoreMember{{Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	slice, err := conn.Zrange(key, 0, -1, true).ToIntfSlice2()
	assertNil(t, err)
	assertEqual(t, slice, [][]interface{}{{"one", float64(1)}, {"uno", float64(1)}, {"two", float64(2)}, {"three", float64(3)}})

	// ZaddCh
	err = conn.Del([]interface{}{key}).Err()
	assertNil(t, err)
	i, err = conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	i, err = conn.ZaddCh(key, []client.ScoreMember{{Score: 2.2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)

	// ZaddNx
	err = conn.Del([]interface{}{key}).Err()
	assertNil(t, err)
	i, err = conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	i, err = conn.ZaddNx(key, []client.ScoreMember{{Score: 1.1, Member: "one"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	slice, err = conn.Zrange(key, 0, -1, true).ToIntfSlice2()
	assertNil(t, err)
	assertEqual(t, slice, [][]interface{}{{"one", float64(1)}, {"two", float64(2)}, {"three", float64(3)}})

	// ZaddXx
	err = conn.Del([]interface{}{key}).Err()
	assertNil(t, err)
	i, err = conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	i, err = conn.ZaddXx(key, []client.ScoreMember{{Score: 1.1, Member: "one"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
	slice, err = conn.Zrange(key, 0, -1, true).ToIntfSlice2()
	assertNil(t, err)
	assertEqual(t, slice, [][]interface{}{{"one", float64(1.1)}, {"two", float64(2)}})

	// ZaddXxCh
	err = conn.Del([]interface{}{key}).Err()
	assertNil(t, err)
	i, err = conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	i, err = conn.ZaddXxCh(key, []client.ScoreMember{{Score: 1.1, Member: "one"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	slice, err = conn.Zrange(key, 0, -1, true).ToIntfSlice2()
	assertNil(t, err)
	assertEqual(t, slice, [][]interface{}{{"one", float64(1.1)}, {"two", float64(2)}})
}

func testZcard(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	i, err = conn.Zcard(key).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
}

func testZcount(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Zcount(key, client.InfNeg, client.InfPos).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Zcount(key, client.Zopen(1), client.Zclose(3)).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
}

func testZincrby(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	f, err := conn.Zincrby(key, 2, "one").ToFloat64()
	assertNil(t, err)
	assertEqual(t, f, 3.0)
	slice, err := conn.Zrange(key, 0, -1, true).ToIntfSlice2()
	assertNil(t, err)
	assertEqual(t, slice, [][]interface{}{{"two", float64(2)}, {"one", float64(3)}})
}

func testZinterstore(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("zset1"), ctx.newKey("zset2")
	i, err := conn.Zadd(key1, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	i, err = conn.Zadd(key2, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	key := ctx.newKey("out")
	i, err = conn.Zinterstore(key, 2, []interface{}{key1, key2}, []int64{2, 3}, nil).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	slice, err := conn.Zrange(key, 0, -1, true).ToIntfSlice2()
	assertNil(t, err)
	assertEqual(t, slice, [][]interface{}{{"one", float64(5)}, {"two", float64(10)}})
}

func testZlexcount(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 0, Member: "a"}, {Score: 0, Member: "b"}, {Score: 0, Member: "c"}, {Score: 0, Member: "d"}, {Score: 0, Member: "e"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 5)
	i, err = conn.Zadd(key, []client.ScoreMember{{Score: 0, Member: "f"}, {Score: 0, Member: "g"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	i, err = conn.Zlexcount(key, "-", "+").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 7)
	i, err = conn.Zlexcount(key, "[b", "[f").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 5)
}

func testZpopmax(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	slice, err := conn.Zpopmax(key, nil).ToIntfSlice()
	assertNil(t, err)
	assertEqual(t, slice, []interface{}{"three", float64(3)})
}

func testZpopmin(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	slice, err := conn.Zpopmin(key, nil).ToIntfSlice()
	assertNil(t, err)
	assertEqual(t, slice, []interface{}{"one", float64(1)})
}

func testZrange(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	slice, err := conn.Zrange(key, 0, -1, false).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"one", "two", "three"})
	slice, err = conn.Zrange(key, 2, 3, false).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"three"})
	slice, err = conn.Zrange(key, -2, -1, false).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"two", "three"})
	slice2, err := conn.Zrange(key, 0, 1, true).ToIntfSlice2()
	assertNil(t, err)
	assertEqual(t, slice2, [][]interface{}{{"one", float64(1)}, {"two", float64(2)}})
}

func testZrangebylex(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 0, Member: "a"}, {Score: 0, Member: "b"}, {Score: 0, Member: "c"}, {Score: 0, Member: "d"}, {Score: 0, Member: "e"}, {Score: 0, Member: "f"}, {Score: 0, Member: "g"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 7)
	slice, err := conn.Zrangebylex(key, "-", "[c", nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"a", "b", "c"})
	slice, err = conn.Zrangebylex(key, "-", "(c", nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"a", "b"})
	slice, err = conn.Zrangebylex(key, "[aaa", "(g", nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"b", "c", "d", "e", "f"})
}

func testZrangebyscore(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	slice, err := conn.Zrangebyscore(key, client.InfNeg, client.InfPos, false, nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"one", "two", "three"})
	slice, err = conn.Zrangebyscore(key, client.Zclose(1), client.Zclose(2), false, nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"one", "two"})
	slice, err = conn.Zrangebyscore(key, client.Zopen(1), client.Zclose(2), false, nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"two"})
	slice, err = conn.Zrangebyscore(key, client.Zopen(1), client.Zopen(2), false, nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{})
}

func testZrank(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Zrank(key, "three").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	b, err := conn.Zrank(key, "four").IsNull()
	assertNil(t, err)
	assertTrue(t, b)
}

func testZrem(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Zrem(key, []interface{}{"two"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	slice, err := conn.Zrange(key, 0, -1, true).ToIntfSlice2()
	assertNil(t, err)
	assertEqual(t, slice, [][]interface{}{{"one", float64(1)}, {"three", float64(3)}})
}

func testZremrangebylex(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 0, Member: "aaaa"}, {Score: 0, Member: "b"}, {Score: 0, Member: "c"}, {Score: 0, Member: "d"}, {Score: 0, Member: "e"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 5)
	i, err = conn.Zadd(key, []client.ScoreMember{{Score: 0, Member: "foo"}, {Score: 0, Member: "zap"}, {Score: 0, Member: "zip"}, {Score: 0, Member: "ALPHA"}, {Score: 0, Member: "alpha"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 5)
	slice, err := conn.Zrange(key, 0, -1, false).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"ALPHA", "aaaa", "alpha", "b", "c", "d", "e", "foo", "zap", "zip"})
	i, err = conn.Zremrangebylex(key, "[alpha", "[omega").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 6)
	slice, err = conn.Zrange(key, 0, -1, false).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"ALPHA", "aaaa", "zap", "zip"})
}

func testZremrangebyrank(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Zremrangebyrank(key, 0, 1).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	slice, err := conn.Zrange(key, 0, -1, true).ToIntfSlice2()
	assertNil(t, err)
	assertEqual(t, slice, [][]interface{}{{"three", float64(3)}})
}

func testZremrangebyscore(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Zremrangebyscore(key, client.InfNeg, client.Zopen(2)).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	slice, err := conn.Zrange(key, 0, -1, true).ToIntfSlice2()
	assertNil(t, err)
	assertEqual(t, slice, [][]interface{}{{"two", float64(2)}, {"three", float64(3)}})
}

func testZrevrange(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	slice, err := conn.Zrevrange(key, 0, -1, false).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"three", "two", "one"})
	slice, err = conn.Zrevrange(key, 2, 3, false).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"one"})
	slice, err = conn.Zrevrange(key, -2, -1, false).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"two", "one"})
}

func testZrevrangebylex(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 0, Member: "a"}, {Score: 0, Member: "b"}, {Score: 0, Member: "c"}, {Score: 0, Member: "d"}, {Score: 0, Member: "e"}, {Score: 0, Member: "f"}, {Score: 0, Member: "g"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 7)
	slice, err := conn.Zrevrangebylex(key, "[c", "-", nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"c", "b", "a"})
	slice, err = conn.Zrevrangebylex(key, "(c", "-", nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"b", "a"})
	slice, err = conn.Zrevrangebylex(key, "(g", "[aaa", nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"f", "e", "d", "c", "b"})
}

func testZrevrangebyscore(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	slice, err := conn.Zrevrangebyscore(key, client.InfPos, client.InfNeg, false, nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"three", "two", "one"})
	slice, err = conn.Zrevrangebyscore(key, client.Zclose(2), client.Zclose(1), false, nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"two", "one"})
	slice, err = conn.Zrevrangebyscore(key, client.Zclose(2), client.Zopen(1), false, nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"two"})
	slice, err = conn.Zrevrangebyscore(key, client.Zopen(2), client.Zopen(1), false, nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{})
}

func testZrevrank(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Zrevrank(key, "one").ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	b, err := conn.Zrank(key, "four").IsNull()
	assertNil(t, err)
	assertTrue(t, b)
}

func testZscan(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	slice, err := conn.Zscan(key, 0, nil, nil).ToSlice()
	assertNil(t, err)
	cursor, err := slice[0].ToInt64()
	assertNil(t, err)
	for cursor != 0 {
		slice, err = conn.Zscan(key, cursor, nil, nil).ToSlice()
		assertNil(t, err)
		cursor, err = slice[0].ToInt64()
		assertNil(t, err)
	}
}

func testZscore(conn client.Conn, ctx *testCTX, t *testing.T) {
	key := ctx.newKey("myZset")
	i, err := conn.Zadd(key, []client.ScoreMember{{Score: 1, Member: "one"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	f, err := conn.Zscore(key, "one").ToFloat64()
	assertNil(t, err)
	assertEqual(t, f, 1.0)
}

func testZunionstore(conn client.Conn, ctx *testCTX, t *testing.T) {
	key1, key2 := ctx.newKey("zset1"), ctx.newKey("zset2")
	i, err := conn.Zadd(key1, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	i, err = conn.Zadd(key2, []client.ScoreMember{{Score: 1, Member: "one"}, {Score: 2, Member: "two"}, {Score: 3, Member: "three"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	key := ctx.newKey("out")
	i, err = conn.Zunionstore(key, 2, []interface{}{key1, key2}, []int64{2, 3}, nil).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	slice, err := conn.Zrange(key, 0, -1, true).ToIntfSlice2()
	assertNil(t, err)
	assertEqual(t, slice, [][]interface{}{{"one", float64(5)}, {"three", float64(9)}, {"two", float64(10)}})
}

// Geo
func testGeoadd(conn client.Conn, ctx *testCTX, t *testing.T) {
	sicily := ctx.newKey("Sicily")
	i, err := conn.Geoadd(sicily, []client.LongitudeLatitudeMember{{Longitude: 13.361389, Latitude: 38.115556, Member: "Palermo"}, {Longitude: 15.087269, Latitude: 37.502669, Member: "Catania"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	f, err := conn.Geodist(sicily, "Palermo", "Catania", nil).ToFloat64()
	assertNil(t, err)
	assertEqual(t, f, 166274.1516)
	slice, err := conn.Georadius(sicily, 15, 37, 100, client.UnitKm, false, false, false, nil, nil, nil, nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"Catania"})
	slice, err = conn.Georadius(sicily, 15, 37, 200, client.UnitKm, false, false, false, nil, nil, nil, nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"Palermo", "Catania"})
}

func testGeodist(conn client.Conn, ctx *testCTX, t *testing.T) {
	sicily := ctx.newKey("Sicily")
	i, err := conn.Geoadd(sicily, []client.LongitudeLatitudeMember{{Longitude: 13.361389, Latitude: 38.115556, Member: "Palermo"}, {Longitude: 15.087269, Latitude: 37.502669, Member: "Catania"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	f, err := conn.Geodist(sicily, "Palermo", "Catania", nil).ToFloat64()
	assertNil(t, err)
	assertEqual(t, f, 166274.1516)
	km, mi := client.UnitKm, client.UnitMi
	f, err = conn.Geodist(sicily, "Palermo", "Catania", &km).ToFloat64()
	assertNil(t, err)
	assertEqual(t, f, 166.2742)
	f, err = conn.Geodist(sicily, "Palermo", "Catania", &mi).ToFloat64()
	assertNil(t, err)
	assertEqual(t, f, 103.3182)
	b, err := conn.Geodist(sicily, "Foo", "Bar", nil).IsNull()
	assertNil(t, err)
	assertTrue(t, b)
}

func testGeohash(conn client.Conn, ctx *testCTX, t *testing.T) {
	sicily := ctx.newKey("Sicily")
	i, err := conn.Geoadd(sicily, []client.LongitudeLatitudeMember{{Longitude: 13.361389, Latitude: 38.115556, Member: "Palermo"}, {Longitude: 15.087269, Latitude: 37.502669, Member: "Catania"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	slice, err := conn.Geohash(sicily, []interface{}{"Palermo", "Catania"}).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"sqc8b49rny0", "sqdtr74hyu0"})
}

func testGeopos(conn client.Conn, ctx *testCTX, t *testing.T) {
	sicily := ctx.newKey("Sicily")
	i, err := conn.Geoadd(sicily, []client.LongitudeLatitudeMember{{Longitude: 13.361389, Latitude: 38.115556, Member: "Palermo"}, {Longitude: 15.087269, Latitude: 37.502669, Member: "Catania"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	slice, err := conn.Geopos(sicily, []interface{}{"Palermo", "Catania", "NonExisting"}).ToIntfSlice2()
	assertNil(t, err)
	assertEqual(t, slice, [][]interface{}{{13.36138933897018433, 38.11555639549629859}, {15.08726745843887329, 37.50266842333162032}, {}})
}

func testGeoradius(conn client.Conn, ctx *testCTX, t *testing.T) {
	sicily := ctx.newKey("Sicily")
	i, err := conn.Geoadd(sicily, []client.LongitudeLatitudeMember{{Longitude: 13.361389, Latitude: 38.115556, Member: "Palermo"}, {Longitude: 15.087269, Latitude: 37.502669, Member: "Catania"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	slice, err := conn.Georadius(sicily, 15, 37, 200, client.UnitKm, false, true, false, nil, nil, nil, nil).ToIntfSlice2()
	assertNil(t, err)
	assertEqual(t, slice, [][]interface{}{{"Palermo", "190.4424"}, {"Catania", "56.4413"}})
	tree, err := conn.Georadius(sicily, 15, 37, 200, client.UnitKm, true, false, false, nil, nil, nil, nil).ToTree()
	assertNil(t, err)
	assertEqual(t, tree, []interface{}{
		[]interface{}{"Palermo", []interface{}{13.36138933897018433, 38.11555639549629859}},
		[]interface{}{"Catania", []interface{}{15.08726745843887329, 37.50266842333162032}},
	})
	tree, err = conn.Georadius(sicily, 15, 37, 200, client.UnitKm, true, true, false, nil, nil, nil, nil).ToTree()
	assertNil(t, err)
	assertEqual(t, tree, []interface{}{
		[]interface{}{"Palermo", "190.4424", []interface{}{13.36138933897018433, 38.11555639549629859}},
		[]interface{}{"Catania", "56.4413", []interface{}{15.08726745843887329, 37.50266842333162032}},
	})
}

func testGeoradiusbymember(conn client.Conn, ctx *testCTX, t *testing.T) {
	sicily := ctx.newKey("Sicily")
	i, err := conn.Geoadd(sicily, []client.LongitudeLatitudeMember{{Longitude: 13.583333, Latitude: 37.316667, Member: "Agrigento"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Geoadd(sicily, []client.LongitudeLatitudeMember{{Longitude: 13.361389, Latitude: 38.115556, Member: "Palermo"}, {Longitude: 15.087269, Latitude: 37.502669, Member: "Catania"}}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	slice, err := conn.Georadiusbymember(sicily, "Agrigento", 100, client.UnitKm, false, false, false, nil, nil, nil, nil).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"Agrigento", "Palermo"})
}

// HyperLogLog
func testPfadd(conn client.Conn, ctx *testCTX, t *testing.T) {
	hll := ctx.newKey("hll")
	i, err := conn.Pfadd(hll, []interface{}{"a", "b", "c", "d", "e", "f", "g"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Pfcount([]interface{}{hll}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 7)
}

func testPfcount(conn client.Conn, ctx *testCTX, t *testing.T) {
	hll, someOtherHll := ctx.newKey("hll"), ctx.newKey("someOtherHll")
	i, err := conn.Pfadd(hll, []interface{}{"foo", "bar", "zap"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Pfadd(hll, []interface{}{"zap", "zap", "zap"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
	i, err = conn.Pfadd(hll, []interface{}{"foo", "bar"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 0)
	i, err = conn.Pfcount([]interface{}{hll}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
	i, err = conn.Pfadd(someOtherHll, []interface{}{"1", "2", "3"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Pfcount([]interface{}{hll, someOtherHll}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 6)
}

func testPfmerge(conn client.Conn, ctx *testCTX, t *testing.T) {
	hll1, hll2, hll3 := ctx.newKey("hll1"), ctx.newKey("hll2"), ctx.newKey("hll3")
	i, err := conn.Pfadd(hll1, []interface{}{"foo", "bar", "zap", "a"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	i, err = conn.Pfadd(hll2, []interface{}{"a", "b", "c", "foo"}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	ok, err := conn.Pfmerge(hll3, []interface{}{hll1, hll2}).ToBool()
	assertNil(t, err)
	assertTrue(t, ok)
	i, err = conn.Pfcount([]interface{}{hll3}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 6)
}

const testScript = "return {KEYS[1],KEYS[2],ARGV[1],ARGV[2]}"

// Scripting
func testEval(conn client.Conn, ctx *testCTX, t *testing.T) {
	slice, err := conn.Eval(testScript, 2, []interface{}{"key1", "key2"}, []interface{}{"first", "second"}).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"key1", "key2", "first", "second"})

}

func testEvalsha(conn client.Conn, ctx *testCTX, t *testing.T) {
	sha, err := conn.ScriptLoad(testScript).ToString()
	assertNil(t, err)
	slice, err := conn.Evalsha(sha, 2, []interface{}{"key1", "key2"}, []interface{}{"first", "second"}).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, slice, []string{"key1", "key2", "first", "second"})
}

func testScriptExists(conn client.Conn, ctx *testCTX, t *testing.T) {
	shaNoExists := ctx.newKey("")
	sha, err := conn.ScriptLoad(testScript).ToString()
	assertNil(t, err)
	slice, err := conn.ScriptExists([]string{sha, shaNoExists}).ToInt64Slice()
	assertNil(t, err)
	assertEqual(t, slice, []int64{1, 0})
}

func testScriptLoad(conn client.Conn, ctx *testCTX, t *testing.T) {
	err := conn.ScriptLoad(testScript).Err()
	assertNil(t, err)
}

// Streams
func testXadd(conn client.Conn, ctx *testCTX, t *testing.T) {
	myStream := ctx.newKey("myStream")
	id1, err := conn.Xadd(myStream, "*", []client.FieldValue{{Field: "name", Value: "Sara"}, {Field: "surname", Value: "OConnor"}}).ToString()
	assertNil(t, err)
	id2, err := conn.Xadd(myStream, "*", []client.FieldValue{{Field: "field1", Value: "value1"}, {Field: "field2", Value: "value2"}, {Field: "field3", Value: "value3"}}).ToString()
	assertNil(t, err)
	i, err := conn.Xlen(myStream).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	slice, err := conn.Xrange(myStream, "-", "+", nil).ToXrange()
	assertNil(t, err)
	assertEqual(t, slice, []client.XItem{
		{ID: id1, Items: []string{"name", "Sara", "surname", "OConnor"}},
		{ID: id2, Items: []string{"field1", "value1", "field2", "value2", "field3", "value3"}},
	})
}

func testXdel(conn client.Conn, ctx *testCTX, t *testing.T) {
	myStream := ctx.newKey("myStream")
	id1, err := conn.Xadd(myStream, "*", []client.FieldValue{{Field: "a", Value: "1"}}).ToString()
	assertNil(t, err)
	id2, err := conn.Xadd(myStream, "*", []client.FieldValue{{Field: "b", Value: "2"}}).ToString()
	assertNil(t, err)
	id3, err := conn.Xadd(myStream, "*", []client.FieldValue{{Field: "c", Value: "3"}}).ToString()
	assertNil(t, err)
	i, err := conn.Xdel(myStream, []string{id2}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)
	slice, err := conn.Xrange(myStream, "-", "+", nil).ToXrange()
	assertNil(t, err)
	assertEqual(t, slice, []client.XItem{{ID: id1, Items: []string{"a", "1"}}, {ID: id3, Items: []string{"c", "3"}}})
}

func testXgroupCreate(conn client.Conn, ctx *testCTX, t *testing.T) {
	myStream := ctx.newKey("myStream")
	myGroup := ctx.newKey("myGroup")
	consumer1, consumer2 := ctx.newKey("consumer1"), ctx.newKey("consumer2")

	ok, err := conn.XgroupCreate(myStream, myGroup, "$", true).ToBool()
	assertNil(t, err)
	assertTrue(t, ok)

	id1, err := conn.Xadd(myStream, "*", []client.FieldValue{{Field: "a", Value: "1"}}).ToString()
	assertNil(t, err)
	id2, err := conn.Xadd(myStream, "*", []client.FieldValue{{Field: "b", Value: "2"}}).ToString()
	assertNil(t, err)
	err = conn.Xadd(myStream, "*", []client.FieldValue{{Field: "c", Value: "3"}}).Err()
	assertNil(t, err)

	err = conn.Xreadgroup(client.GroupConsumer{Group: myGroup, Consumer: consumer1}, client.Int64Ptr(1), nil, false, []interface{}{myStream}, []string{">"}).Err()
	assertNil(t, err)
	err = conn.Xreadgroup(client.GroupConsumer{Group: myGroup, Consumer: consumer2}, client.Int64Ptr(1), nil, false, []interface{}{myStream}, []string{">"}).Err()
	assertNil(t, err)

	m, err := conn.XinfoGroups(myStream).ToStringMapSlice()
	assertNil(t, err)
	assertEqual(t, m, []map[string]interface{}{{"name": myGroup, "consumers": int64(2), "pending": int64(2), "last-delivered-id": id2}})

	slice, err := conn.XinfoConsumers(myStream, myGroup).ToStringMapSlice()
	assertNil(t, err)
	assertEqual(t, len(slice), 2)
	// delete idle time to easier comparison
	for _, m := range slice {
		m["idle"] = int64(0)
	}
	assertEqual(t, slice, []map[string]interface{}{
		{"idle": int64(0), "name": consumer1, "pending": int64(1)},
		{"idle": int64(0), "name": consumer2, "pending": int64(1)},
	})

	tree, err := conn.Xpending(myStream, myGroup, nil, nil).ToTree()
	assertNil(t, err)
	assertEqual(t, tree, []interface{}{
		//strange: number of pending per consumer returned as string?
		int64(2), id1, id2, []interface{}{[]interface{}{consumer1, "1"}, []interface{}{consumer2, "1"}},
	})

	i, err := conn.Xack(myStream, myGroup, []string{id1}).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)

	claimSlice, err := conn.Xclaim(myStream, myGroup, consumer1, "0", []string{id2}, nil, nil, nil, false, true).ToStringSlice()
	assertNil(t, err)
	assertEqual(t, claimSlice, []string{id2})

	i, err = conn.XgroupDelconsumer(myStream, myGroup, consumer1).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)

}

func testXgroupSetid(conn client.Conn, ctx *testCTX, t *testing.T) {
	myStream := ctx.newKey("myStream")
	myGroup := ctx.newKey("myGroup")
	myConsumer := ctx.newKey("myConsumer")
	id, err := conn.Xadd(myStream, "*", []client.FieldValue{{Field: "name", Value: "Sara"}, {Field: "surname", Value: "OConnor"}}).ToString()
	assertNil(t, err)
	ok, err := conn.XgroupCreate(myStream, myGroup, "$", true).ToBool()
	assertNil(t, err)
	assertTrue(t, ok)
	m, err := conn.Xreadgroup(client.GroupConsumer{Group: myGroup, Consumer: myConsumer}, nil, nil, false, []interface{}{myStream}, []string{">"}).ToXread()
	assertNil(t, err)
	assertEqual(t, len(m[myStream]), 0)
	err = conn.XgroupSetid(myStream, myGroup, "0").Err()
	assertNil(t, err)
	m, err = conn.Xreadgroup(client.GroupConsumer{Group: myGroup, Consumer: myConsumer}, nil, nil, false, []interface{}{myStream}, []string{">"}).ToXread()
	assertNil(t, err)
	assertEqual(t, m, map[string][]client.XItem{myStream: {
		{ID: id, Items: []string{"name", "Sara", "surname", "OConnor"}},
	}})
}

func testXgroupDestroy(conn client.Conn, ctx *testCTX, t *testing.T) {
	myStream := ctx.newKey("myStream")
	myGroup := ctx.newKey("myGroup")
	ok, err := conn.XgroupCreate(myStream, myGroup, "$", true).ToBool()
	assertNil(t, err)
	assertTrue(t, ok)
	ok, err = conn.XgroupDestroy(myStream, myGroup).ToBool()
	assertNil(t, err)
	assertTrue(t, ok)
}

func testXgroupHelp(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.XgroupHelp().ToStringSlice()
	assertNil(t, err)
}

func testXinfoStream(conn client.Conn, ctx *testCTX, t *testing.T) {
	myStream := ctx.newKey("myStream")
	err := conn.Xadd(myStream, "*", []client.FieldValue{{Field: "a", Value: "1"}}).Err()
	assertNil(t, err)
	_, err = conn.XinfoStream(myStream).ToStringMap()
	assertNil(t, err)
}

func testXinfoHelp(conn client.Conn, ctx *testCTX, t *testing.T) {
	_, err := conn.XinfoHelp().ToStringSlice()
	assertNil(t, err)
}

func testXlen(conn client.Conn, ctx *testCTX, t *testing.T) {
	myStream := ctx.newKey("myStream")
	err := conn.Xadd(myStream, "*", []client.FieldValue{{Field: "item", Value: "1"}}).Err()
	assertNil(t, err)
	err = conn.Xadd(myStream, "*", []client.FieldValue{{Field: "item", Value: "2"}}).Err()
	assertNil(t, err)
	err = conn.Xadd(myStream, "*", []client.FieldValue{{Field: "item", Value: "3"}}).Err()
	assertNil(t, err)
	i, err := conn.Xlen(myStream).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)
}

func testXrange(conn client.Conn, ctx *testCTX, t *testing.T) {
	writers := ctx.newKey("writers")
	id1, err := conn.Xadd(writers, "*", []client.FieldValue{{Field: "name", Value: "Virginia"}, {Field: "surname", Value: "Woolf"}}).ToString()
	assertNil(t, err)
	id2, err := conn.Xadd(writers, "*", []client.FieldValue{{Field: "name", Value: "Jane"}, {Field: "surname", Value: "Austen"}}).ToString()
	assertNil(t, err)
	err = conn.Xadd(writers, "*", []client.FieldValue{{Field: "name", Value: "Toni"}, {Field: "surname", Value: "Morris"}}).Err()
	assertNil(t, err)
	err = conn.Xadd(writers, "*", []client.FieldValue{{Field: "name", Value: "Agatha"}, {Field: "surname", Value: "Christie"}}).Err()
	assertNil(t, err)
	err = conn.Xadd(writers, "*", []client.FieldValue{{Field: "name", Value: "Ngozi"}, {Field: "surname", Value: "Adichie"}}).Err()
	assertNil(t, err)
	i, err := conn.Xlen(writers).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 5)
	slice, err := conn.Xrange(writers, "-", "+", client.Int64Ptr(2)).ToXrange()
	assertNil(t, err)
	assertEqual(t, slice, []client.XItem{
		{ID: id1, Items: []string{"name", "Virginia", "surname", "Woolf"}},
		{ID: id2, Items: []string{"name", "Jane", "surname", "Austen"}},
	})
}

func testXread(conn client.Conn, ctx *testCTX, t *testing.T) {
	writers := ctx.newKey("writers")
	id1, err := conn.Xadd(writers, "*", []client.FieldValue{{Field: "name", Value: "Virginia"}, {Field: "surname", Value: "Woolf"}}).ToString()
	assertNil(t, err)
	id2, err := conn.Xadd(writers, "*", []client.FieldValue{{Field: "name", Value: "Jane"}, {Field: "surname", Value: "Austen"}}).ToString()
	assertNil(t, err)
	err = conn.Xadd(writers, "*", []client.FieldValue{{Field: "name", Value: "Toni"}, {Field: "surname", Value: "Morris"}}).Err()
	assertNil(t, err)
	err = conn.Xadd(writers, "*", []client.FieldValue{{Field: "name", Value: "Agatha"}, {Field: "surname", Value: "Christie"}}).Err()
	assertNil(t, err)
	err = conn.Xadd(writers, "*", []client.FieldValue{{Field: "name", Value: "Ngozi"}, {Field: "surname", Value: "Adichie"}}).Err()
	assertNil(t, err)
	m, err := conn.Xread(client.Int64Ptr(2), nil, []interface{}{writers}, []string{"0-0"}).ToXread()
	assertNil(t, err)
	assertEqual(t, m, map[string][]client.XItem{writers: {
		{ID: id1, Items: []string{"name", "Virginia", "surname", "Woolf"}},
		{ID: id2, Items: []string{"name", "Jane", "surname", "Austen"}},
	}})
}

func testXrevrange(conn client.Conn, ctx *testCTX, t *testing.T) {
	writers := ctx.newKey("writers")
	err := conn.Xadd(writers, "*", []client.FieldValue{{Field: "name", Value: "Virginia"}, {Field: "surname", Value: "Woolf"}}).Err()
	assertNil(t, err)
	err = conn.Xadd(writers, "*", []client.FieldValue{{Field: "name", Value: "Jane"}, {Field: "surname", Value: "Austen"}}).Err()
	assertNil(t, err)
	err = conn.Xadd(writers, "*", []client.FieldValue{{Field: "name", Value: "Toni"}, {Field: "surname", Value: "Morris"}}).Err()
	assertNil(t, err)
	err = conn.Xadd(writers, "*", []client.FieldValue{{Field: "name", Value: "Agatha"}, {Field: "surname", Value: "Christie"}}).Err()
	assertNil(t, err)
	id5, err := conn.Xadd(writers, "*", []client.FieldValue{{Field: "name", Value: "Ngozi"}, {Field: "surname", Value: "Adichie"}}).ToString()
	assertNil(t, err)
	i, err := conn.Xlen(writers).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 5)
	slice, err := conn.Xrevrange(writers, "+", "-", client.Int64Ptr(1)).ToXrange()
	assertNil(t, err)
	assertEqual(t, slice, []client.XItem{{ID: id5, Items: []string{"name", "Ngozi", "surname", "Adichie"}}})
}

func testXtrim(conn client.Conn, ctx *testCTX, t *testing.T) {
	myStream := ctx.newKey("myStream")
	err := conn.Xadd(myStream, "*", []client.FieldValue{{Field: "field1", Value: "A"}}).Err()
	assertNil(t, err)
	err = conn.Xadd(myStream, "*", []client.FieldValue{{Field: "field2", Value: "B"}}).Err()
	assertNil(t, err)
	id3, err := conn.Xadd(myStream, "*", []client.FieldValue{{Field: "field3", Value: "C"}}).ToString()
	assertNil(t, err)
	id4, err := conn.Xadd(myStream, "*", []client.FieldValue{{Field: "field4", Value: "D"}}).ToString()
	assertNil(t, err)
	i, err := conn.Xtrim(myStream, false, 2).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)
	slice, err := conn.Xrange(myStream, "-", "+", client.Int64Ptr(2)).ToXrange()
	assertNil(t, err)
	assertEqual(t, slice, []client.XItem{
		{ID: id3, Items: []string{"field3", "C"}},
		{ID: id4, Items: []string{"field4", "D"}},
	})
}

// Transaction
// Transaction (cannot be executed in parallel - needs exclusive connection(s))
func testTransaction(conn client.Conn, ctx *testCTX, t *testing.T) {
	conn, err := ctx.dialer.Dial("")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	conn2, err := ctx.dialer.Dial("")
	if err != nil {
		t.Fatal(err)
	}
	defer conn2.Close()

	foo, bar := ctx.newKey("foo"), ctx.newKey("bar")

	// Exec
	ok, err := conn.Multi().ToBool()
	assertNil(t, err)
	assertTrue(t, ok)

	s, err := conn.Incr(foo).ToString()
	assertNil(t, err)
	assertEqual(t, s, client.ReplyQueued)

	s, err = conn.Incr(bar).ToString()
	assertNil(t, err)
	assertEqual(t, s, client.ReplyQueued)

	slice, err := conn.Exec().ToInt64Slice()
	assertNil(t, err)
	assertEqual(t, slice, []int64{1, 1})

	// Discard
	ok, err = conn.Multi().ToBool()
	assertNil(t, err)
	assertTrue(t, ok)

	s, err = conn.Incr(foo).ToString()
	assertNil(t, err)
	assertEqual(t, s, client.ReplyQueued)

	ok, err = conn.Discard().ToBool()
	assertNil(t, err)
	assertTrue(t, ok)

	i, err := conn.Get(foo).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 1)

	// Watch
	ok, err = conn.Watch([]interface{}{foo}).ToBool()
	assertNil(t, err)
	assertTrue(t, ok)

	i, err = conn2.Incr(foo).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 2)

	ok, err = conn.Multi().ToBool()
	assertNil(t, err)
	assertTrue(t, ok)

	s, err = conn.Incr(foo).ToString()
	assertNil(t, err)
	assertEqual(t, s, client.ReplyQueued)

	b, err := conn.Exec().IsNull()
	assertNil(t, err)
	assertTrue(t, b)

	// Unwatch
	ok, err = conn.Watch([]interface{}{foo}).ToBool()
	assertNil(t, err)
	assertTrue(t, ok)

	i, err = conn2.Incr(foo).ToInt64()
	assertNil(t, err)
	assertEqual(t, i, 3)

	ok, err = conn.Unwatch().ToBool()
	assertNil(t, err)
	assertTrue(t, ok)

	ok, err = conn.Multi().ToBool()
	assertNil(t, err)
	assertTrue(t, ok)

	s, err = conn.Incr(foo).ToString()
	assertNil(t, err)
	assertEqual(t, s, client.ReplyQueued)

	slice, err = conn.Exec().ToInt64Slice()
	assertNil(t, err)
	assertEqual(t, slice, []int64{4})
}

//
func interceptSend(tested map[string]bool, t *testing.T) client.SendInterceptor {
	var mu sync.Mutex

	return func(cmdName string, values []interface{}) {
		mu.Lock()
		tested[cmdName] = true
		mu.Unlock()
	}
}

func TestCommand(t *testing.T) {
	tested := map[string]bool{}
	excluded := map[string]bool{
		client.CmdAclLoad:       true,
		client.CmdAclSave:       true,
		client.CmdConfigRewrite: true,
		client.CmdDebugSegfault: true,
		client.CmdFlushall:      true,
		client.CmdFlushdb:       true,
		client.CmdMigrate:       true,
		client.CmdModuleLoad:    true,
		client.CmdModuleUnload:  true,
		client.CmdMonitor:       true,
		client.CmdScriptDebug:   true,
		client.CmdScriptFlush:   true,
		client.CmdScriptKill:    true,
		client.CmdShutdown:      true,
		client.CmdReplicaof:     true,
		client.CmdPsync:         true,
	}

	// exclude cluster commands
	for _, cmd := range client.Groups[client.GroupCluster] {
		excluded[cmd] = true
	}

	for _, cmd := range client.CommandNames {
		tested[cmd] = false
	}

	dialer := client.Dialer{Logger: log.New(os.Stderr, "", log.LstdFlags)}
	dialer.SendInterceptor = interceptSend(tested, t)

	conn, err := dialer.Dial("")
	if err != nil {
		t.Fatal(err)
	}

	ctx := newTestCTX(dialer)

	t.Run("parallel", func(t *testing.T) {
		for _, fct := range fcts {
			if fct.parallel {
				func(f testFct) {
					t.Run(fct.name, func(t *testing.T) {
						t.Parallel() // start tests in parallel
						f(conn, ctx, t)
					})
				}(fct.f)
			}
		}
	})

	t.Run("seriell", func(t *testing.T) {
		for _, fct := range fcts {
			if !fct.parallel {
				t.Run(fct.name, func(t *testing.T) {
					fct.f(conn, ctx, t)
				})
			}
		}
	})

	ctx.cleanup(conn)

	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}

	i, j, l := 0, 0, len(client.CommandNames)
	untested := []string{}
	for _, cmd := range client.CommandNames {
		switch {
		case excluded[cmd]:
			i++
		case !tested[cmd]:
			j++
			untested = append(untested, cmd)
		}
	}
	if j != 0 {
		t.Logf("totals: %d - tested %d excluded %d untested %d => %s", l, l-(i+j), i, j, strings.Join(untested, ", "))
	} else {
		t.Logf("totals: %d - tested %d excluded %d", l, l-(i+j), i)
	}
}

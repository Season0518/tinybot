package tinybot

//func TestOnebot_SendGroupMsg(t *testing.T) {
//	msgBuilder := NewMsgBuilder().Text("Hello").Text("你好").At("excluded").Image("https://bkimg.cdn.bcebos.com/pic/5fdf8db1cb134954c0e116395a4e9258d0094af4", false)
//	Conn, _, err := wsclient.DefaultDialer.Dial("Conn://excluded", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	//defer Conn.Close()
//	bot := wsclient.WSClient{Conn: Conn}
//	go bot.Run()
//	err = bot.SendGroupMsg(0, msgBuilder.Build())
//	if err != nil {
//		t.Fatal(err)
//	}
//}

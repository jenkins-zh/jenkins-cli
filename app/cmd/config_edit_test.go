package cmd

//func TestEditConfig(t *testing.T) {
//	RunEditCommandTest(t, EditCommandTest{
//		Procedure: func(c *expect.Console) {
//			c.ExpectString("Edit config item yourServer")
//			c.SendLine("")
//			go c.ExpectEOF()
//			time.Sleep(time.Millisecond)
//			c.Send("\x1b")
//			c.SendLine(":wq!")
//		},
//		Test: func(stdio terminal.Stdio) (err error) {
//			rootOptions.ConfigFile = "test.yaml"
//			data, err := GenerateSampleConfig()
//			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
//
//			rootCmd.SetArgs([]string{"config", "edit"})
//			configEditOption.Option.Stdio = stdio
//			_, err = rootCmd.ExecuteC()
//			return
//		},
//	})
//}

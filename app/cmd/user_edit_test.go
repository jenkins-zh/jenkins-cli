package cmd

//func TestEditUser(t *testing.T) {
//	RunEditCommandTest(t, EditCommandTest{
//		Procedure: func(c *expect.Console) {
//			c.ExpectString("Edit user description")
//			c.SendLine("")
//			go c.ExpectEOF()
//			time.Sleep(time.Millisecond)
//			c.Send("\x1b")
//			c.SendLine(":wq!")
//		},
//		Test: func(stdio terminal.Stdio) (err error) {
//			configFile := path.Join(os.TempDir(), "fake.yaml")
//			defer os.Remove(configFile)
//
//			data, err := GenerateSampleConfig()
//			err = ioutil.WriteFile(configFile, data, 0664)
//
//			var (
//				description = "fake-description\n"
//			)
//
//			ctrl := gomock.NewController(t)
//			roundTripper := mhttp.NewMockRoundTripper(ctrl)
//
//			client.PrepareGetUser(roundTripper, "http://localhost:8080/jenkins", "admin", "111e3a2f0231198855dceaff96f20540a9")
//
//			client.PrepareForEditUserDesc(roundTripper, "http://localhost:8080/jenkins",
//				"admin", description, "admin", "111e3a2f0231198855dceaff96f20540a9")
//
//			rootCmd.SetArgs([]string{"user", "edit", "--desc", description, "--configFile", configFile})
//
//			userEditOption.RoundTripper = roundTripper
//			userEditOption.Option.Stdio = stdio
//			_, err = rootCmd.ExecuteC()
//			return
//		},
//	})
//}

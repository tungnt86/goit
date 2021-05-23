package must

func NotFail(err error) {
	if nil != err {
		panic(err.Error())
	}
}

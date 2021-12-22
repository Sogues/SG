package battle

type (
	MovePreorder struct {
		velocity Vector3

		//多久完成，单位秒
		inTime float64

		//还有多久移动完成，单位：秒，如果小于1帧的时间但还大于0，就会当做1帧来执行
		duration float64
	}
)

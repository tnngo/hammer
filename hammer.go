package hammer

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/tnngo/hammer/logger"
	"github.com/tnngo/lad"
)

type hammer interface {
	Run()
}

type Builder interface {
	Build()
}

type Hammer struct {
	vrs []hammer

	// 是否已经构建.
	isBuild bool
	// 日志.
	ls []logger.Logger

	// 基础设施配件类型.
	fac facility

	fl facilityLoader
}

func (v *Hammer) preLoad() bool {
	var (
		errLoad error
	)

	var fac facility
	// 检测yaml文件.
	if v.fl == nil {
		y := &yml{}
		fac, errLoad = y.load()
	} else {
		if y, ok := v.fl.(*yml); ok {
			fac, errLoad = y.load()
		} else {
			fac, errLoad = v.fl.load()
		}
	}

	// 没有默认的配置文件
	if fac == nil && errLoad == nil {
		return false
	}

	if errLoad != nil {
		lad.L().Error(errLoad.Error())
		return false
	}

	if fac != nil {
		v.fac = fac
	}

	return true
}

// New 创建英灵殿实例。
func New(vrs ...hammer) *Hammer {
	v := &Hammer{
		vrs: vrs,
	}

	// 先初始化全局日志, 用于打印valhalla在装载配件过程中相关信息.
	lad.ReplaceGlobals(lad.New((&logger.Console{}).Mode(), lad.AddCaller()))

	// 预加载, 当前最高优先级为配置文件, 后期是etcd.
	if v.preLoad() {
		// 开始构建
		v.Build()
		v.isBuild = true
	}

	return v
}

func (v *Hammer) Build() {
	if v.isBuild {
		return
	}

	if v.fac == nil {
		return
	}

	_facMap = v.fac

	/** 构建orm **/
	if o := v.fac.loadOrm(); o != nil {
		o.Build()
	}

	/** 构建http **/
	if h := v.fac.loadHttp(); h != nil {
		if len(v.vrs) != 0 {
			for i, v1 := range v.vrs {
				if _, ok := v1.(*Http); ok {
					h.s = newHttpServer()
					v.vrs[i] = h
				}
			}
		}
	}

}

func (v *Hammer) Http() *Http {
	if len(v.vrs) == 0 {
		v.vrs = append(v.vrs, NewHttpServer())
	}
	for _, v1 := range v.vrs {
		if v2, ok := v1.(*Http); ok {
			if v2.s == nil {
				v2.s = newHttpServer()
			}
			return v2
		}
	}
	return nil
}

func (v *Hammer) Run() {
	for _, v1 := range v.vrs {
		go v1.Run()
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	//等待停止信号
	<-c
}

// 构建执行, 常用于手动配置信息, 而不是通过配置文件.
func (v *Hammer) Builds(builder ...Builder) {
	for _, v1 := range builder {
		v1.Build()
	}
}

func localIP() string {
	localAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, v := range localAddrs {
		if ipnet, ok := v.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

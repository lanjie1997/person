package global
import (
    "wholeally.com/share/v4/dbapi"
    "wholeally.com/share/v4/dbentity"
    "wholeally.com/common/logs"
    "fmt"
)
// 数据库相关接口
// 获取所有onvif设备的
func DbGetOnvifDeviceList()[]string {
    return []string{"172.168.0.229","172.168.0.191","172.168.0.161"}
}

//"172.168.0.154"



// 设置资源同步状态
func DbSetResSyncState(code, curCode string) error {
    req := dbentity.SetResSyncStateRequest{
        ComCode:      curCode,
        ComCodeOther: code,
        Protocol:     1,
        SyncState:    2,
    }
    res := dbentity.SetResSyncStateResponse{}

    err := dbapi.Call("/cascade/setressyncstate", &req, &res)
    if nil != err {
        logs.Waring(err)
        return err
    }

    if 0 != res.Ret {
        err := fmt.Errorf(res.Msg)
        logs.Waring(err)
        return err
    }

    return nil
}
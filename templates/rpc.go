package x

import (
    "strings"
    "github.com/golang/protobuf/proto"
    "errors"
    "ms/sun/config"
    "log"
)

type RPC_UserParam interface {
    GetUserId() int
    IsUser() bool
}

type RPC_ResponseHandlerInterface interface {
    HandleOfflineResult(interface{}, PB_CommandToServer, RPC_UserParam)
    IsUserOnlineResult(interface{}, error)
    HandelError(error)
}

var RPC_ResponseHandler RPC_ResponseHandlerInterface

//note: rpc methods cant have equal name they must be different even in different rpc services
type RPC_AllHandlersInteract interface {
{{range .Services}}
   {{.Name}}
{{- end}}
}

/////////////// Interfaces ////////////////
{{range .Services}}
type {{.Name}} interface {
{{- range .Methods}}
    {{.MethodName}}(*{{.InTypeName}} ) (*{{.OutTypeName}} ,error)
{{- end}}
}
{{end}}

func noDevErr(err error)  {
    if config.IS_DEBUG && err != nil {
        log.Panic(err)
    }
}

////////////// map of rpc methods to all
func HandleRpcs(cmd PB_CommandToServer, params RPC_UserParam, rpcHandler RPC_AllHandlersInteract) {

    splits := strings.Split(cmd.Command, ".")

    if len(splits) != 2 {
        noDevErr(errors.New("HandleRpcs: splic is not 2 parts"))
        return
    }

    switch splits[0] {
{{range .Services}}
    case "{{.Name}}":
            rpc,ok := rpcHandler.({{.Name}})
            if !ok {
                e:=errors.New("rpcHandler could not be cast to : {{.Name}}")
                noDevErr(e)
                RPC_ResponseHandler.HandelError(e)
                return
            }

            switch splits[1]  {
            {{- range .Methods}}
                case "{{.MethodName}}": //each pb_service_method
                    load := &{{.InTypeName}}{}
                    err := proto.Unmarshal(cmd.Data, load)
                    if err == nil {
                        res, err := rpc.{{.MethodName}}(load)
                        RPC_ResponseHandler.HandleOfflineResult(res,cmd, params)
                    }else{
                     RPC_ResponseHandler.HandelError(err)
                    }
            {{- end}}
            }
{{- end}}
}
}
package commands

import (
	"errors"
	"fmt"
	"github.com/apache/brooklyn-client/api/entity_effectors"
	"github.com/apache/brooklyn-client/command_metadata"
	"github.com/apache/brooklyn-client/error_handler"
	"github.com/apache/brooklyn-client/net"
	"github.com/apache/brooklyn-client/scope"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"strings"
)

type Invoke struct {
	network *net.Network
}

type Stop struct {
	Invoke
}

type Start struct {
	Invoke
}

type Restart struct {
	Invoke
}

func NewInvoke(network *net.Network) (cmd *Invoke) {
	cmd = new(Invoke)
	cmd.network = network
	return
}

func NewInvokeStop(network *net.Network) (cmd *Stop) {
	cmd = new(Stop)
	cmd.network = network
	return
}

func NewInvokeStart(network *net.Network) (cmd *Start) {
	cmd = new(Start)
	cmd.network = network
	return
}

func NewInvokeRestart(network *net.Network) (cmd *Restart) {
	cmd = new(Restart)
	cmd.network = network
	return
}

var paramFlags = []cli.Flag{
	cli.StringSliceFlag{
		Name:  "param, P",
		Usage: "Parameter and value separated by '=', e.g. -P x=y",
	},
}

func (cmd *Invoke) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "invoke",
		Description: "Invoke an effector of an application and entity",
		Usage:       "BROOKLYN_NAME EFF-SCOPE invoke [ parameter-options ]",
		Flags:       paramFlags,
	}
}

func (cmd *Stop) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "stop",
		Description: "Invoke stop effector on an application and entity",
		Usage:       "BROOKLYN_NAME ENT-SCOPE stop [ parameter-options ]",
		Flags:       paramFlags,
	}
}

func (cmd *Start) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "start",
		Description: "Invoke start effector on an application and entity",
		Usage:       "BROOKLYN_NAME ENT-SCOPE start [ parameter-options ]",
		Flags:       paramFlags,
	}
}

func (cmd *Restart) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "restart",
		Description: "Invoke restart effector on an application and entity",
		Usage:       "BROOKLYN_NAME ENT-SCOPE restart [ parameter-options ]",
		Flags:       paramFlags,
	}
}

func (cmd *Invoke) Run(scope scope.Scope, c *cli.Context) {
	if err := net.VerifyLoginURL(cmd.network); err != nil {
		error_handler.ErrorExit(err)
	}
	parms := c.StringSlice("param")
	invoke(cmd.network, scope.Application, scope.Entity, scope.Effector, parms)
}

const stopEffector = "stop"

func (cmd *Stop) Run(scope scope.Scope, c *cli.Context) {
	if err := net.VerifyLoginURL(cmd.network); err != nil {
		error_handler.ErrorExit(err)
	}
	parms := c.StringSlice("param")
	invoke(cmd.network, scope.Application, scope.Entity, stopEffector, parms)
}

const startEffector = "start"

func (cmd *Start) Run(scope scope.Scope, c *cli.Context) {
	if err := net.VerifyLoginURL(cmd.network); err != nil {
		error_handler.ErrorExit(err)
	}
	parms := c.StringSlice("param")
	invoke(cmd.network, scope.Application, scope.Entity, startEffector, parms)
}

const restartEffector = "restart"

func (cmd *Restart) Run(scope scope.Scope, c *cli.Context) {
	if err := net.VerifyLoginURL(cmd.network); err != nil {
		error_handler.ErrorExit(err)
	}
	parms := c.StringSlice("param")
	invoke(cmd.network, scope.Application, scope.Entity, restartEffector, parms)
}

func invoke(network *net.Network, application, entity, effector string, parms []string) {
	names, vals, err := extractParams(parms)
	result, err := entity_effectors.TriggerEffector(network, application, entity, effector, names, vals)
	if nil != err {
		error_handler.ErrorExit(err)
	} else {
		if "" != result {
			fmt.Println(result)
		}
	}
}

func extractParams(parms []string) ([]string, []string, error) {
	names := make([]string, len(parms))
	vals := make([]string, len(parms))
	var err error
	for i, parm := range parms {
		index := strings.Index(parm, "=")
		if index < 0 {
			return names, vals, errors.New("Parameter value not provided: " + parm)
		}
		names[i] = parm[0:index]
		vals[i], err = extractParamValue(parm[index+1:])
	}
	return names, vals, err
}

const paramDataPrefix string = "@"

func extractParamValue(rawParam string) (string, error) {
	var err error
	var val string
	if strings.HasPrefix(rawParam, paramDataPrefix) {
		// strip the data prefix from the filename before reading
		val, err = readParamFromFile(rawParam[len(paramDataPrefix):])
	} else {
		val = rawParam
		err = nil
	}
	return val, err
}

// returning a string rather than byte array, assuming non-binary
// TODO - if necessary support binary data sending to effector
func readParamFromFile(filename string) (string, error) {
	dat, err := ioutil.ReadFile(filename)
	return string(dat), err
}

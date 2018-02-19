package main

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/bryanl/ksstdlib/builder"
	"github.com/google/go-jsonnet/ast"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	nm "github.com/ksonnet/ksonnet-lib/ksonnet-gen/nodemaker"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/printer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	var k8sPath string
	flag.StringVar(&k8sPath, "k8sPath", "tmp/k8s.libsonnet", "Path to k8s.libsonnet")

	flag.Parse()

	err := builder.Build(k8sPath)
	if err != nil {
		logrus.WithError(err).Fatal("build")
	}

	o := nm.NewObject()

	out, err := printLib(o)
	if err != nil {
		logrus.WithError(err).Fatal("print library")
	}

	fmt.Println(string(out))
}

func printLib(node nm.Noder) ([]byte, error) {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, node.Node()); err != nil {
		return nil, errors.Wrap(err, "create jsonnet")
	}

	return buf.Bytes(), nil
}

type iterateObjectFn func(string, string, ast.Node) error

func iterateObject(node ast.Node, fn iterateObjectFn) error {
	if node == nil {
		return errors.New("node was nil")
	}

	obj, ok := node.(*astext.Object)
	if !ok {
		return errors.New("node was not an object")
	}

	for _, of := range obj.Fields {
		if of.Hide == ast.ObjectFieldInherit {
			continue
		}

		if of.Kind == ast.ObjectLocal {
			continue
		}

		id := string(*of.Id)
		if id == "hidden" {
			continue
		}

		var comment string
		if of.Comment != nil {
			comment = of.Comment.Text
		}

		if err := fn(id, comment, of.Expr2); err != nil {
			return err
		}
	}

	return nil
}

package builder

import (
	"github.com/bryanl/woowoo/node"
	"github.com/bryanl/woowoo/yaml2jsonnet"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

// Build builds the library.
func Build(path string) error {
	root, err := yaml2jsonnet.ImportJsonnet(path)
	if err != nil {
		return errors.Wrap(err, "import k8s.libsonnet")
	}

	members, err := node.FindMembers(root)
	if err != nil {
		return errors.Wrap(err, "find members in root")
	}

	for _, field := range members.Fields {
		if field == "hidden" {
			continue
		}

		group, err := node.Find(root, field)
		if err != nil {
			return errors.Wrapf(err, "find %s in root", field)
		}

		versionMembers, err := node.FindMembers(group)
		if err != nil {
			return errors.Wrapf(err, "find members in group %s", field)
		}

		for _, versionField := range versionMembers.Fields {
			version, err := node.Find(group, versionField)
			if err != nil {
				return errors.Wrapf(err, "find node %s in group %s", versionField, field)
			}

			kindMembers, err := node.FindMembers(version)
			if err != nil {
				return errors.Wrapf(err, "find members in group %s version %s", field, versionField)
			}

			spew.Dump(kindMembers.Functions)
		}
	}

	spew.Dump(members)

	return nil

}

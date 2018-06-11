// This file was auto-generated from YAML files

package v1_10

func init() {
	Profile["/tosca/kubernetes/1.10/js/plugins/dummy.js"] = `

plugin = {
	name: 'Dummy',
	process: function(clout, v, kubernetes) {
		kubernetes.push({
			apiVersion: 'v1',
			kind: 'dummy',
			data: v.ID
		})
		return kubernetes
	}
};
`
}

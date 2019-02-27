// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/kubernetes/1.0/js/plugins/dummy.js"] = `

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

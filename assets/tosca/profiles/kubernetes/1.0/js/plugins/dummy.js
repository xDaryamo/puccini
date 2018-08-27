
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
